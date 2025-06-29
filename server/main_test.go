package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoutes(t *testing.T) {
	test := []struct {
		name     string
		env      string
		path     string
		wantCode int
		wantBody string
	}{
		{
			name:     "ping",
			env:      "development",
			path:     "/api/ping",
			wantCode: http.StatusOK,
			wantBody: "pong",
		},
		{
			name:     "root in dev",
			env:      "development",
			path:     "/",
			wantCode: http.StatusOK,
			wantBody: "Confido is running in development mode on port 8080\n",
		},
		{
			name:     "unknown route",
			env:      "development",
			path:     "/doest-not-exist",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			handler := buildRouter(tc.env, "8080")

			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Fatalf("status: got %d, want %d", rec.Code, tc.wantCode)
			}

			if tc.wantBody != "" && !strings.Contains(rec.Body.String(), tc.wantBody) {
				t.Fatalf("body mismatch: got %q, want to contain %q", rec.Body.String(), tc.wantBody)
			}

		})
	}
}
