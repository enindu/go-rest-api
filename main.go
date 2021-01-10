package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/forgot-password", forgotPasswordHandler)

	log.Println(fmt.Sprintf("Server listens to port %s", port))

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
