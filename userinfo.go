package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"slices"

	"github.com/kataras/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

func userinfo(w http.ResponseWriter, r *http.Request) {
	/* Find matching user document */
	uuid := r.URL.Query().Get("uuid")
	token := r.URL.Query().Get("token")

	// escape input
	uuid = html.EscapeString(uuid)
	token = html.EscapeString(token)

	log.Println("Received request: ", r.URL)

	// if token is set, get uuid from token
	if len(token) > 0 {
		verifiedToken, err := jwt.Verify(jwt.HS256, []byte(jwtKey), []byte(token))
		if err != nil {
			log.Println(err)
			uuid = "" // this will cause the user not to be found, returning empty data
		} else {
			var claims RulesCustomClaims
			err := verifiedToken.Claims(&claims)

			if err != nil {
				log.Println(err)
				uuid = "" // see comment above
			} else {
				uuid = claims.UUID
			}
		}
	} else {
		http.Error(w, "No token provided", http.StatusBadRequest)
	}

	var user User
	var userinfo UserInfo
	err := collection.FindOne(context.TODO(), bson.D{{Key: "uuid", Value: uuid}}).Decode(&user)

	/* Check if user was found */
	if err != nil {
		userinfo.Exists = false
		userinfo.IsAdmin = false
		userinfo.UUID = ""
	} else {
		userinfo.Exists = true
		userinfo.IsAdmin = slices.Contains(user.Tags, "admin")
		userinfo.UUID = user.Uuid
	}

	/* Convert userinfo to JSON */
	decodedUserinfo, err := json.Marshal(userinfo)
	if err != nil {
		http.Error(w, "Could not decode user info", http.StatusInternalServerError)
		log.Println("Could not decode user info: ", err)
	} else {
		/* Set result headers */
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, string(decodedUserinfo))
	}
}
