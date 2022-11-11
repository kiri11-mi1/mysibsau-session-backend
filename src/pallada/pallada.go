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

func (api *OdooAPI) GetAllSessionsIds(groupId int64) ([]int64, error) {
	options := make(odoo.Options)
	options["fields"] = []string{"session_ids"}
	gia := GroupInfoArray{}
	err := api.Client.Read("info.groups", []int64{groupId}, &options, &gia)
	if err != nil {
		log.Fatal(err)
	}
	if len(gia) == 0 {
		return []int64{}, errors.New("no sessions for this group")
	}
	return gia[0].ConvertIdsToInt(), nil
}

func (api *OdooAPI) GetCurrentYears(groupId int64) (int64, int64, error) {
	options := make(odoo.Options)
	options["fields"] = []string{"cur_year_header"}
	gia := GroupInfoArray{}
	err := api.Client.Read("info.groups", []int64{groupId}, &options, &gia)
	if err != nil {
		log.Fatal(err)
	}
	if len(gia) == 0 {
		return -1, -1, errors.New("no sessions for this group")
	}
	if len(gia[0].CurrentYears) == 0 {
		return -1, -1, errors.New("not current years")
	}
	curYearsString := strings.Split(gia[0].CurrentYears, " - ")
	begin, _ := strconv.Atoi(curYearsString[0])
	end, _ := strconv.Atoi(curYearsString[1])
	return int64(begin), int64(end), nil
}

func (api *OdooAPI) GetSessionByGroupID(groupId int64) ([]Exam, error) {
	exams := Exams{}
	options := make(odoo.Options)
	options["fields"] = []string{
		"year", "group", "employee_name_init", "lesson", "place", "day_week", "time", "date",
	}
	sessionIds, err := api.GetAllSessionsIds(groupId)
	if err != nil {
		return Exams{}, errors.New("Not exams")
	}
	err = api.Client.Read("info.timetable", sessionIds, &options, &exams)
	if err != nil {
		log.Fatal(err)
	}
	curYearStart, curYearFinish, err := api.GetCurrentYears(groupId)
	if err != nil {
		return Exams{}, errors.New("Not exams")
	}
	currentExams := Exams{}
	for _, exam := range exams {
		if exam.Year >= curYearStart && exam.Year <= curYearFinish {
			currentExams = append(currentExams, exam)
		}
	}
	if len(currentExams) == 0 {
		return Exams{}, errors.New("Not exams")
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
