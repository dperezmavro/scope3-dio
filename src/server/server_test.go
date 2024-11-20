package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/dperezmavro/scope3-dio/src/common"
	"github.com/dperezmavro/scope3-dio/src/config"
)

func TestRoutes(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		method  string
		headers map[string]string
		status  int
		body    string
	}{
		{
			name:   "path no exists",
			path:   "/noexist",
			status: http.StatusNotFound,
			method: http.MethodGet,
		},
		{
			name:   "healtch",
			path:   "/healthcheck",
			status: http.StatusOK,
			method: http.MethodGet,
		},
		{
			name:   "missing auth",
			method: http.MethodPost,
			path:   "/v2/measure",
			status: http.StatusUnauthorized,
		},
		{
			name:   "missing rows",
			method: http.MethodPost,
			path:   "/v2/measure",
			status: http.StatusBadRequest,
			headers: map[string]string{
				common.HeaderAuthorization: "dummy",
			},
			body: `{}`,
		},
		{
			name:   "empty rows",
			method: http.MethodPost,
			path:   "/v2/measure",
			status: http.StatusBadRequest,
			headers: map[string]string{
				common.HeaderAuthorization: "dummy",
			},
			body: `{"rows": []}`,
		},
		{
			name:   "empty object",
			method: http.MethodPost,
			path:   "/v2/measure",
			status: http.StatusOK,
			headers: map[string]string{
				common.HeaderAuthorization: "dummy",
			},
			body: `{"rows": [{}]}`,
		},
		{
			method: http.MethodPost,
			name:   "well-formed request",
			path:   "/v2/measure",
			status: http.StatusOK,
			headers: map[string]string{
				common.HeaderAuthorization: "dummy",
			},
			body: `{"rows": [{"propertyId":"nytimes.com","urcDateTime":"2024-10-30","impressions":1000}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body *bytes.Buffer
			body = bytes.NewBuffer([]byte(tt.body))
			request := httptest.NewRequest(tt.method, tt.path, body)
			if tt.headers != nil {
				for k, v := range tt.headers {
					request.Header.Add(k, v)
				}
			}
			responseRecorder := httptest.NewRecorder()

			// mock config
			c := config.Config{
				Environment: config.Environment{
					Name: "local",
				},
				Port: 3000,
				Service: config.Service{
					Name:    "unitTest",
					Version: 1,
				},
				Scope3APIToken: "dummy",
			}

			// mock storage client
			sc := MockStorage{}

			h := CreateRouter(c, sc)
			h.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tt.status {
				t.Errorf("Want status '%d', got '%d'", tt.status, responseRecorder.Code)
			}
		})
	}
}

type MockStorage struct {
}

func (ms MockStorage) Get(ctx context.Context, queries []common.PropertyQuery) []common.PropertyResponse {
	return []common.PropertyResponse{}
}

func (ms MockStorage) Metrics() *ristretto.Metrics {
	return nil
}
