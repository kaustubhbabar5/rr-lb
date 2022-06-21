package server

import (
	"net/http"

	chttp "github.com/kaustubhbabar5/rr-lb/pkg/http"
)

func (s *HTTPServer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			chttp.JSON(w, http.StatusNotAcceptable, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
