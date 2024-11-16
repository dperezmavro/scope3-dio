package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scope3-dio/src/config"
)

func TestRoutes(t *testing.T) {
	tests := []struct {
		path    string
		method  string
		headers map[string]string
		status  int
	}{
		{
			path:   "/noexist",
			status: http.StatusNotFound,
		},
		{
			path:   "/healthcheck",
			status: http.StatusOK,
		},
		// {
		// 	path:   "/v2/measure",
		// 	status: http.StatusUnauthorized,
		// },
		// {
		// 	path:   "/v2/measure",
		// 	status: http.StatusBadRequest,
		// 	headers: map[string]string{
		// 		common.HeaderAuthorization: "dummy",
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			// t.Parallel()
			request := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.headers != nil {
				for k, v := range tt.headers {
					request.Header.Add(k, v)
				}
			}
			responseRecorder := httptest.NewRecorder()

			h := CreateRouter(config.Default, nil)
			h.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tt.status {
				t.Errorf("Want status '%d', got '%d'", tt.status, responseRecorder.Code)
			}
		})
	}
}
