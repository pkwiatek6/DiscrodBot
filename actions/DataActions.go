package actions

import (
	"database/sql"
	"log"

	"github.com/pkwiatek6/DiscrodBot/data"
)

const (
	dbHost = "tcp(127.0.0.1:3307)"
	dbName = "character_list"
	dbUser = "user"
	dbPass = ""

//	MAX_ROWS_PER_THREAD = 100000
//	NUMBER_OF_THREADS   = 10000
)

//SaveCharacter saves player data to a database
func SaveCharacter(character data.Character) error {
	//open db
	//obtain read write lock, ie only it can read and write to the database
	//save to db
	//release locks
	//close db
	//mySQL works better here than noSQL becuase I'm saving and reading structured data
	return nil
}

//LoadCharacter loads a given character
func LoadCharacter(name string) (data.Character, error) {
	//open db
	database, err := sql.Open("mysql", dbUser+":"+dbPass+"@"+dbHost+"/"+dbName+"?charset=utf8")
	if err != nil {
		log.Println(err)
		return data.Character{}, err
	}
	//obtain write locl?, only it can write everyone can read
	//read from db
	//save to variable
	//close db
	//return variable
	//mySQL works better here than noSQL becuase I'm saving and reading structured data
	return data.Character{}, nil
}
