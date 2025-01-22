package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	handler := &myHandler{}
	mux := http.NewServeMux()
	mux.Handle("/path", handler)
	srv := &http.Server{
		Addr:    net.JoinHostPort("", "8080"),
		Handler: clientDisconnectedErrormiddleware(mux),
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("server not started caused by: %v", err)
	}
}

type myHandler struct{}

type response struct {
	Status string `json:"status"`
}

var _ http.Handler = (*myHandler)(nil)

func (s *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	res := &response{"ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func clientDisconnectedErrormiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		done := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err == context.Canceled {
					fmt.Println("client canceld")
					http.Error(w, "client canceled", 444)
				}
			case <-done:
				// respond completely.
			}
		}()
		defer close(done)
		next.ServeHTTP(w, r)
	})
}
