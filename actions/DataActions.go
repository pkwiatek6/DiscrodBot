package actions

import (
	"context"
	"log"

	"github.com/pkwiatek6/DiscrodBot/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//SaveCharacter saves player data to a noSQL db
func SaveCharacter(character data.Character, client *mongo.Client) error {
	collection := client.Database("my_database").Collection("Characters")
	insertResult, err := collection.InsertOne(context.TODO(), character)
	if err != nil {
		return err
	}
	log.Println("Inserted post with ID:", insertResult.InsertedID)
	return nil
}

//LoadCharacter loads a given character
func LoadCharacter(name string, client *mongo.Client) (data.Character, error) {
	collection := client.Database("my_database").Collection("Characters")
	filter := bson.D{}
	var character data.Character
	err := collection.FindOne(nil, filter).Decode(&character)
	if err != nil {
		return character, err
	}
	return character, nil
}

//ConnectDB makes a client that can be called again and again to reference the database, call this first to create a Client
func ConnectDB() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	return client
}
