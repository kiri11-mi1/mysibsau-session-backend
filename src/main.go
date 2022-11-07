package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/skilld-labs/go-odoo"
	"log"
	"mysibsau-session-backend/pallada"
	"net/http"
)

var creds credentials

type Group struct {
	Id   int64  `json:"id"`
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

func foo() {
	// Авторизация под аккаунтом администратора
	if err := env.Parse(&creds); err != nil {
		log.Fatalln(err)
	}
	client, err := odoo.NewClient(&odoo.ClientConfig{
		Admin:    creds.Admin,
		Password: creds.Password,
		Database: creds.Database,
		URL:      fmt.Sprintf("https://%s.pallada.sibsau.ru", creds.Database),
	})
	if err != nil {
		log.Fatal(err)
	}
	var id int64 = 11806
	api := pallada.OdooAPI{Client: client}
	fmt.Println(api.GetSessionByGroupID(id))

}

func main() {
	foo()
	//var router = mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/session", session).Methods("GET")
	//log.Fatal(http.ListenAndServe(":8080", router))
}
