package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/kataras/jwt"
)

var lastJWTAuthentication time.Time

func authenticate(w http.ResponseWriter, r *http.Request) {
	/* Get parameters from URL */
	var tokenString string
	var service string

	// Figure out if this is an authentication request from Traefik
	forwardedURI := r.Header.Get("X-Forwarded-Uri")
	if forwardedURI != "" {
		// Traefik request
		log.Println("Received forwarded URI request: ", forwardedURI)

		// Check if within time frame of 15 seconds (assume resources are loaded)
		if time.Since(lastJWTAuthentication).Seconds() < 15 {
			http.Error(w, "Success", 200)
			return
		}

		// Parse URI
		parsedURI, err := url.Parse(forwardedURI)
		if err != nil {
			// Can not parse URI -> Invalid
			log.Println("Can not parse URI: ", forwardedURI)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get token and service strings from X-Forwarded-Uri header
		tokenString = parsedURI.Query().Get("token")
		service = parsedURI.Query().Get("service")
	} else {
		// Get token and service strings from URL
		log.Println("Received request: ", r.URL)
		tokenString = r.URL.Query().Get("token")
		service = r.URL.Query().Get("service")
	}

	// Prepare response
	result := AuthenticationResponse{false}

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

				// Update time
				if result.Allowed && service == "noVNC" && forwardedURI != "" {
					lastJWTAuthentication = time.Now()
				}
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
