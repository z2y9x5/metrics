package handler

import (
	"net/http"
	"strings"
)

// Удалить повторяющиеся слеши из пути.
func FixDoubleSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := strings.ReplaceAll(r.URL.Path, "//", "/")
		r.URL.Path = cleanPath
		next.ServeHTTP(w, r)
	})
}
