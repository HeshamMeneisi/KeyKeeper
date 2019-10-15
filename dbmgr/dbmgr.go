package dbmgr

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"math/rand"
	"time"
)

func NewMongoClient(containerName string) (*mongo.Client, error) {
	var h string
	if containerName == "" {
		h = "127.0.0.1"
	}else{
		h = containerName
	}
	// Connect to DB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+h+":27017"))
	if err != nil {
		log.Fatal("Failed to connect to DB.")
		log.Fatal(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Failed to ping DB.")
		log.Fatal(err)
	}
	log.Println("Connected to DB")

	return client, err
}

func NextIdForCol(col *mongo.Collection) uint32 {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	options := options.Find()
	options.SetSort(bson.D{{"_id", -1}})
	options.SetLimit(1)
	var id uint32
	cursor, err := col.Find(ctx, bson.D{}, options)
	if err != nil {
		id = rand.Uint32()
	} else {
		for cursor.Next(ctx) {
			var obj struct{ID uint32 `bson:"_id"`}
			if err := cursor.Decode(&obj); err != nil {
				log.Fatal(err)
				id = rand.Uint32()
			} else {
				id = uint32(obj.ID + 1)
			}
		}
	}
	return id
}
