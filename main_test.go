package main

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	// "github.com/gorilla/websocket"
	// "sync"
)

// var once sync.Once

func TestBadPort(t *testing.T) {
	p1 := checkPort("4")

	if p1 != defaultport {
		t.Errorf("got %v want %v", p1, defaultport)
	}

	p2 := checkPort("3000")

	if p2 != "3000" {
		t.Errorf("got %v want 3000", p2)
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
		t.Errorf("got %v want application/javascript", ct[0])
	}
}

func TestSocketHandlerPOST(t *testing.T) {
	handler := socketHandler()

	req, _ := http.NewRequest("POST", "/ws", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Errorf("got: %v wanted an error", w.Code)
	}
}

func TestSocketHandlerUpgrader(t *testing.T) {
	handler := socketHandler()

	req, _ := http.NewRequest("GET", "/ws", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Errorf("got: %v wanted an error", w.Code)
	}
}
