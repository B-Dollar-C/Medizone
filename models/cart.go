package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID          primitive.ObjectID `json: "id" bson: "_id"`
	AuthToken   string             `json: "authtoken,omitempty" bson: "authtoken,omitempty"`
	Medicine    Medicine           `json: "medicine" bson: "medicine"`
	CreatedAt   time.Time          `json: "createdat,omitempty" bson: "createdat,omitempty"`
	UpdatedAt   time.Time          `json: "updatedat,omitempty" bson: "updatedat,omitempty"`
	DeletedAt   time.Time          `json: "deletedat,omitempty" bson: "deletedat,omitempty"`
	IsDeleted   bool               `json: "isdeleted" bson: "isdeleted"`
	UserID      string             `json: "userid" bson: "userid"`
	Quantity    int                `json: "quantity,omitempty" bson: "quantity,omitempty"`
	Total       float32            `json: "total,omitempty" bson: "total,omitempty"`
	OrderPlaced bool               `json: "orderplaced" bson: "orderplaced"`
}
