package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/northern-ai/medizone/controllers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	route := httprouter.New()
	uc := controllers.NewUserController(getClient())
	route.POST("/api/SignUp", uc.SignUp)
	route.POST("/api/SignIn", uc.SignIn)
	route.POST("/api/CreateOwner", uc.CreateOwner)
	route.POST("/api/CreateMedicine", uc.CreateMedicine)
	route.GET("/api/GetMedicines", uc.GetMedicines)
	route.POST("/api/AddToCart", uc.AddToCart)
	route.PUT("/api/RemoveFromCart", uc.RemoveFromCart)
	route.POST("/api/Checkout", uc.Checkout)
	route.PUT("/api/PlaceOrder", uc.PlaceOrder)
	route.GET("/api/GetCarts/", uc.GetCarts)
	fmt.Println("Starting server at port 8080")
	http.ListenAndServe(":8080", route)

}

func getClient() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}

	return client

}
