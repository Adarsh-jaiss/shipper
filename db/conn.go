package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

// ConnectToDB connects to the database

var (
	DBURL       string
	DB_PASSWORD string
	DBNAME      string
	USERNAME    string
	HOST        string
	PORT        string
)

func ConnectToDB() (*mongo.Client, error) {
	url, err := CreateDBURL()
	if err!= nil {
		return nil,err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		fmt.Errorf("error connecting to database: %v", err)
		return nil, err
	}

	return client, nil
}

func CreateDatabaseAndCollection(client *mongo.Client, dbName, collectionName string) error {
	collection := client.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), map[string]string{"key": "value"})
	if err != nil {
		return fmt.Errorf("error creating collection: %v", err)
	}
	return nil
}

func CreateDBURL() (string, error) {
	if os.Getenv(DBURL) == "" || len(os.Getenv(DBURL)) == 0 {

		fmt.Printf("%v is not present, creating a new URL with the config", DBURL)

		if os.Getenv(USERNAME) == "" || len(os.Getenv(USERNAME)) == 0 {
			return "", fmt.Errorf("please set %s env variable for the database", USERNAME)
		}

		if os.Getenv(DB_PASSWORD) == "" || len(os.Getenv(DB_PASSWORD)) == 0 {
			return "", fmt.Errorf("please set %s env variable for the database", DB_PASSWORD)
		}

		if os.Getenv(HOST) == "" || len(os.Getenv(HOST)) == 0 {
			return "", fmt.Errorf("please set %s env variable for the database", HOST)
		}

		if os.Getenv(PORT) == "" || len(os.Getenv(PORT)) == 0 {
			return "", fmt.Errorf("please set %s env variable for the database", PORT)
		}

		if os.Getenv(DBNAME) == "" || len(os.Getenv(DBNAME)) == 0 {
			return "", fmt.Errorf("please set %s env variable for the database", DBNAME)
		}

		DB_PASSWORD = os.Getenv(DB_PASSWORD)
		connString := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", USERNAME, DB_PASSWORD, HOST, PORT, DBNAME)
		return connString, nil

	}
	DBURL = os.Getenv(DBURL)
	return DBURL, nil

}
