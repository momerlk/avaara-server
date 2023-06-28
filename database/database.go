package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client 				*mongo.Client
	mongoDB 			*mongo.Database
	Uri 				string
}
func (db *Database) Connect(uri , dbName string) error {
	db.Uri = uri	
	
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	
	db.Client = client
	db.mongoDB = db.Client.Database(dbName)
	
	return nil
}

func (db *Database) Close(l *log.Logger){
	if err := db.Client.Disconnect(context.TODO()); err != nil {
		l.Fatal(err)
	}
}

func (db *Database) Store(collName string , data any) error {
	var err error = nil
	coll := db.mongoDB.Collection(collName)
	_ , err = coll.InsertOne(context.TODO() , data)
	
	return err
}

func (db *Database) Update(collName string , data any , filter any) error {
	coll := db.mongoDB.Collection(collName)
	_ , err := coll.UpdateOne(context.TODO() , filter , bson.D{
		{"$set",data},
	})
	return err
}

func Get[T any](db *Database , collName string , filter interface{}) ([]T , error) {
	var err error = nil
	var results []T
	
	coll := db.mongoDB.Collection(collName)
	cur , err := coll.Find(context.TODO() , filter)
	if err != nil {
		return results ,  err
	}
	
	
	err = cur.All(context.TODO() , &results)
	if err != nil {
		return results , err
	}
	
	return results , err
}

func GetOne[T any](db *Database, collName string , filter interface{}) (T , error){
	var NULL T
	res , err := Get[T](db , collName , filter)
	if err != nil {
		return NULL , err
	}
	if len(res) == 0 {
		return NULL , err
	}
	return res[0] , err
}


