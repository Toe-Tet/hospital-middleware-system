package staffs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apperrors "hospital-middleware-system/src/errors"
	"hospital-middleware-system/src/helper"
	"hospital-middleware-system/src/middleware"
	"hospital-middleware-system/src/modules/staffs/dto"
	"hospital-middleware-system/src/modules/staffs/model"
	"hospital-middleware-system/src/modules/staffs/serializer"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	createFn  func(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error)
	getByIDFn func(ctx context.Context, id int) (*model.Staff, error)
	loginFn   func(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error)
	meFn      func(ctx context.Context, staffID int) (any, error)
}

func (m *mockService) Create(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error) {
	return m.createFn(ctx, req)
}

func (m *mockService) GetByID(ctx context.Context, id int) (*model.Staff, error) {
	if m.getByIDFn == nil {
		return nil, nil
	}
	return m.getByIDFn(ctx, id)
}

func (m *mockService) Login(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error) {
	return m.loginFn(ctx, req)
}

func (m *mockService) Me(ctx context.Context, staffID int) (any, error) {
	return m.meFn(ctx, staffID)
}

func staffStrPtr(s string) *string { cp := s; return &cp }
func staffTimePtr(t time.Time) *time.Time { return &t }

func fixtureStaff() *model.Staff {
	t := time.Date(2026, 8, 14, 10, 30, 0, 0, time.UTC)
	return &model.Staff{
		ID:         7,
		HospitalID: 1,
		Username:   "alice.carter",
		Name:       staffStrPtr("Dr. Alice Carter"),
		Password:   "hashed-secret",
		Status:     "active",
		CreatedAt:  staffTimePtr(t),
		UpdatedAt:  staffTimePtr(t.Add(time.Hour)),
	}
}

func newStaffsRouter(svc Service, authFn gin.HandlerFunc) (*gin.Engine, *Controller) {
	ctrl := &Controller{service: svc}
	r := gin.New()

	r.POST("/staffs", ctrl.Create)
	r.POST("/staffs/login", ctrl.Login)

	me := r.Group("/staffs")
	if authFn != nil {
		me.Use(authFn)
	}
	me.GET("/me", ctrl.Me)

	return r, ctrl
}

func fakeStaffJWT(staffID int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(string(middleware.StaffIDKey), staffID)
		c.Set(string(middleware.HospitalIDKey), 1)
		c.Set(string(middleware.RoleKey), "admin")
		c.Next()
	}
}

