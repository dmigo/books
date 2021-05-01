package books

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title" binding:"required"`
	Author      string             `json:"author" bson:"author" binding:"required"`
	PublishYear int                `json:"publish_year" bson:"publish_year" binding:"required"`
	Status      Status             `json:"status" bson:"status"`
}

type Status int

const (
	WANT_TO_READ Status = 0
	READING_NOW  Status = 1
	ALREADY_READ Status = 2
)

type NewStatus struct {
	Status int `json:"status" bson:"status" binding:"required"`
}
