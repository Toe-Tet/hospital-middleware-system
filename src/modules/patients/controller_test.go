package patients

import (
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
	"hospital-middleware-system/src/modules/patients/dto"
	"hospital-middleware-system/src/modules/patients/model"
	"hospital-middleware-system/src/modules/patients/serializer"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	getByIDFn func(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error)
	listFn    func(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error)
}

func (m *mockService) GetByID(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error) {
	return m.getByIDFn(ctx, hospitalID, identifier)
}

func (m *mockService) List(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
	return m.listFn(ctx, hospitalID, filters, page, perPage)
}

const testHospitalID = 7

func newPatientsRouter(svc Service, authFn gin.HandlerFunc) (*gin.Engine, *Controller) {
	ctrl := &Controller{service: svc}
	r := gin.New()
	g := r.Group("/patients")
	if authFn != nil {
		g.Use(authFn)
	}
	g.GET("/", ctrl.List)
	g.GET("/:id", ctrl.GetByID)
	return r, ctrl
}

func fakeJWTMiddleware(hospitalID int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(string(middleware.StaffIDKey), 42)
		c.Set(string(middleware.HospitalIDKey), hospitalID)
		c.Next()
	}
}

func strPtr(s string) *string { return &s }
func timePtr(t time.Time) *time.Time { return &t }

func fixturePatients() []*model.Patient {
	t := time.Date(2026, 8, 14, 10, 30, 0, 0, time.UTC)
	dob := time.Date(1990, 5, 20, 0, 0, 0, 0, time.UTC)
	dob2 := time.Date(1985, 12, 5, 0, 0, 0, 0, time.UTC)
	return []*model.Patient{
		{
			ID:           101,
			HospitalID:   testHospitalID,
			PatientHN:    "HN-0001",
			FirstNameTh:  strPtr("สมชาย"),
			LastNameTh:   strPtr("ไทย"),
			FirstNameEn:  "John",
			MiddleNameEn: strPtr("Michael"),
			LastNameEn:   "Doe",
			NationalID:   strPtr("1100100100123"),
			PassportID:   strPtr("AB1234567"),
			DateOfBirth:  dob,
			Gender:       "M",
			Email:        "john@example.com",
			PhoneNumber:  "0812345678",
			Status:       "active",
			CreatedAt:    t,
			UpdatedAt:    timePtr(t.Add(time.Hour)),
		},
		{
			ID:           102,
			HospitalID:   testHospitalID,
			PatientHN:    "HN-0002",
			FirstNameEn:  "Jane",
			LastNameEn:   "Smith",
			DateOfBirth:  dob2,
			Gender:       "F",
			Email:        "jane@example.com",
			PhoneNumber:  "0898765432",
			Status:       "active",
			CreatedAt:    t.Add(-2 * time.Hour),
		},
	}
}

func TestController_GetByID(t *testing.T) {
	patients := fixturePatients()

	tests := []struct {
		name         string
		hospitalID   int
		pathID       string
		setupMock    func(m *mockService)
		wantStatus   int
		wantSuccess  bool
		wantItem     bool
		wantPatient  *model.Patient
		wantErrorCode string
	}{
		{
			name:       "national_id lookup returns scoped patient",
			hospitalID: testHospitalID,
			pathID:     "1100100100123",
			setupMock: func(m *mockService) {
				m.getByIDFn = func(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error) {
					assert.Equal(t, testHospitalID, hospitalID)
					assert.Equal(t, "1100100100123", identifier)
					return patients[0], nil
				}
			},
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantItem:    true,
			wantPatient: patients[0],
		},
		{
			name:       "passport_id lookup returns patient",
			hospitalID: testHospitalID,
			pathID:     "AB1234567",
			setupMock: func(m *mockService) {
				m.getByIDFn = func(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error) {
					assert.Equal(t, testHospitalID, hospitalID)
					assert.Equal(t, "AB1234567", identifier)
					return patients[1], nil
				}
			},
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantItem:    true,
			wantPatient: patients[1],
		},
		{
			name:       "identifier not in hospital scope returns 404",
			hospitalID: testHospitalID,
			pathID:     "ZZ-NOPE-999",
			setupMock: func(m *mockService) {
				m.getByIDFn = func(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error) {
					return nil, apperrors.NewNotFound("Patient not found")
				}
			},
			wantStatus:    http.StatusNotFound,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeNotFound,
		},
		{
			name:       "service returns internal db error",
			hospitalID: testHospitalID,
			pathID:     "1100100100123",
			setupMock: func(m *mockService) {
				m.getByIDFn = func(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error) {
					return nil, apperrors.NewInternal(errors.New("pgx: timeout"))
				}
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
		{
			name:       "empty identifier param returns 400 from service",
			hospitalID: testHospitalID,
			pathID:     "_EMPTY_",
			setupMock: func(m *mockService) {
				m.getByIDFn = func(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error) {
					return nil, apperrors.NewBadRequest("Patient identifier is required")
				}
			},
			wantStatus:    http.StatusBadRequest,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{}
			tt.setupMock(svc)
			r, _ := newPatientsRouter(svc, fakeJWTMiddleware(tt.hospitalID))

			req := httptest.NewRequest(http.MethodGet, "/patients/"+tt.pathID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp helper.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantSuccess, resp.Success)

			if tt.wantItem {
				raw, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "expected data to be object")
				id, _ := raw["id"].(float64)
				assert.EqualValues(t, tt.wantPatient.ID, int(id))
				assert.Equal(t, tt.wantPatient.FirstNameEn, raw["first_name_en"])
				assert.Equal(t, tt.wantPatient.LastNameEn, raw["last_name_en"])
				assert.Equal(t, tt.wantPatient.PhoneNumber, raw["phone_number"])
				assert.Equal(t, tt.wantPatient.Status, raw["status"])
				assert.Equal(t, tt.wantPatient.DateOfBirth.Format("2006-01-02"), raw["date_of_birth"])
				assert.Equal(t, tt.wantPatient.HospitalID, int(raw["hospital_id"].(float64)))
				hn, _ := raw["patient_hn"].(string)
				assert.Equal(t, tt.wantPatient.PatientHN, hn)
			}

			if !tt.wantSuccess {
				errMap, ok := resp.Error.(map[string]interface{})
				if !ok {
					t.Fatalf("expected error object, got %T", resp.Error)
				}
				assert.Equal(t, tt.wantErrorCode, errMap["code"])
			}
		})
	}
}