func jsonBody(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

type errorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestController_Create(t *testing.T) {
	staff := fixtureStaff()
	validReq := dto.CreateStaffRequest{
		HospitalID: 1,
		Username:   "alice.carter",
		Password:   "SecurePass123!",
	}

	tests := []struct {
		name          string
		body          any
		setupMock     func(m *mockService)
		wantStatus    int
		wantSuccess   bool
		checkData     bool
		wantStaff     *model.Staff
		wantErrorCode string
	}{
		{
			name: "valid request creates staff via serializer",
			body: validReq,
			setupMock: func(m *mockService) {
				m.createFn = func(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error) {
					assert.Equal(t, validReq.HospitalID, req.HospitalID)
					assert.Equal(t, validReq.Username, req.Username)
					assert.Equal(t, validReq.Password, req.Password)
					return staff, nil
				}
			},
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			checkData:   true,
			wantStaff:   staff,
		},
		{
			name: "payload missing required fields returns internal error",
			body: struct{ Junk string }{Junk: "junk"},
			setupMock: func(m *mockService) {
				m.createFn = func(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error) {
					t.Fatal("service should not be called on bind failure")
					return nil, nil
				}
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
		{
			name: "missing required password returns internal error",
			body: dto.CreateStaffRequest{
				HospitalID: 1,
				Username:   "alice.carter",
			},
			setupMock: func(m *mockService) {
				m.createFn = func(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error) {
					t.Fatal("service should not be called on validation failure")
					return nil, nil
				}
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
		{
			name: "username too short returns internal error",
			body: dto.CreateStaffRequest{
				HospitalID: 1,
				Username:   "ab",
				Password:   "SecurePass123!",
			},
			setupMock: func(m *mockService) {
				m.createFn = func(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error) {
					t.Fatal("service should not be called on validation failure")
					return nil, nil
				}
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
		{
			name: "duplicate username → 409 CONFLICT",
			body: validReq,
			setupMock: func(m *mockService) {
				m.createFn = func(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error) {
					return nil, apperrors.NewConflict("staff with this username already exists")
				}
			},
			wantStatus:    http.StatusConflict,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeConflict,
		},
		{
			name: "hospital does not exist → 400 BAD_REQUEST",
			body: validReq,
			setupMock: func(m *mockService) {
				m.createFn = func(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error) {
					return nil, apperrors.NewBadRequest("hospital does not exist")
				}
			},
			wantStatus:    http.StatusBadRequest,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeBadRequest,
		},
		{
			name: "db internal error → 500",
			body: validReq,
			setupMock: func(m *mockService) {
				m.createFn = func(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error) {
					return nil, apperrors.NewInternal(errors.New("pgx: connection refused"))
				}
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{}
			tt.setupMock(svc)
			r, _ := newStaffsRouter(svc, nil)

			req := httptest.NewRequest(http.MethodPost, "/staffs", jsonBody(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantSuccess {
				var resp struct {
					Success bool                           `json:"success"`
					Data    serializer.CreateStaffResponse `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Success)

				if tt.checkData {
					expected := serializer.SerializeCreateStaff(tt.wantStaff)
					assert.Equal(t, expected, &resp.Data)
				}
			} else {
				var resp errorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
				if tt.wantErrorCode != "" {
					assert.Equal(t, tt.wantErrorCode, resp.Error.Code)
				}
			}
		})
	}
}

func TestController_Login(t *testing.T) {
	staff := fixtureStaff()
	loginReq := dto.LoginRequest{Username: "alice.carter", Password: "SecurePass123!"}
	tokenResp := &serializer.LoginResponse{
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.stub",
		ExpiresAt: 1_700_000_000,
		Staff:     serializer.MeResponse(staff),
	}

	tests := []struct {
		name          string
		body          any
		setupMock     func(m *mockService)
		wantStatus    int
		wantSuccess   bool
		checkToken    bool
		wantErrorCode string
		wantMsgSubstr string
	}{
		{
			name: "valid credentials returns token + staff profile",
			body: loginReq,
			setupMock: func(m *mockService) {
				m.loginFn = func(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error) {
					assert.Equal(t, loginReq.Username, req.Username)
					assert.Equal(t, loginReq.Password, req.Password)
					return tokenResp, nil
				}
			},
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			checkToken:  true,
		},
		{
			name: "missing password returns internal error",
			body: dto.LoginRequest{Username: "alice.carter"},
			setupMock: func(m *mockService) {
				m.loginFn = func(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error) {
					t.Fatal("service should not be called")
					return nil, nil
				}
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
		{
			name: "empty JSON object returns internal error",
			body: struct{}{},
			setupMock: func(m *mockService) {
				m.loginFn = func(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error) {
					t.Fatal("service should not be called")
					return nil, nil
				}
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
		{
			name: "wrong credentials → 401 Unauthorized",
			body: dto.LoginRequest{Username: "alice.carter", Password: "wrong"},
			setupMock: func(m *mockService) {
				m.loginFn = func(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error) {
					return nil, apperrors.NewUnauthorized("Invalid username or password")
				}
			},
			wantStatus:    http.StatusUnauthorized,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeUnauthorized,
		},
		{
			name: "inactive account → 401 Unauthorized",
			body: loginReq,
			setupMock: func(m *mockService) {
				m.loginFn = func(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error) {
					return nil, apperrors.NewUnauthorized("Account is inactive")
				}
			},
			wantStatus:    http.StatusUnauthorized,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeUnauthorized,
		},
		{
			name: "token generation failure → 500",
			body: loginReq,
			setupMock: func(m *mockService) {
				m.loginFn = func(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error) {
					return nil, apperrors.NewInternal(errors.New("signing key missing"))
				}
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{}
			tt.setupMock(svc)
			r, _ := newStaffsRouter(svc, nil)

			req := httptest.NewRequest(http.MethodPost, "/staffs/login", jsonBody(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantSuccess {
				var resp struct {
					Success bool                    `json:"success"`
					Data    serializer.LoginResponse `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Success)

				if tt.checkToken {
					assert.Equal(t, tokenResp.Token, resp.Data.Token)
					assert.EqualValues(t, tokenResp.ExpiresAt, resp.Data.ExpiresAt)
					assert.Equal(t, tokenResp.Staff, resp.Data.Staff)
				}
			} else {
				var resp errorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
				if tt.wantErrorCode != "" {
					assert.Equal(t, tt.wantErrorCode, resp.Error.Code)
				}
			}
		})
	}
}

func TestController_Me(t *testing.T) {
	staff := fixtureStaff()
	meResp := serializer.MeResponse(staff)

	tests := []struct {
		name          string
		staffID       int
		auth          gin.HandlerFunc
		setupMock     func(m *mockService)
		wantStatus    int
		wantSuccess   bool
		checkMe       bool
		wantMe        *serializer.StaffMeResponse
		wantErrorCode string
	}{
		{
			name:    "authenticated call returns MeResponse for staff_id=7",
			staffID: 7,
			auth:    fakeStaffJWT(7),
			setupMock: func(m *mockService) {
				m.meFn = func(ctx context.Context, staffID int) (any, error) {
					assert.Equal(t, 7, staffID)
					return meResp, nil
				}
			},
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			checkMe:     true,
			wantMe:      meResp,
		},
		{
			name:    "authenticated call with different staff_id passed through",
			staffID: 42,
			auth:    fakeStaffJWT(42),
			setupMock: func(m *mockService) {
				m.meFn = func(ctx context.Context, staffID int) (any, error) {
					assert.Equal(t, 42, staffID)
					return meResp, nil
				}
			},
			wantStatus:  http.StatusOK,
			wantSuccess: true,
		},
		{
			name:    "staff not found → 404",
			staffID: 999,
			auth:    fakeStaffJWT(999),
			setupMock: func(m *mockService) {
				m.meFn = func(ctx context.Context, staffID int) (any, error) {
					return nil, apperrors.NewNotFound("Staff not found")
				}
			},
			wantStatus:    http.StatusNotFound,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeNotFound,
		},
		{
			name:    "db error bubbles 500",
			staffID: 7,
			auth:    fakeStaffJWT(7),
			setupMock: func(m *mockService) {
				m.meFn = func(ctx context.Context, staffID int) (any, error) {
					return nil, apperrors.NewInternal(errors.New("tx failed"))
				}
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{}
			tt.setupMock(svc)
			r, _ := newStaffsRouter(svc, tt.auth)

			req := httptest.NewRequest(http.MethodGet, "/staffs/me", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantSuccess {
				var resp struct {
					Success bool                      `json:"success"`
					Data    serializer.StaffMeResponse `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.True(t, resp.Success)

				if tt.checkMe {
					assert.Equal(t, tt.wantMe, &resp.Data)
				}
			} else {
				var resp errorResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.False(t, resp.Success)
				if tt.wantErrorCode != "" {
					assert.Equal(t, tt.wantErrorCode, resp.Error.Code)
				}
			}
		})
	}
}

func TestController_CreateSerializerMatches(t *testing.T) {
	staff := fixtureStaff()
	svc := &mockService{
		createFn: func(ctx context.Context, req *dto.CreateStaffRequest) (*model.Staff, error) {
			return staff, nil
		},
	}
	r, _ := newStaffsRouter(svc, nil)

	body := dto.CreateStaffRequest{
		HospitalID: 1,
		Username:   "alice.carter",
		Password:   "SecurePass123!",
	}
	req := httptest.NewRequest(http.MethodPost, "/staffs", jsonBody(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Success bool                           `json:"success"`
		Data    serializer.CreateStaffResponse `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, serializer.SerializeCreateStaff(staff), &resp.Data)
}

func TestController_LoginSerializerMatches(t *testing.T) {
	staff := fixtureStaff()
	loginResp := &serializer.LoginResponse{
		Token:     "stub-token",
		ExpiresAt: 987654321,
		Staff:     serializer.MeResponse(staff),
	}
	svc := &mockService{
		loginFn: func(ctx context.Context, req *dto.LoginRequest) (*serializer.LoginResponse, error) {
			return loginResp, nil
		},
	}
	r, _ := newStaffsRouter(svc, nil)

	req := httptest.NewRequest(http.MethodPost, "/staffs/login", jsonBody(dto.LoginRequest{
		Username: "alice.carter",
		Password: "SecurePass123!",
	}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Success bool                    `json:"success"`
		Data    serializer.LoginResponse `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, loginResp, &resp.Data)
}

func TestController_MeSerializerMatches(t *testing.T) {
	staff := fixtureStaff()
	svc := &mockService{
		meFn: func(ctx context.Context, staffID int) (any, error) {
			return serializer.MeResponse(staff), nil
		},
	}
	r, _ := newStaffsRouter(svc, fakeStaffJWT(7))

	req := httptest.NewRequest(http.MethodGet, "/staffs/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Success bool                      `json:"success"`
		Data    serializer.StaffMeResponse `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, serializer.MeResponse(staff), &resp.Data)
}

var _ = helper.OK
