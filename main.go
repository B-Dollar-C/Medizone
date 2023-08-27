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

	corsMiddleware := func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r, ps)
		}
	}

	route.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	route.POST("/api/SignUp", corsMiddleware(uc.SignUp))
	route.POST("/api/SignIn", corsMiddleware(uc.SignIn))
	route.POST("/api/CreateOwner", corsMiddleware(uc.CreateOwner))
	route.POST("/api/CreateMedicine", corsMiddleware(uc.CreateMedicine))
	route.GET("/api/GetMedicines", corsMiddleware(uc.GetMedicines))
	route.POST("/api/AddToCart", corsMiddleware(uc.AddToCart))
	route.PUT("/api/RemoveFromCart", corsMiddleware(uc.RemoveFromCart))
	route.POST("/api/Checkout", corsMiddleware(uc.Checkout))
	route.PUT("/api/PlaceOrder", corsMiddleware(uc.PlaceOrder))
	route.GET("/api/GetCarts/", corsMiddleware(uc.GetCarts))
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
