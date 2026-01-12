package api_test

import (
	"backend-example/src/api"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	api.HandleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if body := rr.Body.String(); body != "OK" {
		t.Fatalf(`expected body "OK", got %q`, body)
	}
}

func TestHandleSum_ValidParameters(t *testing.T) {
	cases := []struct {
		x, y string
		sum  int
	}{
		{"2", "3", 5},
		{"-1", "4", 3},
		{"0", "0", 0},
		{"10", "20", 30},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("x=%s_y=%s", tc.x, tc.y), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sum?x=%s&y=%s", tc.x, tc.y), nil)
			rr := httptest.NewRecorder()

			api.HandleSum(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
			}
			if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
				t.Fatalf("expected Content-Type application/json, got %q", ct)
			}
			expected := fmt.Sprintf(`{"sum": %d}`, tc.sum)
			if body := strings.TrimSpace(rr.Body.String()); body != expected {
				t.Fatalf("expected body %q, got %q", expected, body)
			}
		})
	}
}

func TestHandleSum_InvalidParameters(t *testing.T) {
	cases := []struct {
		name string
		url  string
	}{
		{"missing_x", "/sum?y=1"},
		{"missing_y", "/sum?x=1"},
		{"missing_both", "/sum"},
		{"non_numeric_x", "/sum?x=foo&y=1"},
		{"non_numeric_y", "/sum?x=1&y=bar"},
		{"blank_x", "/sum?x=%20&y=1"},
		{"blank_y", "/sum?x=1&y=%20"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			rr := httptest.NewRecorder()

			api.HandleSum(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
			}
			if body := strings.TrimSpace(rr.Body.String()); body != "Invalid parameters" {
				t.Fatalf("expected body %q, got %q", "Invalid parameters", body)
			}
			if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") && ct != "" {
				t.Fatalf("expected Content-Type to be text/plain or empty, got %q", ct)
			}
		})
	}
}