func TestController_List(t *testing.T) {
	patients := fixturePatients()

	tests := []struct {
		name            string
		hospitalID      int
		query           string
		setupMock       func(m *mockService)
		wantStatus      int
		wantSuccess     bool
		checkItems      bool
		wantTotal       int64
		wantPage        int
		wantPerPage     int
		wantTotalPages  int
		wantItemCount   int
		wantErrorCode   string
	}{
		{
			name:       "default page/per_page passes through with hospital tenancy",
			hospitalID: testHospitalID,
			query:      "",
			setupMock: func(m *mockService) {
				m.listFn = func(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
					assert.Equal(t, testHospitalID, hospitalID)
					assert.Equal(t, 1, page)
					assert.Equal(t, 10, perPage)
					assert.NotNil(t, filters)
					return patients, 2, nil
				}
			},
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			checkItems:     true,
			wantTotal:      2,
			wantPage:       1,
			wantPerPage:    10,
			wantTotalPages: 1,
			wantItemCount:  2,
		},
		{
			name:       "all filter query params bound to PatientFilters",
			hospitalID: testHospitalID,
			query:      "?national_id=1100100100123&passport_id=AB1234567&first_name=John&middle_name=Michael&last_name=Doe&date_of_birth=1990-05-20&email=john%40example.com&phone_number=081&page=2&per_page=5",
			setupMock: func(m *mockService) {
				m.listFn = func(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
					assert.Equal(t, testHospitalID, hospitalID)
					assert.Equal(t, "1100100100123", filters.NationalID)
					assert.Equal(t, "AB1234567", filters.PassportID)
					assert.Equal(t, "John", filters.FirstName)
					assert.Equal(t, "Michael", filters.MiddleName)
					assert.Equal(t, "Doe", filters.LastName)
					assert.Equal(t, "1990-05-20", filters.DateOfBirth)
					assert.Equal(t, "john@example.com", filters.Email)
					assert.Equal(t, "081", filters.PhoneNumber)
					assert.Equal(t, 2, page)
					assert.Equal(t, 5, perPage)
					return patients[:1], 1, nil
				}
			},
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			checkItems:     true,
			wantTotal:      1,
			wantPage:       2,
			wantPerPage:    5,
			wantTotalPages: 1,
			wantItemCount:  1,
		},
		{
			name:       "page 0 clamped to 1, per_page 0 to 10",
			hospitalID: testHospitalID,
			query:      "?page=0&per_page=0",
			setupMock: func(m *mockService) {
				m.listFn = func(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
					assert.Equal(t, 1, page)
					assert.Equal(t, 10, perPage)
					return []*model.Patient{}, 0, nil
				}
			},
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			checkItems:     true,
			wantTotal:      0,
			wantPage:       1,
			wantPerPage:    10,
			wantTotalPages: 0,
			wantItemCount:  0,
		},
		{
			name:       "per_page 500 clamped to 100",
			hospitalID: testHospitalID,
			query:      "?page=1&per_page=100",
			setupMock: func(m *mockService) {
				m.listFn = func(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
					assert.Equal(t, 100, perPage)
					return []*model.Patient{}, 250, nil
				}
			},
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			wantTotal:      250,
			wantPage:       1,
			wantPerPage:    100,
			wantTotalPages: 3,
		},
		{
			name:       "invalid per_page triggers validation error 400",
			hospitalID: testHospitalID,
			query:      "?per_page=-5",
			setupMock: func(m *mockService) {
				m.listFn = func(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
					t.Fatal("service should not be called on validation failure")
					return nil, 0, nil
				}
			},
			wantStatus:    http.StatusBadRequest,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeValidationError,
		},
		{
			name:       "per_page 501 triggers validation error",
			hospitalID: testHospitalID,
			query:      "?per_page=501",
			setupMock: func(m *mockService) {
				m.listFn = func(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
					t.Fatal("service should not be called on validation failure")
					return nil, 0, nil
				}
			},
			wantStatus:    http.StatusBadRequest,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeValidationError,
		},
		{
			name:       "invalid date_of_birth bubbles service 400",
			hospitalID: testHospitalID,
			query:      "?date_of_birth=not-a-date",
			setupMock: func(m *mockService) {
				m.listFn = func(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
					assert.Equal(t, "not-a-date", filters.DateOfBirth)
					return nil, 0, apperrors.NewBadRequest("Invalid date format, use YYYY-MM-DD")
				}
			},
			wantStatus:    http.StatusBadRequest,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeBadRequest,
		},
		{
			name:       "database internal error",
			hospitalID: testHospitalID,
			query:      "",
			setupMock: func(m *mockService) {
				m.listFn = func(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
					return nil, 0, apperrors.NewInternal(errors.New("connection refused"))
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
			r, _ := newPatientsRouter(svc, fakeJWTMiddleware(tt.hospitalID))

			req := httptest.NewRequest(http.MethodGet, "/patients/"+tt.query, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp helper.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantSuccess, resp.Success)

			if tt.wantSuccess {
				dataMap, ok := resp.Data.(map[string]interface{})
				if !ok {
					t.Fatalf("expected data to be object, got %T", resp.Data)
				}

				total, _ := dataMap["total"].(float64)
				page, _ := dataMap["page"].(float64)
				perPage, _ := dataMap["per_page"].(float64)
				totalPages, _ := dataMap["total_pages"].(float64)

				assert.EqualValues(t, tt.wantTotal, total)
				assert.EqualValues(t, tt.wantPage, page)
				assert.EqualValues(t, tt.wantPerPage, perPage)
				assert.EqualValues(t, tt.wantTotalPages, totalPages)

				if tt.checkItems {
					rawItems, ok := dataMap["items"].([]interface{})
					assert.True(t, ok, "expected items array, got %T", dataMap["items"])
					assert.Len(t, rawItems, tt.wantItemCount)
				}
			} else {
				errMap, ok := resp.Error.(map[string]interface{})
				assert.True(t, ok, "expected error object, got %T", resp.Error)
				assert.Equal(t, tt.wantErrorCode, errMap["code"])
			}
		})
	}
}

func TestController_ListSerializerMatches(t *testing.T) {
	patients := fixturePatients()
	svc := &mockService{
		listFn: func(ctx context.Context, hospitalID int, filters *dto.PatientFilters, page, perPage int) ([]*model.Patient, int64, error) {
			return patients, 2, nil
		},
	}
	r, _ := newPatientsRouter(svc, fakeJWTMiddleware(testHospitalID))

	req := httptest.NewRequest(http.MethodGet, "/patients/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			Items serializer.ListPatientsResponseItems `json:"items"`
			Total int64                                `json:"total"`
			Page  int                                  `json:"page"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.EqualValues(t, 2, resp.Data.Total)
	assert.Equal(t, 1, resp.Data.Page)
	assert.Equal(t, serializer.SerializeListPatients(patients), resp.Data.Items)
}

func TestController_GetByIDSerializerMatches(t *testing.T) {
	patients := fixturePatients()
	svc := &mockService{
		getByIDFn: func(ctx context.Context, hospitalID int, identifier string) (*model.Patient, error) {
			assert.Equal(t, testHospitalID, hospitalID)
			assert.Equal(t, "HN-0001", identifier)
			return patients[0], nil
		},
	}
	r, _ := newPatientsRouter(svc, fakeJWTMiddleware(testHospitalID))

	req := httptest.NewRequest(http.MethodGet, "/patients/HN-0001", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Success bool                         `json:"success"`
		Data    serializer.GetPatientResponse `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	expected := serializer.SerializeGetPatient(patients[0])
	assert.Equal(t, expected, &resp.Data)
}
