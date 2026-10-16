package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://127.0.0.1:27017"
	}

	client, err := connectMongo(uri)
	if err != nil {
		log.Fatalf("could not connect to MongoDB: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	defer client.Disconnect(context.Background())
	log.Printf("DB Connected Successfully")

	router := gin.Default()

	db := client.Database("doctorfirst")
	patientHandler := NewPatientHandler(db.Collection("patients"))
	healthHandler := NewHealthHandler(client)
	router.POST("/patients", patientHandler.Create)
	router.GET("/patients/:id", patientHandler.GetByID)
	router.GET("/patients", patientHandler.List)
	router.PATCH("/patients/:id", patientHandler.Update)
	router.DELETE("/patients/:id", patientHandler.Delete)
	router.GET("/health", healthHandler.Check)

	log.Fatal(router.Run(":" + port))

}
