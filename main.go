package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var items []Item

func main() {
	router := mux.NewRouter()

	items = append(items, Item{ID: 1, Name: "santoor"}, Item{ID: 2, Name: "nirma"})

	router.HandleFunc("/items", getItems).Methods("GET")
	router.HandleFunc("/items/{id}", getItem).Methods("GET")
	router.HandleFunc("/items", createItems).Methods("POST")
	router.HandleFunc("/items/{id}", updateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	for _, item := range items {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.NotFound(w, r)
}

func createItems(w http.ResponseWriter, r *http.Request) {
	var newItems []Item
	err := json.NewDecoder(r.Body).Decode(&newItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process each item
	for _, newItem := range newItems {
		newItem.ID = len(items) + 1 // simple ID generation, needs improvement for real applications
		items = append(items, newItem)
	}

	json.NewEncoder(w).Encode(newItems)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	for i, item := range items {
		if item.ID == id {
			var updatedItem Item
			if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			items[i] = updatedItem // Update the item
			json.NewEncoder(w).Encode(updatedItem)
			return
		}
	}
	http.NotFound(w, r)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}
