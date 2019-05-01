package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// InternalAuth - Object that implements the Authenticator interface
type InternalAuth struct {
	Client     *mongo.Database
	FailLimit  int
	Inactive   float64
	Expiration float64
}

type passHash struct {
	Pass string `json:"pass,omitempty" bson:"pass,omitempty"`
	Hash string `json:"hash,omitempty" bson:"hash,omitempty"`
}

// internalUser - Internally managed user
type credentials struct {
	UID      string    `json:"uid,omitempty" bson:"uid,omitempty"`
	Fname    string    `json:"fname,omitempty" bson:"fname,omitempty"`
	Lname    string    `json:"lname,omitempty" bson:"lname,omitempty"`
	Email    string    `json:"email,omitempty" bson:"email,omitempty"`
	Pass     string    `json:"pass,omitempty" bson:"pass,omitempty"`
	Hash     string    `json:"hash,omitempty" bson:"hash,omitempty"`
	Role     string    `json:"role,omitempty" bson:"role,omitempty"`
	Failures int       `json:"failures" bson:"failures"`
	Enabled  bool      `json:"enabled" bson:"enabled"`
	Last     time.Time `json:"last,omitempty" bson:"last,omitempty"`
}

// DisplayCredentials - struct to hold user information safe for display
type DisplayCredentials struct {
	UID      string    `json:"uid,omitempty" bson:"uid,omitempty"`
	Fname    string    `json:"fname,omitempty" bson:"fname,omitempty"`
	Lname    string    `json:"lname,omitempty" bson:"lname,omitempty"`
	Email    string    `json:"email,omitempty" bson:"email,omitempty"`
	Role     string    `json:"role,omitempty" bson:"role,omitempty"`
	Failures int       `json:"failures" bson:"failures"`
	Enabled  bool      `json:"enabled" bson:"enabled"`
	Last     time.Time `json:"last,omitempty" bson:"last,omitempty"`
}

type mcredentials []DisplayCredentials

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

func checkTime(last time.Time) (since float64) {
	since = time.Since(last).Hours() / 24
	return
}

func genString() (s string) {
	n := 12
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s = fmt.Sprintf("%X", b)
	return
}

// GenDefaultUser - Function to generate a default user
func (ia InternalAuth) GenDefaultUser() {
	var administrator = credentials{
		UID:      "administrator",
		Fname:    "Default",
		Lname:    "Admin",
		Role:     "admin",
		Failures: 0,
		Enabled:  true,
		Last:     time.Now(),
	}
	check, err := ia.retrieveCreds(administrator.UID)
	if check.UID == administrator.UID {
		return
	}

	administrator.Pass = genString()

	_, err = ia.createUser(administrator)
	if err != nil {
		fmt.Printf("ERROR: Could not create default admin user: %v", err)
		return
	}
	fmt.Printf("INFO: Set Default Password to %s \n", administrator.Pass)
}

// Authenticate - Check password validity
func (ia InternalAuth) Authenticate(uid, pass string) (u User) {
	u.Username = uid
	u.Authenticated = false
	creds, err := ia.retrieveCreds(uid)
	if err != nil {
		log.Printf("ERROR: %s while retrieving user credentials", err)
	}
	u.Role = creds.Role
	ok := CheckPasswordHash(pass, creds.Hash)
	if !ok {
		creds.Failures++
		fmt.Printf("INFO: Failed login for %s \n", uid)
	} else {
		creds.Failures = 0
		creds.Last = time.Now()
	}

	if creds.Failures > ia.FailLimit {
		creds.Enabled = false
		fmt.Printf("INFO: Disabling %s account due to excessive login failures \n", uid)
	}

	since := checkTime(creds.Last)
	if since > ia.Inactive {
		creds.Enabled = false
		fmt.Printf("INFO: Disabling %s account due to inactivity \n", uid)
	}

	stored, err := ia.storeCreds(creds)
	if err != nil {
		log.Printf("ERROR: %s while updating user credentials", err)
	}

	if ok && creds.Enabled && stored {
		u.Authenticated = true
	}

	return u
}

func (ia InternalAuth) retrieveCreds(uid string) (creds credentials, err error) {
	filter := bson.M{"uid": uid}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	err = ia.Client.Collection("users").FindOne(ctx, filter).Decode(&creds)
	return
}

func (ia InternalAuth) storeCreds(creds credentials) (ok bool, err error) {
	ok = false
	if creds.Pass != "" {
		creds.Hash, err = HashPassword(creds.Pass)
		if err != nil {
			log.Printf("ERROR: %s while generating password hash for %s \n", err, creds.UID)
		}
	}
	filter := bson.M{"uid": creds.UID}
	update := bson.M{"$set": creds}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	err = ia.Client.Collection("users").FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		fmt.Printf("ERROR: Got %s while storing creds \n", err)
	} else {
		ok = true
	}
	return
}

func (ia InternalAuth) retrieveUser(uid string) (creds DisplayCredentials, err error) {
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
		var result DisplayCredentials
		err = c.Decode(&result)
		results = append(results, result)
	}
	err = c.Err()
	return
}

func (ia InternalAuth) createUser(new credentials) (res *mongo.InsertOneResult, err error) {
	new.Hash, err = HashPassword(new.Pass)
	if err != nil {
		log.Printf("ERROR: %s while generating password hash for %s \n", err, new.UID)
	}
	new.Failures = 0
	new.Last = time.Now()
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	res, err = ia.Client.Collection("users").InsertOne(ctx, new)
	return
}

func (ia InternalAuth) deleteUser(uid string) (count int64, err error) {
	filter := bson.M{"uid": uid}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	res, err := ia.Client.Collection("users").DeleteOne(ctx, filter)
	count = res.DeletedCount
	return
}
