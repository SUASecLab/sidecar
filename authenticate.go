package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kataras/jwt"
)

func authenticate(w http.ResponseWriter, r *http.Request) {
	/* Get parameters from URL */
	tokenString := r.URL.Query().Get("token")
	service := r.URL.Query().Get("service")
	result := AuthenticationResponse{false}

	log.Println("Received request: ", r.URL)

	/* Check if token and service are present */
	if tokenString != "" && service != "" {
		/* Verify token */
		verifiedToken, err := jwt.Verify(jwt.HS256, []byte(jwtKey), []byte(tokenString))

		if err != nil {
			log.Println(err)
		} else {
			var claims RulesCustomClaims
			err := verifiedToken.Claims(&claims)

			if err != nil {
				log.Println(err)
			} else {
				result.Allowed = decide(service, claims)
			}
		}
	} else {
		http.Error(w, "Token and service were not provided", http.StatusBadRequest)
	}

	resultText, err := json.Marshal(result)
	if err != nil {
		errorMessage := "Error while encoding json data"
		log.Println(errorMessage, err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
	} else {
		/* Set result headers */
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		/* Show result */
		fmt.Fprint(w, string(resultText))
	}
}
