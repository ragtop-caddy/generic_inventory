package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// InternalAuth - Object that implements the Authenticator interface
type InternalAuth struct {
	Client *mongo.Database
}

// internalUser - Internally managed user
type credentials struct {
	UID        string    `json:"uid,omitempty" bson:"uid,omitempty"`
	Fname      string    `json:"fname,omitempty" bson:"fname,omitempty"`
	Lname      string    `json:"lname,omitempty" bson:"lname,omitempty"`
	Email      string    `json:"email,omitempty" bson:"email,omitempty"`
	Pass       string    `json:"pass,omitempty" bson:"pass,omitempty"`
	Hash       string    `json:"hash,omitempty" bson:"hash,omitempty"`
	Role       string    `json:"role,omitempty" bson:"role,omitempty"`
	State      string    `json:"state,omitempty" bson:"state,omitempty"`
	Inactive   int       `json:"inactive,omitempty" bson:"inactive,omitempty"`
	Expiration int       `json:"expiration,omitempty" bson:"expiration,omitempty"`
	Last       time.Time `json:"last,omitempty" bson:"last,omitempty"`
}

type mcredentials []credentials

// HashPassword - Create a hash of the provided password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash - Compare the hashed and provided passwords
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenDefaultUser - Function to generate a default user
func (ia InternalAuth) GenDefaultUser() {
	filter := bson.M{"uid": "Administrator"}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	n := 12
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	fmt.Println(s)
	passhash, err := HashPassword(s)
	if err != nil {
		log.Fatalf("ERROR: Unable to generate default password: %s", err)
	}
	defaultUser := bson.M{
		"$set": bson.M{
			"uid":        "Administrator",
			"fname":      "Inventory",
			"lname":      "Administrator",
			"email":      "administrator@localhost",
			"hash":       passhash,
			"role":       "admin",
			"state":      "enabled",
			"inactive":   60,
			"expiration": 120,
		},
		"$currentDate": bson.M{
			"last": true,
		},
	}
	upsert := true
	udopts := &options.UpdateOptions{Upsert: &upsert}
	res, err := ia.Client.Collection("users").UpdateOne(ctx, filter, defaultUser, udopts)
	if err != nil {
		log.Fatalf("ERROR: Unable to initialize the default account: %s", err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", res.MatchedCount, res.ModifiedCount)
}

// Authenticate - Check password validity
func (ia InternalAuth) Authenticate(uid, pass string) (u User) {
	u.Username = uid
	u.Authenticated = false
	creds, err := ia.retrieveUser(uid)
	if err != nil {
		log.Printf("ERROR: %s while retrieving user credentials", err)
	}
	if CheckPasswordHash(pass, creds.Hash) {
		u.Authenticated = true
	}
	u.Role = creds.Role
	return u
}

func (ia InternalAuth) retrieveUser(uid string) (creds credentials, err error) {
	filter := bson.M{"uid": uid}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	err = ia.Client.Collection("users").FindOne(ctx, filter).Decode(&creds)
	return
}

func (ia InternalAuth) retrieveUsers() (results mcredentials, err error) {
	filter := bson.D{}
	ctx, close := context.WithTimeout(context.Background(), 30*time.Second)
	defer close()

	c, err := ia.Client.Collection("users").Find(ctx, filter)
	defer c.Close(ctx)

	for c.Next(ctx) {
		var result credentials
		err = c.Decode(&result)
		results = append(results, result)
	}
	err = c.Err()
	return
}

func (ia InternalAuth) createUser(new credentials) (res *mongo.InsertOneResult, err error) {
	hash, err := HashPassword(new.Pass)
	if err != nil {
		log.Printf("ERROR: %s while storing user credentials", err)
	}
	newuserDoc := bson.M{
		"$set": bson.M{
			"uid":        new.UID,
			"fname":      new.Fname,
			"lname":      new.Lname,
			"email":      new.Email,
			"hash":       hash,
			"role":       new.Role,
			"state":      new.State,
			"inactive":   new.Inactive,
			"expiration": new.Expiration,
		},
		"$currentDate": bson.M{
			"last": true,
		},
	}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	res, err = ia.Client.Collection("users").InsertOne(ctx, newuserDoc)
	return
}

// DeleteEntry - Delete an entry
//func (ia InternalAuth) deleteUser(uid string) (count int64, err error) {
//	filter := bson.M{"sku": sku}
//	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
//	defer close()
//	res, err := ia.Client.Collection("entries").DeleteOne(ctx, filter)
//	count = res.DeletedCount
//	return
//}

//func storeCreds(creds credentials, client *mongo.Database) (res *mongo.InsertOneResult, err error) {
//	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
//	defer close()
//	res, err = client.Collection("users").InsertOne(ctx, creds)
//	return
//}

//func updateCreds(creds credentials, ) (ok bool, err error) {
//
//}

//func deleteCreds(creds credentials) (ok bool, err error) {
//
//}
