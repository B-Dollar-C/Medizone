package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Medicine struct {
	ID         primitive.ObjectID `json: "id" bson: "_id"`
	AuthToken  string             `json: "authtoken,omitempty" bson: "authtoken,omitempty"`
	Name       string             `json: "name,omitempty" bson: "name,omitempty"`
	Price      float32            `json: "price,omitempty" bson: "price,omitempty"`
	ImageUrl   string             `json: "imageurl,omitempty" bson: "imageurl,omitempty"`
	CreatedAt  time.Time          `json: "createdat" bson: "createdat"`
	UpdatedAt  time.Time          `json: "updatedat" bson: "updatedat"`
	DeletedAt  time.Time          `json: "deletedat,omitempty" bson: "deletedat,omitempty"`
	IsDeleted  bool               `json: "isdeleted" bson: "isdeleted"`
	MedOwnerID string             `json: "medownerid,omitempty" bson: "medownerid,omitempty"`
}
