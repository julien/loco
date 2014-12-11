package main

import (
	//"fmt"
	// "io/ioutil"
	// "log"
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
		t.Errorf("File server not working: %v", w.Code)
	}
	// body, err := ioutil.ReadAll(w.Body)
	// if err != nil {
	//     fmt.Println(err)
	// }
	// fmt.Println(string(body))
}

func TestScriptHandler(t *testing.T) {
	handler := scriptHandler()

	req, _ := http.NewRequest("GET", "/livereload.js", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	ct := w.Header()["Content-Type"]

	if w.Code != http.StatusOK {
		t.Errorf("got: %v", w.Code, "want 200")
	}

	if len(ct) < 1 || ct[0] != "application/javascript" {
		t.Errorf("got %v", ct[0], "want application/javascript")
	}
}
