package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// create a gzipped response
			wrw := NewWrappedResponseWriter(w)
			wrw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(wrw, r)
			defer wrw.Flush()

			return
		}
		// handle normal response
		next.ServeHTTP(w, r)
	})
}

type WrappedResponseWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(w)
	return &WrappedResponseWriter{w: w, gw: gw}
}

func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.w.Header()
}

func (wr *WrappedResponseWriter) Write(b []byte) (int, error) {
	return wr.gw.Write(b)
}

func (wr *WrappedResponseWriter) WriteHeader(statusCode int) {
	wr.w.WriteHeader(statusCode)
}

func (wr *WrappedResponseWriter) Flush() {
	wr.gw.Flush()
	wr.gw.Close()
}
