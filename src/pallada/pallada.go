package pallada

import (
	"github.com/skilld-labs/go-odoo"
	"log"
	"strconv"
	"strings"
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

type Exam struct {
	Year          int64         `xmlrpc:"year"`
	Group         []interface{} `xmlrpc:"group"`
	ProfessorName string        `xmlrpc:"employee_name_init"`
	Subject       []interface{} `xmlrpc:"lesson"`
	Room          string        `xmlrpc:"place"`
	DayWeek       string        `xmlrpc:"day_week"`
	Time          string        `xmlrpc:"time"`
	Date          string        `xmlrpc:"date"`
}
type Exams []Exam

func (api *OdooAPI) GetAllSessionsIds(groupId int64) []interface{} {
	options := make(odoo.Options)
	options["fields"] = []string{"session_ids"}
	groups := Groups{}
	err := api.Client.Read("info.groups", []int64{groupId}, &options, &groups)
	if err != nil {
		log.Fatal(err)
	}
	return groups[0].SessionsIds
}

func (api *OdooAPI) GetCurrentYears(groupId int64) (int64, int64) {
	options := make(odoo.Options)
	options["fields"] = []string{"cur_year_header"}
	groups := Groups{}
	err := api.Client.Read("info.groups", []int64{groupId}, &options, &groups)
	if err != nil {
		log.Fatal(err)
	}
	curYearsString := strings.Split(groups[0].CurrentYears, " - ")
	begin, _ := strconv.Atoi(curYearsString[0])
	end, _ := strconv.Atoi(curYearsString[1])
	return int64(begin), int64(end)
}

func (api *OdooAPI) GetSessionByGroupID(groupId int64) Exams {
	exams := Exams{}
	options := make(odoo.Options)
	options["fields"] = []string{
		"year", "group", "employee_name_init", "lesson", "place", "day_week", "time", "date",
	}
	sessionsArray := api.GetAllSessionsIds(groupId)
	newSessionsIds := make([]int64, len(sessionsArray))
	for i, sessionId := range api.GetAllSessionsIds(groupId) {
		newSessionsIds[i] = sessionId.(int64)
	}

	err := api.Client.Read("info.timetable", newSessionsIds, &options, &exams)
	if err != nil {
		log.Fatal(err)
	}
	curYearStart, curYearFinish := api.GetCurrentYears(groupId)
	currentExams := make(Exams, len(exams))
	for _, exam := range exams {
		if exam.Year >= curYearStart && exam.Year <= curYearFinish {
			currentExams = append(currentExams, exam)
		}
	}
	return currentExams
}
