package api

import (
	"context"
	"encoding/json"
	"generic_inventory/web"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// GetEntries - Return a json object containing people
func GetEntries(w http.ResponseWriter, r *http.Request) {
	var entries []Entry
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx, close := context.WithTimeout(context.Background(), 30*time.Second)
	defer close()
	c, err := InventoryDB.Collection("entries").Find(ctx, bson.D{})
	if err != nil {
		w.WriteHeader(404) // Not Found
	}
	defer c.Close(ctx)

	for c.Next(ctx) {
		var result Entry
		err := c.Decode(&result)
		if err != nil {
			w.WriteHeader(500) // Internal error
		}
		entries = append(entries, result)
	}
	if err := c.Err(); err != nil {
		w.WriteHeader(500) // Internal error
	}

	if err := json.NewEncoder(w).Encode(entries); err != nil {
		w.WriteHeader(500) // Internal error
	}
}

// GetEntry - Return a json object containing one person
func GetEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var params = mux.Vars(r)
	var result Entry
	filter := bson.M{"sku": params["sku"]}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	err := InventoryDB.Collection("entries").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		w.WriteHeader(404) // Not Found
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		w.WriteHeader(500) // Internal error
	}
}

// CreateEntry - Create a json object containing one person
func CreateEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var entry Entry
	params := mux.Vars(r)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		w.WriteHeader(413) // Too Large
	}

	if err := r.Body.Close(); err != nil {
		w.WriteHeader(500) // Internal error
	}

	if err := json.Unmarshal(body, &entry); err != nil {
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			w.WriteHeader(500) // Internal error
		}
	}

	entry.SKU = params["sku"]
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	res, err := InventoryDB.Collection("entries").InsertOne(ctx, entry)
	id := res.InsertedID

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		w.WriteHeader(500) // Internal error
	}
}

// DeleteEntry - Delete an entry
func DeleteEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	filter := bson.M{"sku": params["sku"]}
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	res, err := InventoryDB.Collection("entries").DeleteOne(ctx, filter)
	if err != nil {
		w.WriteHeader(500) // Internal error
	}
	count := res.DeletedCount

	if count == 0 {
		w.WriteHeader(404) // Not Found
	}

	if err := json.NewEncoder(w).Encode(count); err != nil {
		w.WriteHeader(500) // Internal error
	}
}

// GetIndex - Return the main HTML page for the site
func GetIndex(w http.ResponseWriter, r *http.Request) {
	p := &web.Page{Title: "Hello", Body: []byte("This is a sample Page.")}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	web.RenderTemplate(w, "view", p)
}

// GetCSS - Return CSS Files from the Filesystem
func GetCSS(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "C:/Users/dusti/go/src/generic_inventory/web/static/css/"+params["cssfile"])
}

// GetJS - Return Javascript Files from the Filesystem
func GetJS(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, "C:/Users/dusti/go/src/generic_inventory/web/static/js/"+params["jsfile"])
}

// GetIMG - Return Image Files from the Filesystem
func GetIMG(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, "C:/Users/dusti/go/src/generic_inventory/web/static/img/"+params["imgfile"])
}
