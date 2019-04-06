package dao

import (
	"context"
	"encoding/json"
	"generic_inventory/conf"
	"io"
	"io/ioutil"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConfigDB - Function to set up a DB client
func ConfigDB(conf *conf.ServerConf) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+conf.DBHost+":"+conf.DBPort))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	conf.DBClient = client.Database(conf.DBName)
}

// GetEntries - Return a json object containing people
func GetEntries(conf conf.ServerConf) (results []Entry, err error) {
	ctx, close := context.WithTimeout(context.Background(), 30*time.Second)
	defer close()
	c, err := conf.DBClient.Collection("entries").Find(ctx, bson.D{})
	defer c.Close(ctx)

	for c.Next(ctx) {
		var result Entry
		err = c.Decode(&result)
		results = append(results, result)
	}
	err = c.Err()
	return
}

// GetEntry - Return a json object containing one person
func GetEntry(sku string, conf conf.ServerConf) (result Entry, err error) {
	filter := bson.M{"sku": sku}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	err = conf.DBClient.Collection("entries").FindOne(ctx, filter).Decode(&result)
	return
}

// CreateEntry - Create a json object containing one person
func CreateEntry(reqbody io.ReadCloser, sku string, conf conf.ServerConf) (res *mongo.InsertOneResult, err error) {
	var entry Entry
	entry.SKU = sku
	body, err := ioutil.ReadAll(io.LimitReader(reqbody, 1048576))
	err = json.Unmarshal(body, &entry)
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	res, err = conf.DBClient.Collection("entries").InsertOne(ctx, entry)
	return
}

// DeleteEntry - Delete an entry
func DeleteEntry(sku string, conf conf.ServerConf) (count int64, err error) {
	filter := bson.M{"sku": sku}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	res, err := conf.DBClient.Collection("entries").DeleteOne(ctx, filter)
	count = res.DeletedCount
	return
}
