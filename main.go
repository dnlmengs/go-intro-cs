package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	serverAddr string = "127.0.0.1:8081"
	mdbConn    string = "mongodb+srv://wadus:8ZDiaNPDZyI5mbTU@test-00.zkittdf.mongodb.net/?retryWrites=true&w=majority&appName=test-00"
)

type Scope struct {
	Project string
	Area    string
}

type Note struct {
	Title string
	Text  string
	Tags  []string
	Scope Scope
}

func createNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Note: %+v", note)

	var mdbClient *mongo.Client
	var err error
	ctxBg := context.Background()
	mdbClient, err = mongo.Connect(ctxBg, options.Client().ApplyURI(mdbConn))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = mdbClient.Disconnect(ctxBg); err != nil {
			panic(err)
		}
	}()

	notesCollection := mdbClient.Database("NoteKeeper").Collection("Notes")
	result, err := notesCollection.InsertOne(r.Context(), note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Id: %v", result.InsertedID)
}

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hola, perro."))
	})

	http.HandleFunc("POST /notes", createNote)

	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
