package api

import (
	"context"
	"encoding/json"
	"generic_inventory/web"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// GetEntries - Return a json object containing people
func GetEntries(w http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(context.Background(), 30*time.Second)
	defer close()
	c, err := InventoryDB.Collection("entries").Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	defer c.Close(ctx)

	for c.Next(ctx) {
		var result Entry
		err := c.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		entries = append(entries, result)
	}
	if err := c.Err(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		panic(err)
	}
}

// GetEntry - Return a json object containing one person
//func (c *Controller) GetEntry(w http.ResponseWriter, r *http.Request) {
//	var params = mux.Vars(r)
//	result := c.Repository.GetDBEntry(params["sku"])
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusOK)
//	if err := json.NewEncoder(w).Encode(result); err != nil {
//		panic(err)
//	}
//
//}

// CreateEntry - Create a json object containing one person
//func CreateEntry(w http.ResponseWriter, r *http.Request) {
//	var entry Entry
//	params := mux.Vars(r)
//	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
//
//	if err != nil {
//		panic(err)
//	}
//
//	if err := r.Body.Close(); err != nil {
//		panic(err)
//	}
//
//	if err := json.Unmarshal(body, &entry); err != nil {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(422) // unprocessable entity
//		if err := json.NewEncoder(w).Encode(err); err != nil {
//			panic(err)
//		}
//	}

//	entry.SKU = params["sku"]
//	entries = append(entries, entry)
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusCreated)
//	if err := json.NewEncoder(w).Encode(entries); err != nil {
//		panic(err)
//	}
//}

// DeleteEntry - Delete a person
//func DeleteEntry(w http.ResponseWriter, r *http.Request) {
//	params := mux.Vars(r)
//	for index, entry := range entries {
//		if entry.SKU == params["sku"] {
//			entries = append(entries[:index], entries[index+1:]...)
//			break
//		}
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//		if err := json.NewEncoder(w).Encode(entry); err != nil {
//			panic(err)
//		}
//	}
//}

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
	http.ServeFile(w, r, "./css/"+params["cssfile"])
}

// GetJS - Return Javascript Files from the Filesystem
func GetJS(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, "./js/"+params["jsfile"])
}

// GetIMG - Return Image Files from the Filesystem
func GetIMG(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, "./img/"+params["imgfile"])
}
