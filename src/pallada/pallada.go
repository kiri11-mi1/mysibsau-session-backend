package pallada

import (
	"errors"
	"github.com/skilld-labs/go-odoo"
	"log"
	"strconv"
	"strings"
)

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
	exams := Exams{}
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
	currentExams := Exams{}
	for _, exam := range exams {
		if exam.Year >= curYearStart && exam.Year <= curYearFinish {
			currentExams = append(currentExams, exam)
		}
	}
	if len(currentExams) == 0 {
		return []Exam{}, errors.New("Not exams")
	}
	currentExams.ConvertDayWeek()
	currentExams.MakeValidTime()
	return currentExams.SortingByDate(), nil
}

func (api *OdooAPI) GetGroupIdByName(nameGroup string) (int64, error) {
	nameGroup = strings.ToUpper(nameGroup)
	options := make(odoo.Options)
	options["fields"] = []string{
		"name", "id",
	}
	var criteries = odoo.Criteria{{"name", "=", nameGroup}}
	groups := []Group{}
	err := api.Client.SearchRead("info.groups", &criteries, &options, &groups)
	if err != nil {
		log.Fatal(err)
	}
	if len(groups) == 0 {
		return -1, errors.New("there is no group")
	}
	return groups[0].Id, nil
}
