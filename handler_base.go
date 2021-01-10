package main

import (
	"encoding/json"
	"net/http"
)

func baseHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	if request.RequestURI != "/" {
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(&Response{
			StatusCode: http.StatusNotFound,
			Message:    "Not found",
		})
		return
	}

	if request.Method != "GET" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(response).Encode(&Response{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    "Method not allowed",
		})
		return
	}

	json.NewEncoder(response).Encode("This is a dummy API")
}
