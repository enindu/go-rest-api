package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/rs/xid"
)

// RegisterRequest struct
type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email,max=191"`
	Username string `json:"username" validate:"required,alpha,min=3,max=6"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

func registerHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Accept", "application/json")
	response.Header().Set("Content-Type", "application/json")

	serverValidation := validateServer(request, "PUT")
	if serverValidation.StatusCode != 200 {
		response.WriteHeader(serverValidation.StatusCode)
		json.NewEncoder(response).Encode(serverValidation)
		return
	}

	registerRequest := &RegisterRequest{}
	json.NewDecoder(request.Body).Decode(registerRequest)

	dataValidation := validateData(registerRequest)
	if dataValidation.StatusCode != 200 {
		response.WriteHeader(dataValidation.StatusCode)
		json.NewEncoder(response).Encode(dataValidation)
		return
	}

	database, issue := connectDatabase()
	if issue != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(&Response{
			StatusCode: http.StatusInternalServerError,
			Message:    issue.Error(),
		})
		return
	}

	account := &Account{}
	accountRead := database.Where("email = @Email or username = @Username", registerRequest).First(account)
	if accountRead.Error == nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(&Response{
			StatusCode: http.StatusBadRequest,
			Message:    "There is an account already using that email or username",
		})
		return
	}

	password, issue := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), 10)
	if issue != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(&Response{
			StatusCode: http.StatusInternalServerError,
			Message:    issue.Error(),
		})
		return
	}

	uniqueID := xid.New().String()
	database.Create(&Account{
		UniqueID: uniqueID,
		Email:    registerRequest.Email,
		Username: registerRequest.Username,
		Password: string(password),
	})

	json.NewEncoder(response).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
		Data:       uniqueID,
	})
}
