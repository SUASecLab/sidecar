package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	collection *mongo.Collection

	jwtKey      string
	jitsiIssuer string
	jitsiKey    string
)

func main() {
	/* Read environment variables */
	mongoUri := os.Getenv("MONGODB_URI")

	jwtKey = os.Getenv("JWT_KEY")
	jitsiIssuer = os.Getenv("JITSI_ISS")
	jitsiKey = os.Getenv("SECRET_JITSI_KEY")

	/* Check if all variables are set */
	if mongoUri == "" {
		log.Println("Database information incomplete")
		os.Exit(1)
	}

	if jwtKey == "" || jitsiIssuer == "" || jitsiKey == "" {
		log.Println("Jitsi information incomplete")
		os.Exit(1)
	}

	/* Connect to database */
	log.Println("Connecting to database...")
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	collection = mongoClient.Database("workadventure").Collection("users")

	log.Println("Connected to database, setting up API routes...")

	/* Set up router */
	r := mux.NewRouter()
	r.HandleFunc("/auth", authenticate).Methods("GET")
	r.HandleFunc("/issuance", issuance).Methods("GET")
	r.HandleFunc("/userinfo", userinfo).Methods("GET")
	r.HandleFunc("/validate", validate).Methods("GET")

	log.Println("Starting up sidecar")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("Sidecar failed: ", err)
	}
}
