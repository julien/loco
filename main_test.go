package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileHandler(t *testing.T) {
	handler := fileHandler(".")
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got: %v", w.Code, "wanted 200")
	}
}

func TestScriptHandler(t *testing.T) {
	handler := scriptHandler()

	req, _ := http.NewRequest("GET", "/livereload.js", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	ct := w.Header()["Content-Type"]

	if w.Code != http.StatusOK {
		t.Errorf("got: %v want 200", w.Code)
	}

	if len(ct) < 1 || ct[0] != "application/javascript" {
		t.Errorf("got %v", ct[0], "want application/javascript")
	}
}

func TestSocketHandler(t *testing.T) {
	handler := socketHandler()

	req, _ := http.NewRequest("GET", "/ws", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got: %v want 200", w.Code)
	}
}
