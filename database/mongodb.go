package database

import (
	"database/sql"
	"fmt"
	"log"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Db *sql.DB
}

func NewMongoDB(dbhost, dbport, dbuser, dbpass, dbname string) MongoDB {
	mongodb := MongoDB{}

	mgoInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport,
		dbuser, dbpass, dbname)
	log.Println(mgoInfo)

	db, err := sql.Open("mongodb", mgoInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Database connection established!")

	mongodb.Db = db

	return mongodb
}

/*
// Returns string representation of the key stored in mongoDB.
func (m *MongoDB) Get(key string) string {
	return m.Db.Get(key).Val()
}

// Sets the key to the mongoDB.
func (m *MongoDB) Set(key string, value interface{}) {
	m.Db.Set(key, value, 0)
}
*/
