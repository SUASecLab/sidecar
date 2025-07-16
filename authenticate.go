package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"

	"github.com/kataras/jwt"
)

func authenticate(w http.ResponseWriter, r *http.Request) {
	/* GET parameters from URL */
	var tokenString string
	var service string

	// Variables
	var parsedURI *url.URL
	var err error

	// Figure out if this is an authentication request from Traefik
	// If it is from Traefik, we recognize it by the X-Forwarded-Uri header being set
	forwardedURI := r.Header.Get("X-Forwarded-Uri")
	if forwardedURI != "" {
		// This is a Traefik request

		// Parse URI first
		parsedURI, err = url.Parse(forwardedURI)
		if err != nil {
			// Can not parse URI -> Invalid
			log.Println("Can not parse URI: ", forwardedURI)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set service (we can do this hardcoded here as the noVNC service is the only consumer of these statements)
		service = "noVNC"

		// Check if the VNC cookie is already present
		cookie, err := r.Cookie("vnc_token")

		// Cookie is not set -> Store the JWT token as VNC cookie and retry
		if err != nil {
			// Get token from parsed URI
			tokenString = parsedURI.Query().Get("token")

			// Set cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "vnc_token",
				Value:    tokenString,
				Path:     "/",
				MaxAge:   3600 * 12,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})

			// To be able to actually store the cookie we must redirect now
			// As we want to retry with the original request, go back to the original URL now
			// Reconstruct original URL from headers and noVNC subdir variable
			x_forwarded_proto := r.Header.Get("X-Forwarded-Proto")
			x_forwarded_host := r.Header.Get("X-Forwarded-Host")
			x_forwarded_uri := r.Header.Get("X-Forwarded-Uri")
			redirectTargetURL := fmt.Sprintf("%s://%s%s%s",
				x_forwarded_proto, x_forwarded_host, noVNCSubdir, x_forwarded_uri)

			// Send the redirect
			http.Redirect(w, r, redirectTargetURL, http.StatusFound)

			// The browser now redirects back and sets the cookie
			// By the redirection, the request will be sent once again
			// As the cookie is now set, this routine will be skipped then
			return
		}

		// Cookie is set -> hand over cookie value as JWT token value
		tokenString = cookie.Value
	} else {
		// Get token and service strings from URL
		log.Println("Received request: ", r.URL)
		tokenString = r.URL.Query().Get("token")
		service = r.URL.Query().Get("service")
	}

	// escape input
	tokenString = html.EscapeString(tokenString)
	service = html.EscapeString(service)

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
			}
		}
	} else {
		http.Error(w, "Token and/or service were not provided", http.StatusBadRequest)
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
