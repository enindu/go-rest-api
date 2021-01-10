package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/rs/xid"
)

// ForgotPasswordRequest struct
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=191"`
}

func forgotPasswordHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Accept", "application/json")
	response.Header().Set("Content-Type", "application/json")

	serverValidation := validateServer(request, "PUT")
	if serverValidation.StatusCode != 200 {
		response.WriteHeader(serverValidation.StatusCode)
		json.NewEncoder(response).Encode(serverValidation)
		return
	}

	forgotPasswordRequest := &ForgotPasswordRequest{}
	json.NewDecoder(request.Body).Decode(forgotPasswordRequest)

	dataValidation := validateData(forgotPasswordRequest)
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
	accountRead := database.Where("email = @Email", forgotPasswordRequest).First(account)
	if accountRead.Error != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(&Response{
			StatusCode: http.StatusBadRequest,
			Message:    "There is no account found using that email",
		})
		return
	}

	token := xid.New().String()
	mail := smtp.PlainAuth("PLAIN", "f204e4687aa3cb", "3f23587357e4e7", "smtp.mailtrap.io")
	from := "noreply@example.com"
	to := forgotPasswordRequest.Email
	subject := "Password Reset"
	message := []byte("From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		"https://www.example.com/password-reset?username=" + account.Username + "&token=" + token)
	if smtp.SendMail("smtp.mailtrap.io:2525", mail, from, []string{to}, message) != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(&Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong while sending password reset email",
		})
		return
	}

	account.Token = token
	accountRead.Save(account)

	json.NewEncoder(response).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
		Data:       fmt.Sprintf("%s,%s", account.Username, token),
	})
}
