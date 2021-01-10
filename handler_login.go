package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest struct
type LoginRequest struct {
	Username string `json:"username" validate:"required,alpha,min=3,max=6"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Accept", "application/json")
	response.Header().Set("Content-Type", "application/json")

	serverValidation := validateServer(request, "POST")
	if serverValidation.StatusCode != 200 {
		response.WriteHeader(serverValidation.StatusCode)
		json.NewEncoder(response).Encode(serverValidation)
		return
	}

	loginRequest := &LoginRequest{}
	json.NewDecoder(request.Body).Decode(loginRequest)

	dataValidation := validateData(loginRequest)
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
	accountRead := database.Where("username = @Username", loginRequest).First(account)
	if accountRead.Error != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(&Response{
			StatusCode: http.StatusBadRequest,
			Message:    "There is no account found using that username",
		})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(loginRequest.Password)) != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(&Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Password is invalid",
		})
		return
	}

	json.NewEncoder(response).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
		Data:       account.UniqueID,
	})
}
