package main

import (
	// "fmt"

	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	// "github.com/gorilla/websocket"
	// "sync"
)

func TestMain(m *testing.M) {

	// check defaults
	if port != "3000" {
		log.Fatalf("expected \"3000\" got %v", port)
		os.Exit(1)
	}
	// if root != "." {
	// 	m.Errorf("got %v want \".\"", root)
	// }

	os.Exit(m.Run())
}

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

func TestGZWithHeader(t *testing.T) {
	handler := gzHandler(fileHandler("."))
	req, _ := http.NewRequest("GET", "/test.js", nil)
	req.Header["Accept-Encoding"] = []string{"gzip"}

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got: %v want 200", w.Code)
	}

	if h, ok := w.HeaderMap["Content-Encoding"]; ok != true {
		t.Errorf("got: %v want \"Content-Encoding\"", ok)
	} else if h[0] != "gzip" {
		t.Errorf("got: %v want \"gzip\"", h[0])

	}
}

func TestGZWithoutHeader(t *testing.T) {
	handler := gzHandler(fileHandler("."))
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
