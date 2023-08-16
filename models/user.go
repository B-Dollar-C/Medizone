package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	AuthToken        string             `json: "authtoken,omitempty" bson: "authtoken,omitempty"`
	FirstName        string             `json: "firstname,omitempty" bson: "firstname,omitempty"`
	LastName         string             `json: "lastname,omitempty" bson: "lastname,omitempty"`
	Email            string             `json: "email" bson: "email"`
	Password         string             `json: "-" bson: "-"`
	MedicalStoreName string             `json: "medicalstorename,omitempty" bson: "medicalstorename,omitempty"`
	Country          string             `json: "country,omitempty" bson: "country,omitempty"`
	StreetAddress_v1 string             `json: "streetaddress_v1,omitempty" bson: "streetaddress_v1,omitempty"`
	StreetAddress_v2 string             `json: "streetaddress_v2,omitempty" bson: "streetaddress_v2,omitempty"`
	Town             string             `json: "town,omitempty" bson: "town,omitempty"`
	State            string             `json: "state,omitempty" bson: "state,omitempty"`
	PostCode         string             `json: "postcode,omitempty" bson: "postcode,omitempty"`
	Phone            string             `json: "phone,omitempty" bson: "phone,omitempty"`
	ID               primitive.ObjectID `json: "id" bson: "_id"`
	Role             string             `json: "role,omitempty" bson: "role,omitempty"`
	CartItems        []string           `json: "cartitems" bson: "cartitems"`
	CreatedAt        time.Time          `json: "createdat,omitempty" bson: "createdat,omitempty"`
	UpdatedAt        time.Time          `json: "updatedat,omitempty" bson: "updatedat,omitempty"`
	DeletedAt        time.Time          `json: "deletedat,omitempty" bson: "deletedat,omitempty"`
	IsDeleted        bool               `json: "isdeleted,omitempty" bson: "isdeleted,omitempty"`
}
