package pallada

import (
	"errors"
	"github.com/skilld-labs/go-odoo"
	"log"
	"strconv"
	"strings"
)

type OdooAPI struct {
	Client *odoo.Client
}

type SessionInfo struct {
	Id           int64         `xmlrpc:"id"`
	SessionsIds  []interface{} `xmlrpc:"session_ids"`
	CurrentYears string        `xmlrpc:"cur_year_header"`
}

type Exam struct {
	Year      int64         `xmlrpc:"year" json:"-"`
	Professor string        `xmlrpc:"employee_name_init" json:"professor"`
	Subject   []interface{} `xmlrpc:"lesson" json:"subject"`
	Room      string        `xmlrpc:"place" json:"room"`
	DayWeek   string        `xmlrpc:"day_week" json:"day_week"`
	Time      string        `xmlrpc:"time" json:"time"`
	Date      string        `xmlrpc:"date" json:"date"`
}

type Group struct {
	Name string `xmlrpc:"name"`
	Id   int64  `xmlrpc:"id"`
}
type Groups []Group

func (api *OdooAPI) GetAllSessionsIds(groupId int64) []interface{} {
	options := make(odoo.Options)
	options["fields"] = []string{"session_ids"}
	sessionsInfo := []SessionInfo{}
	err := api.Client.Read("info.groups", []int64{groupId}, &options, &sessionsInfo)
	if err != nil {
		log.Fatal(err)
	}
	return sessionsInfo[0].SessionsIds
}

func (api *OdooAPI) GetCurrentYears(groupId int64) (int64, int64) {
	options := make(odoo.Options)
	options["fields"] = []string{"cur_year_header"}
	sessionsInfo := []SessionInfo{}
	err := api.Client.Read("info.groups", []int64{groupId}, &options, &sessionsInfo)
	if err != nil {
		log.Fatal(err)
	}
	curYearsString := strings.Split(sessionsInfo[0].CurrentYears, " - ")
	begin, _ := strconv.Atoi(curYearsString[0])
	end, _ := strconv.Atoi(curYearsString[1])
	return int64(begin), int64(end)
}

func (api *OdooAPI) GetSessionByGroupID(groupId int64) ([]Exam, error) {
	exams := []Exam{}
	options := make(odoo.Options)
	options["fields"] = []string{
		"year", "group", "employee_name_init", "lesson", "place", "day_week", "time", "date",
	}
	newSessionsIds := make([]int64, 1)
	for _, sessionId := range api.GetAllSessionsIds(groupId) {
		newSessionsIds = append(newSessionsIds, sessionId.(int64))
	}

	err := api.Client.Read("info.timetable", newSessionsIds, &options, &exams)
	if err != nil {
		log.Fatal(err)
	}
	curYearStart, curYearFinish := api.GetCurrentYears(groupId)
	currentExams := []Exam{}
	for _, exam := range exams {
		if exam.Year >= curYearStart && exam.Year <= curYearFinish {
			currentExams = append(currentExams, exam)
		}
	}
	if len(currentExams) == 0 {
		return []Exam{}, errors.New("Not exams")
	}
	return currentExams, nil
}

func (api *OdooAPI) GetGroupIdByName(nameGroup string) (int64, error) {
	options := make(odoo.Options)
	options["fields"] = []string{
		"name", "id",
	}
	var criteries = odoo.Criteria{{"name", "=", nameGroup}}
	groups := Groups{}
	err := api.Client.SearchRead("info.groups", &criteries, &options, &groups)
	if err != nil {
		log.Fatal(err)
	}
	if len(groups) == 0 {
		return -1, errors.New("there is no group")
	}
	return groups[0].Id, nil
}
