package api

import (
	"context"
	"log"
	"time"

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
