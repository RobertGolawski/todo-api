package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	filePath string = "./todos.json"
	mu       sync.Mutex
)

func serverSave(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Server save error during read body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var data TodoList

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	writeData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error marshalling the data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(filePath, writeData, 0644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func serverSend(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Printf("File not found for pull request: %v", err)
		http.Error(w, "Todo list not found on server", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error stating file for pull: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Error writing response data: %v", err)
	}
}

func main() {
	http.HandleFunc("/push", serverSave)
	http.HandleFunc("/pull", serverSend)

	log.Println("Listening")

	log.Fatal(http.ListenAndServe(":8081", nil))

}
