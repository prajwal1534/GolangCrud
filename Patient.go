package main

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Patient struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Age       int           `json:"age" bson:"age"`
	Gender    string        `json:"gender" bson:"gender"`
	Phone     string        `json:"phone" bson:"phone"`
	Email     string        `json:"email" bson:"email"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

type CreatePatientRequest struct {
	Name   string `json:"name" binding:"required"`
	Age    *int   `json:"age" binding:"required,gte=0,lte=150"`
	Gender string `json:"gender" binding:"omitempty,oneof=male female other"`
	Phone  string `json:"phone" binding:"required"`
	Email  string `json:"email" binding:"omitempty,email"`
}

type UpdatePatientRequest struct {
	Name   *string `json:"name" binding:"omitempty,min=1"`
	Age    *int    `json:"age" binding:"omitempty,gte=0,lte=150"`
	Gender *string `json:"gender" binding:"omitempty,oneof=male female other"`
	Phone  *string `json:"phone" binding:"omitempty,min=1"`
	Email  *string `json:"email" binding:"omitempty,email"`
}
