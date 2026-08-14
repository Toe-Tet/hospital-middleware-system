package hospitals

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
	"hospital-middleware-system/src/modules/hospitals/model"
	"hospital-middleware-system/src/modules/hospitals/serializer"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	listFn func(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error)
}

func (m *mockService) List(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
	return m.listFn(ctx, page, perPage)
}

func newTestController(svc Service) *Controller {
	return &Controller{service: svc}
}

func performRequest(method, target string, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/hospitals", handler)
	req := httptest.NewRequest(method, target, nil)
	r.ServeHTTP(w, req)
	return w
}

func fixtureHospitals() []*model.Hospital {
	t := time.Date(2026, 8, 14, 10, 30, 0, 0, time.UTC)
	return []*model.Hospital{
		{
			ID:        1,
			Name:      "St. Mary's General Hospital",
			Address:   "123 Medical Lane, Springfield",
			Phone:     "+12175550123",
			Status:    "active",
			CreatedAt: t,
			UpdatedAt: t,
		},
		{
			ID:        2,
			Name:      "Bangkok General",
			Address:   "456 Sukhumvit Rd, Bangkok",
			Phone:     "+6621234567",
			Status:    "active",
			CreatedAt: t.Add(-time.Hour),
			UpdatedAt: t.Add(-time.Hour),
		},
	}
}

func TestController_List(t *testing.T) {
	hospitals := fixtureHospitals()

	tests := []struct {
		name           string
		query          string
		mockList       func(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error)
		wantStatus     int
		wantSuccess    bool
		checkItems     bool
		wantTotal      int64
		wantPage       int
		wantPerPage    int
		wantTotalPages int
		wantItemCount  int
		wantErrorCode  string
	}{
		{
			name:  "default pagination returns 200 with serialized items",
			query: "",
			mockList: func(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
				assert.Equal(t, 1, page)
				assert.Equal(t, 10, perPage)
				return hospitals, 2, nil
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
			name:  "custom page and per_page passed to service",
			query: "?page=2&per_page=5",
			mockList: func(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
				assert.Equal(t, 2, page)
				assert.Equal(t, 5, perPage)
				return hospitals[:1], 7, nil
			},
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			checkItems:     true,
			wantTotal:      7,
			wantPage:       2,
			wantPerPage:    5,
			wantTotalPages: 2,
			wantItemCount:  1,
		},
		{
			name:  "page below 1 clamps to 1",
			query: "?page=0&per_page=10",
			mockList: func(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
				assert.Equal(t, 1, page)
				assert.Equal(t, 10, perPage)
				return []*model.Hospital{}, 0, nil
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
			name:  "per_page over 100 clamps to 100",
			query: "?page=1&per_page=500",
			mockList: func(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
				assert.Equal(t, 1, page)
				assert.Equal(t, 100, perPage)
				return []*model.Hospital{}, 250, nil
			},
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			wantTotal:      250,
			wantPage:       1,
			wantPerPage:    100,
			wantTotalPages: 3,
		},
		{
			name:  "invalid page query ignores and uses default",
			query: "?page=notanumber&per_page=abc",
			mockList: func(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
				assert.Equal(t, 1, page)
				assert.Equal(t, 10, perPage)
				return hospitals, 2, nil
			},
			wantStatus:     http.StatusOK,
			wantSuccess:    true,
			wantTotal:      2,
			wantPage:       1,
			wantPerPage:    10,
			wantTotalPages: 1,
		},
		{
			name:  "service returns AppError maps to status and code",
			query: "",
			mockList: func(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
				return nil, 0, apperrors.NewInternal(errors.New("pgx conn busy"))
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
		{
			name:  "service returns raw error wrapped as internal 500",
			query: "",
			mockList: func(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
				return nil, 0, errors.New("boom")
			},
			wantStatus:    http.StatusInternalServerError,
			wantSuccess:   false,
			wantErrorCode: apperrors.CodeInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{listFn: tt.mockList}
			ctrl := newTestController(svc)

			w := performRequest(http.MethodGet, "/hospitals"+tt.query, ctrl.List)

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
					if !ok {
						t.Fatalf("expected items to be array, got %T", dataMap["items"])
					}
					assert.Len(t, rawItems, tt.wantItemCount)

					if len(rawItems) > 0 {
						first, _ := rawItems[0].(map[string]interface{})
						id, _ := first["id"].(float64)
						assert.EqualValues(t, hospitals[0].ID, int(id))
						assert.Equal(t, hospitals[0].Name, first["name"])
						assert.Equal(t, hospitals[0].Address, first["address"])
						assert.Equal(t, hospitals[0].Phone, first["phone"])
						assert.Equal(t, hospitals[0].Status, first["status"])
						assert.Equal(t,
							hospitals[0].CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
							first["created_at"],
						)
						assert.Equal(t,
							hospitals[0].UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
							first["updated_at"],
						)
					}
				}
			} else {
				errMap, ok := resp.Error.(map[string]interface{})
				if !ok {
					t.Fatalf("expected error to be object, got %T", resp.Error)
				}
				assert.Equal(t, tt.wantErrorCode, errMap["code"])
			}
		})
	}
}

func TestController_ListSerializerMatches(t *testing.T) {
	svc := &mockService{
		listFn: func(ctx context.Context, page, perPage int) ([]*model.Hospital, int64, error) {
			return fixtureHospitals(), 2, nil
		},
	}
	ctrl := newTestController(svc)

	w := performRequest(http.MethodGet, "/hospitals", ctrl.List)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Success bool                           `json:"success"`
		Data    struct {
			Items serializer.ListHospitalsResponseItems `json:"items"`
			Total int64                                 `json:"total"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.EqualValues(t, 2, resp.Data.Total)
	assert.Equal(t, serializer.SerializeListHospitals(fixtureHospitals()), resp.Data.Items)
}
