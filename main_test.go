package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	if port != "8000" {
		log.Fatalf("expected \"8000\" got %v", port)
		os.Exit(1)
	}

	if root != "." {
		log.Fatalf("expected \".\" got %v", root)
	}

	os.Exit(m.Run())
}

func TestBadPort(t *testing.T) {
	p1 := checkPort("4")

	if p1 != defaultport {
		t.Errorf("got %v want %v", p1, defaultport)
	}

	p2 := checkPort("8000")

	if p2 != "8000" {
		t.Errorf("got %v want 8000", p2)
	}
}

func TestFileHandler(t *testing.T) {
	handler := fileHandler(".")
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got: %v want 200", w.Code)
	}
}

func TestGZWithoutHeader(t *testing.T) {
	handler := fileHandler(".")
	req, _ := http.NewRequest("GET", "/test.js", nil)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got: %v want 200", w.Code)
	}

	if _, ok := w.HeaderMap["Content-Encoding"]; ok != false {
		t.Errorf("got \"Content-Encoding\"")
	}
}

func TestFavicon(t *testing.T) {
	handler := noIconHandler(fileHandler("."))
	req, _ := http.NewRequest("GET", "/favicon.ico", nil)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got: %v want 200", w.Code)
	}

}
