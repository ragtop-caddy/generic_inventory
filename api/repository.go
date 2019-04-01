package api

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// InventoryDAO - Type to hold a couple of parameters needed to set up a DB connection
type InventoryDAO struct {
	URI string
}

// InventoryDB - exported db client
var InventoryDB *mongo.Database

// ConfigDB - Function to set up a DB client
func (d *InventoryDAO) ConfigDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.URI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	InventoryDB = client.Database("inventory")
}

// GetEntries - Return a json object containing people
func GetEntries() (results []Entry, err error) {
	ctx, close := context.WithTimeout(context.Background(), 30*time.Second)
	defer close()
	c, err := InventoryDB.Collection("entries").Find(ctx, bson.D{})
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
func GetEntry(sku string) (result Entry, err error) {
	filter := bson.M{"sku": sku}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	err = InventoryDB.Collection("entries").FindOne(ctx, filter).Decode(&result)
	return
}

// CreateEntry - Create a json object containing one person
func CreateEntry(reqbody io.ReadCloser, sku string) (res *mongo.InsertOneResult, err error) {
	var entry Entry
	entry.SKU = sku
	body, err := ioutil.ReadAll(io.LimitReader(reqbody, 1048576))
	err = json.Unmarshal(body, &entry)
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	res, err = InventoryDB.Collection("entries").InsertOne(ctx, entry)
	return
}

// DeleteEntry - Delete an entry
func DeleteEntry(sku string) (count int64, err error) {
	filter := bson.M{"sku": sku}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	res, err := InventoryDB.Collection("entries").DeleteOne(ctx, filter)
	count = res.DeletedCount
	return
}
