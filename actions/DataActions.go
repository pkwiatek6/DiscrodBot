package actions

import (
	"context"
	"log"

	"github.com/pkwiatek6/DiscrodBot/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//Database connectiong to
	Database = "Characters"
	//Collection reading from
	Collection = "Sheets"
)

//SaveCharacter saves player data to a noSQL db
func SaveCharacter(character data.Character, client *mongo.Client) error {
	collection := client.Database(Database).Collection(Collection)
	//will get readded when I figure out why discordgo isn't giving be user discriminator
	//filter := bson.D{{Key: "name", Value: character.Name}}
	filter := bson.M{"name": character.Name, "user": character.User}
	update := bson.M{"$set": character}
	updateResult, err1 := collection.UpdateOne(context.TODO(), filter, update)
	if err1 != nil {
		return err1
		//checks if there was a document that was updated and if so finsih saving
	} else if updateResult.MatchedCount == 0 {
		log.Println("Failed to find matching document, making a new one")
	} else if updateResult.MatchedCount == 1 {
		log.Println("Document Was updated")
		return nil
	}
	//Creates a new document if there wasn't one already
	insertResult, err2 := collection.InsertOne(context.TODO(), character)
	if err2 != nil {
		return err2
	}
	log.Println("Inserted post with ID:", insertResult.InsertedID)
	return nil
}

//LoadCharacter loads a given character by name, I'm probably also gonna require it to look up User ID
func LoadCharacter(name string, user string, client *mongo.Client) (*data.Character, error) {
	filter := bson.M{"name": name, "user": user}
	collection := client.Database(Database).Collection(Collection)
	var character data.Character
	//Finds a document with name and decodes it into the variable character
	err := collection.FindOne(context.TODO(), filter).Decode(&character)
	if err != nil {
		return &character, err
	}
	return &character, nil
}

//LoadAllCharacters loads all the characters into memory
func LoadAllCharacters(client *mongo.Client) (map[string]*data.Character, error) {
	var results []data.Character
	var toReturn = make(map[string]*data.Character)
	cursor, err := client.Database(Database).Collection(Collection).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	cursor.All(context.TODO(), &results)
	for _, character := range results {
		toReturn[character.User] = &character
	}
	return toReturn, nil
}

//SaveAllCharacters saves all the characters to the DB
func SaveAllCharacters(Characters map[string]*data.Character, client *mongo.Client) error {
	var err error
	for _, character := range Characters {
		err = SaveCharacter(*character, client)
		if err != nil {
			return err
		}
	}
	return nil
}

//ConnectDB makes a client that can be called again and again to reference the database, call this first to create a Client
func ConnectDB() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}
	return client, nil
}
