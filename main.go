package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
)

type customer struct {
	Id        string `bson:"_id,omitempty"`
	FirstName string `bson:"firstName,omitempty"`
	LastName  string `bson:"lastName,omitempty"`
	Age       int    `bson:"age,omitempty"`
}

const uri = "mongodb://127.0.0.1:27017"
const database = "pub"
const collectionPrimary = "pub_customers"

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("[MONGODB] Successfully connected and pinged.")
	fmt.Println("first findOne:")

	//print getted document
	//documentRaw, err := getCustomer(*client, "10")
	//document, err := json.MarshalIndent(documentRaw, "", "   ")
	//fmt.Printf("%s\n", document)

	router := gin.Default()

	router.POST("/visitPub", postCustomer)
	router.Run("localhost:8080")

}
func insertCustomer(client mongo.Client, customer customer) {
	coll := client.Database(database).Collection(collectionPrimary)

	foundCustomer, _ := getCustomer(client, customer.Id)

	if foundCustomer == nil {
		fmt.Printf("inserting: id: %s, first name: %s, last name: %s, age: %d\n", customer.Id, customer.FirstName, customer.LastName, customer.Age)
		result, err := coll.InsertOne(context.TODO(), customer)
		if err != nil {
			fmt.Println("error occured!", err)
		}
		fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	} else {
		fmt.Printf("customer with id %s already exists!", customer.Id)

	}

}

func getCustomer(client mongo.Client, id string) (bson.M, error) {
	coll := client.Database(database).Collection(collectionPrimary)
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("mongo doc with id [%s] not found!", id)
			return nil, err
		}
		fmt.Println("unexpected error! \n", err)
	}
	return result, err
}

func postCustomer(c *gin.Context) {
	var newCustomer customer

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newCustomer); err != nil {
		return
	}

	// Add the new album to the slice.

	c.IndentedJSON(http.StatusCreated, newCustomer)
	fmt.Println(newCustomer.Id)
}
