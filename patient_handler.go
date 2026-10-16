package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PatientHandler struct {
	collection *mongo.Collection
}

func NewPatientHandler(col *mongo.Collection) *PatientHandler {
	return &PatientHandler{
		collection: col,
	}
}

func (h *PatientHandler) Create(c *gin.Context) {
	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the patient into the database
	now := time.Now().UTC()

	patient := Patient{
		Name:      req.Name,
		Age:       *req.Age,
		Gender:    req.Gender,
		Phone:     req.Phone,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	result, err := h.collection.InsertOne(ctx, patient)

	if err != nil {
		log.Printf("Insert patient failed %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create patient"})
		return

	}

	patient.ID = result.InsertedID.(bson.ObjectID)
	c.JSON(http.StatusCreated, patient)

}

func (h *PatientHandler) GetByID(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid iD"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var patient Patient
	err = h.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&patient)

	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	if err != nil {
		log.Printf("Cannot retrieve the patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get the patient details"})
		return
	}

	c.JSON(http.StatusOK, patient)

}

func (h *PatientHandler) List(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetLimit(50).SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{}

	phone := c.Query("phone")
	if phone != "" {
		filter["phone"] = phone
	}
	cursor, err := h.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("list patients: find failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list patients"})
		return
	}

	defer cursor.Close(ctx)

	patients := []Patient{}

	err = cursor.All(ctx, &patients)

	if err != nil {
		log.Printf("list patients: decode failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list patients"})
		return
	}

	c.JSON(http.StatusOK, patients)

}

func (h *PatientHandler) Update(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(c.Param(("id")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	set := bson.M{}

	if req.Name != nil {
		set["name"] = *req.Name
	}

	if req.Age != nil {
		set["age"] = *req.Age
	}

	if req.Gender != nil {
		set["gender"] = *req.Gender
	}

	if req.Phone != nil {
		set["phone"] = *req.Phone
	}

	if req.Email != nil {
		set["email"] = *req.Email
	}

	if len(set) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}
	set["updated_at"] = time.Now().UTC()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var patient Patient
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = h.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": set}, opts).Decode(&patient)

	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	if err != nil {
		log.Printf("Cannot update the patient details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update the patient details"})
		return
	}

	c.JSON(http.StatusOK, patient)

}

func (h *PatientHandler) Delete(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Printf("object id invalid: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	result, err := h.collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		log.Printf("Cannot delete the patient details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete the patient details"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	c.Status(http.StatusNoContent)

}
