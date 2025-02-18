package middleware

import (
	"context"
	"fmt"
	"net/http"
)

func Cors(ctx context.Context, url string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("r : %+v\n", r)
		enableCors(&w, url)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func enableCors(w *http.ResponseWriter, url string) {
	(*w).Header().Set("Access-Control-Allow-Origin", url)
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	// Allow specified headers
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
}
