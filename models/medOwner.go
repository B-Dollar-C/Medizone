package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MedOwner struct {
	ID         primitive.ObjectID `json: "id" bson: "_id"`
	AuthToken  string             `json: "authtoken,omitempty" bson: "authtoken,omitempty"`
	OwnerName  string             `json: "ownername,omitempty" bson: "ownername,omitempty"`
	StoreName  string             `json: "storename" bson: "storename"`
	StoreEmail string             `json: "storeemail" bson: "storeemail"`
	Password   string             `json: "-" bson: "-"`
	Medicines  []Medicine         `json: "medicines" bson: "medicines"`
	Role       string             `json: "role,omitempty" bson: "role,omitempty"`
	ImageUrl   string             `json: "imageurl,omitempty" bson: "imageurl,omitempty"`
	CreatedAt  time.Time          `json: "createdat" bson: "createdat"`
	UpdatedAt  time.Time          `json: "updatedat" bson: "updatedat"`
	DeletedAt  time.Time          `json: "deletedat,omitempty" bson: "deletedat,omitempty"`
	IsDeleted  bool               `json: "isdeleted" bson: "isdeleted"`
}
