package main

import (
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

func validateServer(request *http.Request, method string) *Response {
	if request.Method != method {
		return &Response{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    "Method not allowed",
		}
	}

	if request.Header.Get("Content-Type") != "application/json" {
		return &Response{
			StatusCode: http.StatusExpectationFailed,
			Message:    "Expectation failed",
		}
	}

	return &Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
}

func validateData(data interface{}) *Response {
	validation := validator.New().Struct(data)
	if validation != nil {
		return &Response{
			StatusCode: http.StatusBadRequest,
			Message:    "One or more fields invalid",
		}
	}

	return &Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
}
