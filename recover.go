package main

import (
	"log"
	"net/http"
)

func recoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\nPath: %s", err, r.URL.Path)
				http.Error(w, "Internal Server Error", 500)
			}
		}()
		next(w, r)
	}
}
