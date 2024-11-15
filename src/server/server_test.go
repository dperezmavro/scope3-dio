package server

import (
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	tests := []struct {
		path    string
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
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {

		})
	}
}
