package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"github.com/skilld-labs/go-odoo"
	"log"
	"mysibsau-session-backend/constants"
	"mysibsau-session-backend/pallada"
	"net/http"
)

var creds credentials

func session(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recover() != nil {
			http.Error(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
		}
	}()
	queryParams := mux.Vars(r)
	groupName := queryParams["groupName"]
	log.Print(groupName)

	if err := env.Parse(&creds); err != nil {
		log.Print(err)
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
	api := pallada.OdooAPI{Client: client}
	groupId, err := api.GetGroupIdByName(groupName)
	if err != nil {
		resp, _ := json.Marshal(constants.GROUP_NOT_FOUND)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}
	sessionTimetable, err := api.GetSessionByGroupID(groupId)
	if err != nil {
		resp, _ := json.Marshal(constants.SESSION_NOT_FOUND)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(&sessionTimetable)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func main() {
	var router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/session/{groupName}", session).Methods("GET")
	log.Print(http.ListenAndServe(":8080", router))
}
