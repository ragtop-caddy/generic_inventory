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
	"go.mongodb.org/mongo-driver/mongo"
)

// CrudHandle - Function to call other crud funtions
func CrudHandle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	switch params["action"] {
	case "show":
		if params["sku"] == "all" {
			results, err := GetEntries()
			if err != nil {
				w.WriteHeader(404) // Not Found
			}
			if err := json.NewEncoder(w).Encode(results); err != nil {
				w.WriteHeader(500) // Internal error
			}
		} else {
			result, err := GetEntry(params["sku"])
			if err != nil {
				w.WriteHeader(404) // Not Found
			} else {
				if err := json.NewEncoder(w).Encode(result); err != nil {
					w.WriteHeader(500) // Internal error
				}
			}
		}
	case "add":
		id, err := CreateEntry(r, params["sku"])
		if err != nil {
			w.WriteHeader(500) // Internal error
			if err := json.NewEncoder(w).Encode(err); err != nil {
				w.WriteHeader(500) // Internal error
			}
		} else {
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(id); err != nil {
				w.WriteHeader(500) // Internal error
			}
		}
	case "remove":
		count, err := DeleteEntry(params["sku"])
		if err != nil {
			w.WriteHeader(500) // Internal error
			if err := json.NewEncoder(w).Encode(err); err != nil {
				w.WriteHeader(500) // Internal error
			}
		} else {
			if count == 0 {
				w.WriteHeader(404) // Not Found
			}
			if err := json.NewEncoder(w).Encode(count); err != nil {
				w.WriteHeader(500) // Internal error
			}
		}
	default:
		if err := json.NewEncoder(w).Encode(params); err != nil {
			w.WriteHeader(500) // Internal error
		}
	}
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
func CreateEntry(req *http.Request, sku string) (res *mongo.InsertOneResult, err error) {
	var entry Entry
	entry.SKU = sku
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
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
