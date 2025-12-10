package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string {
	return e.Msg
}

type NotFoundError struct {
	Msg string
}

func (e NotFoundError) Error() string {
	return e.Msg
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func WithErrorHandling(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			log.Printf("[ERROR] Request %s %s failed: %v", r.Method, r.URL.Path, err)
			var statusCode int
			var clientMsg string

			switch e := err.(type) {
			case ValidationError:
				statusCode = http.StatusBadRequest
				clientMsg = e.Msg
			case NotFoundError:
				statusCode = http.StatusNotFound
				clientMsg = e.Msg
			default:
				statusCode = http.StatusInternalServerError
				clientMsg = "Internal server error"
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(APIResponse{
				Success: false,
				Error:   clientMsg,
			})
			return
		}
	}
}