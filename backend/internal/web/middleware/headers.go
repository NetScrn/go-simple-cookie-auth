package middleware

import "net/http"

func CommonHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "application/json;charset=utf-8")
		header.Set("Accept", "application/json")
		header.Set("X-Content-Type-Options", "nosniff")
		header.Set("X-Frame-Options", "DENY")
		header.Set("X-XSS-Protection", "0")
		header.Set("Cache-Control", "no-store")
		header.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; sandbox")

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		h.ServeHTTP(w, r)
	})
}
