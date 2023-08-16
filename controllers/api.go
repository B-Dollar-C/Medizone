package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/northern-ai/medizone/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}

func (uc UserController) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" || user.Password == "" || user.FirstName == "" {
		response := map[string]interface{}{
			"code":    400,
			"status":  false,
			"message": "SignUp Details are required!",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}

	user.Password = string(hashedPassword)

	filter := bson.M{"email": user.Email}
	collection := uc.client.Database("medizone").Collection("users")
	var refUser models.User
	err = collection.FindOne(r.Context(), filter).Decode(&refUser)
	if err == nil {
		response := map[string]interface{}{
			"code":    400,
			"status":  false,
			"message": "Email Already In Use",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	parts := regexp.MustCompile(" ").Split(user.FirstName, 2)
	if len(parts) != 1 {
		user.FirstName = parts[0]
		user.LastName = parts[1]
	}

	result, err := collection.InsertOne(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authToken := (result.InsertedID.(primitive.ObjectID)).Hex()
	fil := bson.M{"_id": result.InsertedID.(primitive.ObjectID)}
	update := bson.M{"$set": bson.M{
		"authtoken": authToken,
		"role":      "customer",
		"cartitems": []models.Cart{},
		"createdat": time.Now(),
		"updatedat": time.Now(),
	}}

	_, red := collection.UpdateOne(r.Context(), fil, update)
	if red != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)

}

func (uc UserController) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	if user.Email == "" || user.Password == "" {
		response := map[string]interface{}{
			"code":    400,
			"status":  false,
			"message": "SignIn Details Required!",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	emailLogin := false
	passLogin := false
	responseUser := models.User{}
	collection := uc.client.Database("medizone").Collection("users")
	err := collection.FindOne(r.Context(), bson.M{"email": user.Email}).Decode(&responseUser)
	if err != nil {
		response := map[string]interface{}{
			"code":    400,
			"status":  false,
			"message": "Account Not Found With This Email",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return

	} else {
		emailLogin = true
	}

	if passErr := bcrypt.CompareHashAndPassword([]byte(responseUser.Password), []byte(user.Password)); passErr != nil {
		response := map[string]interface{}{
			"code":    400,
			"status":  false,
			"message": "Invalid Password",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return

	} else {
		passLogin = true
	}

	if emailLogin == true && passLogin == true {
		uj, err := json.Marshal(responseUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusFound)
		w.Write(uj)
	}

}

func (uc UserController) CreateOwner(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	owner := models.MedOwner{}
	json.NewDecoder(r.Body).Decode(&owner)
	if owner.StoreEmail == "" || owner.Password == "" || owner.StoreName == "" {
		response := map[string]interface{}{
			"code":    400,
			"status":  false,
			"message": "Owner Mandate Details are required!",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(owner.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}

	owner.Password = string(hashedPassword)

	filter := bson.M{"storeemail": owner.StoreEmail}
	collection := uc.client.Database("medizone").Collection("medowners")
	var refUser models.MedOwner
	err = collection.FindOne(r.Context(), filter).Decode(&refUser)
	if err == nil {
		response := map[string]interface{}{
			"code":    400,
			"status":  false,
			"message": "Email Already In Use",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	result, err := collection.InsertOne(r.Context(), owner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authToken := (result.InsertedID.(primitive.ObjectID)).Hex()
	fil := bson.M{"_id": result.InsertedID.(primitive.ObjectID)}
	update := bson.M{"$set": bson.M{
		"role":      "vendor",
		"authtoken": authToken,
		"medicines": []models.Medicine{},
		"createdat": time.Now(),
		"updatedat": time.Now(),
	}}

	_, red := collection.UpdateOne(r.Context(), fil, update)
	if red != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)

}

func (uc UserController) CreateMedicine(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	med := models.Medicine{}
	json.NewDecoder(r.Body).Decode(&med)
	db := uc.client.Database("medizone")
	ownerCollection := db.Collection("medowners")
	medicineCollection := db.Collection("medicines")
	var owner models.MedOwner

	err := ownerCollection.FindOne(r.Context(), bson.M{"authtoken": med.MedOwnerID, "role": "vendor"}).Decode(&owner)
	if err == nil {
		result, medErr := medicineCollection.InsertOne(r.Context(), med)
		if medErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		uj, respErr := json.Marshal(result)
		if respErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authToken := (result.InsertedID.(primitive.ObjectID)).Hex()
		filter := bson.M{"_id": result.InsertedID.(primitive.ObjectID)}
		update := bson.M{
			"$set": bson.M{
				"authtoken": authToken,
				"createdat": time.Now(),
				"updatedat": time.Now(),
			},
		}

		_, updMedErr := medicineCollection.UpdateOne(r.Context(), filter, update)
		if updMedErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var medUpdated models.Medicine
		medUpdErr := medicineCollection.FindOne(r.Context(), filter).Decode(&medUpdated)
		if medUpdErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		_, errUPD := ownerCollection.UpdateOne(r.Context(), bson.M{"authtoken": med.MedOwnerID}, bson.M{"$push": bson.M{"medicines": medUpdated}})
		if errUPD != nil {

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s\n", uj)

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (uc UserController) AddToCart(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cart := models.Cart{}
	json.NewDecoder(r.Body).Decode(&cart)
	db := uc.client.Database("medizone")
	cartCollection := db.Collection("carts")
	medicineCollection := db.Collection("medicines")
	userCollection := db.Collection("users")
	user := models.User{}
	userErr := userCollection.FindOne(r.Context(), bson.M{"authtoken": cart.UserID}).Decode(&user)
	medicine := models.Medicine{}

	err := medicineCollection.FindOne(r.Context(), bson.M{"authtoken": cart.Medicine.AuthToken}).Decode(&medicine)
	if err == nil && userErr == nil {
		foundCart := models.Cart{}
		cartFilter := bson.M{"medicine.authtoken": bson.M{"$eq": cart.Medicine.AuthToken}, "userid": bson.M{"$eq": cart.UserID}, "orderplaced": false, "isdeleted": false}
		fmt.Println(cartFilter)
		cartErr := cartCollection.FindOne(r.Context(), cartFilter).Decode(&foundCart)
		if cartErr == nil {
			addCart := bson.M{"$set": bson.M{"quantity": foundCart.Quantity + 1, "total": foundCart.Total + medicine.Price, "updatedat": time.Now()}}
			_, updErr := cartCollection.UpdateOne(r.Context(), cartFilter, addCart)
			if updErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var updCart models.Cart
			errUpd := cartCollection.FindOne(r.Context(), bson.M{"authtoken": foundCart.AuthToken}).Decode(&updCart)

			if errUpd != nil {
				w.WriteHeader((http.StatusInternalServerError))
				return
			}

			totalCartPrice := uc.cartTotal(user.CartItems, r)
			updatedCart := map[string]interface{}{
				"code":       200,
				"status":     true,
				"message":    "Cart Updated Successfully",
				"data":       updCart,
				"Cart Total": totalCartPrice,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedCart)

		} else {

			insertedCart := models.Cart{}
			insertedCart.Medicine = medicine
			insertedCart.Quantity = 1
			insertedCart.Total = medicine.Price
			insertedCart.UserID = user.AuthToken
			result, insertErr := cartCollection.InsertOne(r.Context(), insertedCart)
			if insertErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			cartAuthToken := (result.InsertedID.(primitive.ObjectID)).Hex()
			cartFilter := bson.M{"_id": result.InsertedID.(primitive.ObjectID)}
			cartUpdate := bson.M{"$set": bson.M{"authtoken": cartAuthToken, "createdat": time.Now(), "updatedat": time.Now()}}
			_, updErr := cartCollection.UpdateOne(r.Context(), cartFilter, cartUpdate)
			if updErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, userUpdErr := userCollection.UpdateOne(r.Context(), bson.M{"authtoken": cart.UserID}, bson.M{"$push": bson.M{"cartitems": cartAuthToken}})
			if userUpdErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var userCartUpd models.User
			userCartUpdErr := userCollection.FindOne(r.Context(), bson.M{"authtoken": cart.UserID}).Decode(&userCartUpd)
			if userCartUpdErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			totalCartPrice := uc.cartTotal(userCartUpd.CartItems, r)

			insertCart := map[string]interface{}{
				"code":       200,
				"status":     true,
				"message":    "New Item Added to your Cart Successfully",
				"data":       result,
				"Cart Total": totalCartPrice,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(insertCart)

		}

	} else {

		errorCart := map[string]interface{}{
			"code":    404,
			"status":  false,
			"message": "Medicine or User Not Found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorCart)

	}

}

func (uc UserController) RemoveFromCart(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authToken := r.URL.Query().Get("cartToken")
	status := r.URL.Query().Get("status")
	db := uc.client.Database("medizone")
	cartCollection := db.Collection("carts")
	userCollection := db.Collection("users")
	var user models.User
	var cart models.Cart
	var cartUpdated models.Cart
	fmt.Println(authToken)
	fmt.Println(status)
	filter := bson.M{"authtoken": authToken}
	cartFindErr := cartCollection.FindOne(r.Context(), filter).Decode(&cart)
	userErr := userCollection.FindOne(r.Context(), bson.M{"authtoken": cart.UserID}).Decode(&user)
	if cartFindErr == nil && userErr == nil {
		if cart.Quantity != 1 && status == "true" {
			cartCollection.UpdateOne(r.Context(), filter, bson.M{"$set": bson.M{"quantity": cart.Quantity - 1, "total": cart.Total - cart.Medicine.Price, "updatedat": time.Now()}})
			cartUpdatedErr := cartCollection.FindOne(r.Context(), filter).Decode(&cartUpdated)
			if cartUpdatedErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			totalCartPrice := uc.cartTotal(user.CartItems, r)
			response := map[string]interface{}{
				"code":       200,
				"status":     true,
				"message":    "Item Removed from Your Cart Successfully",
				"data":       cartUpdated,
				"Total Cart": totalCartPrice,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)

		} else if cart.Quantity == 1 || status == "false" {
			cartCollection.UpdateOne(r.Context(), filter, bson.M{"$set": bson.M{"isdeleted": true, "deletedat": time.Now(), "updatedat": time.Now()}})
			cartUpdatedErr := cartCollection.FindOne(r.Context(), filter).Decode(&cartUpdated)
			if cartUpdatedErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, userUpdErr := userCollection.UpdateOne(r.Context(), bson.M{"authtoken": cart.UserID}, bson.M{"$pull": bson.M{"cartitems": authToken}})
			if userUpdErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var userCartUpd models.User
			userCartUpdErr := userCollection.FindOne(r.Context(), bson.M{"authtoken": cart.UserID}).Decode(&userCartUpd)
			if userCartUpdErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			totalCartPrice := uc.cartTotal(userCartUpd.CartItems, r)
			response := map[string]interface{}{
				"code":       200,
				"status":     true,
				"message":    "Item Deleted from Your Cart Successfully",
				"data":       cartUpdated,
				"Total Cart": totalCartPrice,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
	} else {
		errorCart := map[string]interface{}{
			"code":    404,
			"status":  false,
			"message": "Cart Not Found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorCart)
	}

}

func (uc UserController) Checkout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := models.User{}
	order := models.Order{}
	json.NewDecoder(r.Body).Decode(&user)
	db := uc.client.Database("medizone")
	orderCollection := db.Collection("orders")
	userCollection := db.Collection("users")
	cartCollection := db.Collection("carts")

	for _, cart := range user.CartItems {
		var cartObj models.Cart
		foundNot := cartCollection.FindOne(r.Context(), bson.M{"authtoken": cart}).Decode(&cartObj)
		if foundNot == nil && cartObj.UserID == user.AuthToken {
			order.CartItems = append(order.CartItems, cartObj)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
	var reqUser models.User
	err := userCollection.FindOne(r.Context(), bson.M{"authtoken": user.AuthToken}).Decode(&reqUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	order.User = reqUser
	order.Tax = "18% GST"
	totalCartPrice := uc.cartTotal(user.CartItems, r)
	order.Total = totalCartPrice
	order.SubTotal = totalCartPrice + totalCartPrice*0.18
	order.PaymentType = "Cash on Delivery"
	order.ShippingCost = 0.0

	result, orderInsert := orderCollection.InsertOne(r.Context(), order)
	if orderInsert != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authToken := (result.InsertedID.(primitive.ObjectID)).Hex()
	filter := bson.M{"_id": result.InsertedID.(primitive.ObjectID)}
	update := bson.M{
		"$set": bson.M{
			"authtoken": authToken,
			"createdat": time.Now(),
			"updatedat": time.Now(),
		},
	}

	_, updOrderErr := orderCollection.UpdateOne(r.Context(), filter, update)
	if updOrderErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var orderUpdated models.Order
	orderUpdErr := orderCollection.FindOne(r.Context(), filter).Decode(&orderUpdated)
	if orderUpdErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	response := map[string]interface{}{
		"code":    200,
		"status":  true,
		"message": "Order Created Successfully",
		"data":    orderUpdated,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (uc UserController) PlaceOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	reqOrder := models.Order{}
	authtoken := r.URL.Query().Get("orderToken")

	json.NewDecoder(r.Body).Decode(&reqOrder)
	db := uc.client.Database("medizone")
	orderCollection := db.Collection("orders")
	userCollection := db.Collection("users")
	cartCollection := db.Collection("carts")

	var order models.Order

	orderErr := orderCollection.FindOne(r.Context(), bson.M{"authtoken": authtoken}).Decode(&order)
	emailUpdReq := uc.validateEmail(reqOrder.User.Email, order.User.Email, r, w)
	phoneUpdReq := uc.validatePhone(reqOrder.User.Phone, order.User.Phone, r, w)

	if emailUpdReq == false {
		response := map[string]interface{}{
			"code":    400,
			"status":  false,
			"message": "Email Already In Use",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	if phoneUpdReq == false {
		response := map[string]interface{}{
			"code":    400,
			"status":  false,
			"message": "Phone Number Already In Use",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	if orderErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, cart := range order.CartItems {
		var cartObj models.Cart
		foundNot := cartCollection.FindOne(r.Context(), bson.M{"authtoken": cart.AuthToken}).Decode(&cartObj)
		if foundNot == nil {
			_, updErr := cartCollection.UpdateOne(r.Context(), bson.M{"authtoken": cart.AuthToken}, bson.M{"$set": bson.M{"orderplaced": true, "updatedat": time.Now()}})
			if updErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Println("rugby")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	update := bson.M{
		"$set": bson.M{
			"firstname":        reqOrder.User.FirstName,
			"lastname":         reqOrder.User.LastName,
			"email":            reqOrder.User.Email,
			"medicalstorename": reqOrder.User.MedicalStoreName,
			"country":          reqOrder.User.Country,
			"streetaddress_v1": reqOrder.User.StreetAddress_v1,
			"streetaddress_v2": reqOrder.User.StreetAddress_v2,
			"town":             reqOrder.User.Town,
			"state":            reqOrder.User.State,
			"postcode":         reqOrder.User.PostCode,
			"phone":            reqOrder.User.Phone,
			"updatedat":        time.Now(),
		},
	}
	_, err := userCollection.UpdateOne(r.Context(), bson.M{"authtoken": order.User.AuthToken}, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var respUser models.User
	respErr := userCollection.FindOne(r.Context(), bson.M{"authtoken": order.User.AuthToken}).Decode(&respUser)
	if respErr != nil {
		fmt.Println(order.User.AuthToken)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, orderUpdErr := orderCollection.UpdateOne(r.Context(), bson.M{"authtoken": authtoken}, bson.M{"$set": bson.M{"cartitems.$[].orderplaced": true, "updatedat": time.Now(), "user": respUser}})
	if orderUpdErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code":    200,
		"status":  true,
		"message": "Order Placed Successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (uc UserController) GetCarts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := models.User{}
	authToken := r.URL.Query().Get("userToken")
	db := uc.client.Database("medizone")
	userCollection := db.Collection("users")
	cartCollection := db.Collection("carts")
	userCollection.FindOne(r.Context(), bson.M{"authtoken": authToken}).Decode(&user)
	var respAllCarts []models.Cart

	for _, cart := range user.CartItems {
		var cartObj models.Cart
		foundNot := cartCollection.FindOne(r.Context(), bson.M{"authtoken": cart}).Decode(&cartObj)
		if foundNot == nil {
			respAllCarts = append(respAllCarts, cartObj)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	response := map[string]interface{}{
		"code":    200,
		"status":  true,
		"records": len(respAllCarts),
		"data":    respAllCarts,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (uc UserController) GetMedicines(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	medicineCollection := uc.client.Database("medizone").Collection("medicines")
	medDoc, medErr := medicineCollection.Find(r.Context(), bson.M{})
	if medErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var medicines []models.Medicine
	err := medDoc.All(r.Context(), &medicines)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code":    200,
		"status":  true,
		"records": len(medicines),
		"data":    medicines,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (uc UserController) cartTotal(cartArray []string, r *http.Request) float32 {
	collection := uc.client.Database("medizone").Collection("carts")
	var totalCart float32
	for _, cart := range cartArray {
		var cartObj models.Cart
		foundNot := collection.FindOne(r.Context(), bson.M{"authtoken": cart}).Decode(&cartObj)
		if foundNot == nil {
			totalCart += cartObj.Total
		} else {
			return 0.0
		}
	}
	return totalCart
}

func (uc UserController) validateEmail(givenEmail string, presentEmail string, r *http.Request, w http.ResponseWriter) bool {
	var updateRequired bool
	var dummy []models.User
	userCollection := uc.client.Database("medizone").Collection("users")
	emailFound, emailFoundErr := userCollection.Find(r.Context(), bson.M{"email": givenEmail})
	if emailFoundErr != nil {
		updateRequired = true
	} else {
		countErr := emailFound.All(r.Context(), &dummy)
		if countErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}
		if len(dummy) == 1 {
			if presentEmail == givenEmail {
				updateRequired = true

			} else {
				updateRequired = false
			}
		} else if len(dummy) == 0 {
			updateRequired = true
		} else {
			updateRequired = false
		}
	}
	return updateRequired
}

func (uc UserController) validatePhone(givenPhone string, presentPhone string, r *http.Request, w http.ResponseWriter) bool {
	var updateRequired bool
	var dummy []models.User
	userCollection := uc.client.Database("medizone").Collection("users")
	emailFound, emailFoundErr := userCollection.Find(r.Context(), bson.M{"phone": givenPhone})
	if emailFoundErr != nil {
		updateRequired = true
	} else {
		countErr := emailFound.All(r.Context(), &dummy)
		if countErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}
		if len(dummy) == 1 {
			if presentPhone == givenPhone {
				updateRequired = true

			} else {
				updateRequired = false
			}
		} else if len(dummy) == 0 {
			updateRequired = true
		} else {
			updateRequired = false
		}
	}
	return updateRequired
}
