package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/kataras/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

func issuance(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	name := r.URL.Query().Get("name")
	token := r.URL.Query().Get("token")

	// escape input
	uuid = html.EscapeString(uuid)
	name = html.EscapeString(name)
	token = html.EscapeString(token)

	currentTime := time.Now()
	result := IssuanceResult{"", ""}

	log.Println("Received request: ", r.URL)

	if uuid != "" {
		var user User
		err := collection.FindOne(context.TODO(), bson.D{{Key: "uuid", Value: uuid}}).Decode(&user)
		if err != nil {
			errorMessage := "User does not exist"
			http.Error(w, errorMessage, http.StatusForbidden)
			result.Error = errorMessage
			log.Println(errorMessage, err)
		} else {
			signedToken, err := jwt.Sign(jwt.HS256, []byte(jwtKey), map[string]interface{}{
				"uuid":      user.Uuid,
				"tags":      user.Tags,
				"moderator": slices.Contains(user.Tags, "admin"),
				"exp":       currentTime.Unix() + 60*60*24,
				"nbf":       currentTime.Unix() - 10,
				"iat":       currentTime.Unix(),
			})

			if err != nil {
				errorMessage := "Could not sign authentication token"
				log.Println(errorMessage, err)
				http.Error(w, errorMessage, http.StatusInternalServerError)
				result.Error = errorMessage
			} else {
				result.Token = string(signedToken)
			}
		}
	} else if name != "" && token != "" {
		verifiedToken, err := jwt.Verify(jwt.HS256, []byte(jwtKey), []byte(token))

		if err != nil {
			errorMessage := "Could not verify provided token"
			http.Error(w, errorMessage, http.StatusBadRequest)
			log.Println(errorMessage, err)
		} else {
			var claims RulesCustomClaims
			err = verifiedToken.Claims(&claims)

			if err != nil {
				errorMessage := "Could not fetch claims"
				http.Error(w, errorMessage, http.StatusInternalServerError)
				log.Println(errorMessage, err)
			} else {
				isModerator := decide("jitsiModerator", claims)

				signedToken, err := jwt.Sign(jwt.HS256, []byte(jitsiKey), map[string]interface{}{
					"context": map[string]interface{}{
						"user": map[string]interface{}{
							"id":   jitsiId,
							"name": name,
						},
					},
					"nbf":       currentTime.Unix() - 10,
					"aud":       "jitsi",
					"iss":       jitsiIssuer,
					"room":      "*",
					"moderator": isModerator,
					"sub":       jitsiUrl,
					"iat":       currentTime.Unix(),
					"exp":       currentTime.Unix() + 60,
				})

				jitsiId += 1

				if err != nil {
					errorMessage := "Could not sign new token"
					log.Println(errorMessage, err)
					http.Error(w, errorMessage, http.StatusInternalServerError)
				} else {
					if slices.Contains(claims.Tags, "exam") {
						signedToken = []byte("")
					}

					/* Set result headers */
					result.Token = string(signedToken)
				}
			}
		}
	} else {
		errorMessage := "No 'uuid || name && token' provided"
		http.Error(w, errorMessage, http.StatusBadRequest)
		result.Error = errorMessage
	}

	/* Encode token */
	jsonToken, err := json.Marshal(result)
	if err != nil {
		errorMessage := "Could not encode token"
		http.Error(w, errorMessage, http.StatusInternalServerError)
		log.Println(errorMessage, err)
	} else {
		/* Return issued token */
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, string(jsonToken))
	}
}
