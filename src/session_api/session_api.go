package session_api

import (
	"github.com/skilld-labs/go-odoo"
	"log"
)

type OdooAPI struct {
	Client *odoo.Client
}

type GroupInfo struct {
	Id           int64         `xmlrpc:"id"`
	SessionsIds  []interface{} `xmlrpc:"session_ids"`
	CurrentYears string        `xmlrpc:"cur_year_header"`
}

type Groups []GroupInfo

func (api *OdooAPI) GetSessionByGroup(id int64) {
	return
}

func (api *OdooAPI) GetAllSessionsIds(id int64) Groups {
	options := make(odoo.Options)
	options["fields"] = []string{"session_ids"}
	groups := Groups{}
	err := api.Client.Read("info.groups", []int64{5080}, &options, &groups)
	if err != nil {
		log.Fatal(err)
	}
	return groups
}

func (api *OdooAPI) GetCurrentYears(id int64) {
	return
}
