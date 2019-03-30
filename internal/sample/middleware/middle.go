package middleware

import (
	"fmt"
	"net/http"
)

// SampleOne middleware
func SampleOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(" SampleOne Middle ")
		next.ServeHTTP(w, r)
	})
}

// SampleTwo middleware
func SampleTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(" SampleTwo Middle ")
		next.ServeHTTP(w, r)
	})
}

// SampleFour middleware
func SampleFour(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(" SampleFour Middle ")
		next.ServeHTTP(w, r)
	})
}
