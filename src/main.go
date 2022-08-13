package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func session(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recover() != nil {
			http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
		}
	}()
	g := Group{Id: 4, Name: "РС18-01"}
	err := json.NewEncoder(w).Encode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	var router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/session", session).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
