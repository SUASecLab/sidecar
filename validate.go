package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kataras/jwt"
)

func validate(w http.ResponseWriter, r *http.Request) {
	/* Get token */
	token := r.URL.Query().Get("token")
	result := ValidationResult{"", false}
	log.Println("Received request: ", r.URL)

	/* Check if token was provided, otherwise return an error message */
	if len(token) > 0 {
		/* Check if provided token is valid, otherwise return an error message */
		_, err := jwt.Verify(jwt.HS256, []byte(jwtKey), []byte(token))

		if err != nil {
			errorMessage := "Failed to validate handed over JWT token"
			log.Println(errorMessage)
			result.Error = errorMessage
		} else {
			result.Valid = true
		}
	} else {
		errorMessage := "No token provided"
		http.Error(w, errorMessage, http.StatusBadRequest)
		result.Error = errorMessage
	}

	/* Decode and send result */
	decodedResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Could not decode result", http.StatusInternalServerError)
		log.Println("Could not decode result: ", err)
	} else {
		/* Set result headers */
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		/* Print result */
		fmt.Fprintf(w, string(decodedResult))
	}
}
