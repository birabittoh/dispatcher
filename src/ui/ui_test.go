package ui_test

import (
	"backend-example/src/ui"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ui.HandleIndex)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedContentType := "text/html; charset=utf-8"
	if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
		t.Fatalf("handler returned wrong Content-Type: got %v want %v",
			ct, expectedContentType)
	}

	expectedBodySubstring := "<!DOCTYPE html>"
	if !strings.Contains(rr.Body.String(), expectedBodySubstring) {
		t.Fatalf("handler returned unexpected body: got %v want substring %v",
			rr.Body.String(), expectedBodySubstring)
	}
}
