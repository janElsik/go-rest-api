package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"strconv"
)

type customer struct {
	Id        string `bson:"_id,omitempty"`
	FirstName string `bson:"firstName,omitempty"`
	LastName  string `bson:"lastName,omitempty"`
	Age       int    `bson:"age,omitempty"`
}

var Uri = os.Getenv("MONGO_URI")
var Database = os.Getenv("MONGO_DB_NAME")
var CollectionPrimary = os.Getenv("PRIMARY_COLLECTION")

func main() {

	if Uri == "" {
		Uri = "mongodb://127.0.0.1:27017"
	}
	if Database == "" {
		Database = "pub"
	}

	if CollectionPrimary == "" {
		CollectionPrimary = "pub_customers"
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(Uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the database
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("[MONGODB] Successfully connected and pinged.")

	// start go-fiber
	app := fiber.New()

	// POST method, can be called with curl, such as
	// curl -X POST http://localhost:3000/visitPub?id=3458&firstName=Jane&lastName=Doe&age=15
	app.Post("/visitPub", func(c *fiber.Ctx) error {
		newCustomer := new(customer)

		newCustomer.Id = c.Query("id")
		newCustomer.FirstName = c.Query("firstName")
		newCustomer.LastName = c.Query("lastName")
		newCustomer.Age, err = strconv.Atoi(c.Query("age"))

		if err != nil {
			fmt.Println("error whilst converting age from string to Int!, setting customer age to 1", err)
			newCustomer.Age = 1
		}

		insertCustomer(*client, *newCustomer)
		return c.SendString("hello " + newCustomer.FirstName)
	})

	// GET method, can be called with curl, such as
	// curl -X GET 'http://localhost:3000/getCustomer?id=109403948483'
	app.Get("/getCustomer", func(c *fiber.Ctx) error {
		existingCustomerIdRaw := c.Query("id")

		customer, err := getCustomer(*client, existingCustomerIdRaw)

		if err != nil {
			fmt.Printf("couldn't find customer with id %s!", customer.Id)
		}

		fmt.Printf("id: %s \n", customer.Id)
		fmt.Printf("first name: %s \n", customer.FirstName)
		fmt.Printf("last name: %s \n", customer.LastName)
		fmt.Printf("age: %d \n", customer.Age)

		return c.SendString("customer is named: " + customer.FirstName)
	})

	log.Fatal(app.Listen(":8080"))

}
func insertCustomer(client mongo.Client, customer customer) {
	coll := client.Database(Database).Collection(CollectionPrimary)

	// first step is to check if customer exists, in order to not duplicate inserts to database
	_, err := getCustomer(client, customer.Id)

	if err != nil {
		// ErrNoDocuments is good here since it explicitly tells us that the document doesn't exist
		if err == mongo.ErrNoDocuments {
			fmt.Printf("inserting: id: %s, first name: %s, last name: %s, age: %d\n", customer.Id, customer.FirstName, customer.LastName, customer.Age)
			result, err := coll.InsertOne(context.TODO(), customer)
			if err != nil {
				fmt.Println("error with inserting customer!", err)

			} else {
				fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
			}
		}
	} else {
		fmt.Printf("customer with %s already exists! Please enter new valid personal id. \n", customer.Id)

	}

}

func getCustomer(client mongo.Client, id string) (customer, error) {
	coll := client.Database(Database).Collection(CollectionPrimary)
	var result customer
	err := coll.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("mongo doc with id [%s] not found!", id)
			return result, err
		}
		fmt.Println("unexpected error! \n", err)
	}
	return result, err
}
