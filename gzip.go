package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gw := gzResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gw, r)

	})
}
