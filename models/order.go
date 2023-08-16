package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID           primitive.ObjectID `json: "id" bson: "_id"`
	AuthToken    string             `json: "authtoken,omitempty" bson: "authtoken,omitempty"`
	User         User               `json: "user" bson: "user"`
	CartItems    []Cart             `json: "cartitems" bson: "cartitems"`
	CreatedAt    time.Time          `json: "createdat" bson: "createdat"`
	UpdatedAt    time.Time          `json: "updatedat" bson: "updatedat"`
	DeletedAt    time.Time          `json: "deletedat,omitempty" bson: "deletedat,omitempty"`
	IsDeleted    bool               `json: "isdeleted" bson: "isdeleted"`
	Total        float32            `json: "total,omitempty" bson: "total,omitempty"`
	SubTotal     float32            `json: "subtotal,omitempty" bson: "subtotal,omitempty"`
	Tax          string             `json: "tax,omitempty" bson: "tax,omitempty"`
	ShippingCost float32            `json: "shippingcost,omitempty" bson: "shippingcost,omitempty"`
	PaymentType  string             `json: "paymenttype,omitempty" bson: "paymenttype,omitempty"`
}
