package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
)

func userinfo(w http.ResponseWriter, r *http.Request) {
	/* Find matching user document */
	uuid := r.URL.Query().Get("uuid")
	log.Println("Received request: ", r.URL)

	var user User
	var userinfo UserInfo
	err := collection.FindOne(context.TODO(), bson.D{{Key: "uuid", Value: uuid}}).Decode(&user)

	/* Check if user was found */
	if err != nil {
		userinfo.Exists = false
		userinfo.IsAdmin = false
	} else {
		userinfo.Exists = true
		userinfo.IsAdmin = slices.Contains(user.Tags, "admin")
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

		fmt.Fprintf(w, string(decodedUserinfo))
	}
}
