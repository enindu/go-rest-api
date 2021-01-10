package main

import (
	"log"
	"net/http"

	"gorm.io/gorm"
)

// Response struct
type Response struct {
	StatusCode int    `json:"status-code"`
	Message    string `json:"message"`
	Data       string `json:"data"`
}

// Account struct
type Account struct {
	gorm.Model
	UniqueID string
	Token    string
	Email    string
	Username string
	Password string
}

func main() {
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/forgot-password", forgotPasswordHandler)

	log.Println("Server listens to port 8000")

	http.ListenAndServe(":8000", nil)
}
