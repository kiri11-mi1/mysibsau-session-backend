package pallada

import (
	"encoding/json"
	"fmt"
	"github.com/skilld-labs/go-odoo"
)

type OdooAPI struct {
	Client *odoo.Client
}

type SessionInfo struct {
	Id           int64         `xmlrpc:"id"`
	SessionsIds  []interface{} `xmlrpc:"session_ids"`
	CurrentYears string        `xmlrpc:"cur_year_header"`
}

type Group struct {
	Name string `xmlrpc:"name"`
	Id   int64  `xmlrpc:"id"`
}

type SubjectArray []interface{}

func (sa SubjectArray) MarshalJSON() ([]byte, error) {
	if len(sa) < 2 {
		return nil, fmt.Errorf("cannot marshal subject array value %v into a string", sa)
	}
	return json.Marshal(sa[1])
}

type Exam struct {
	Year      int64        `xmlrpc:"year" json:"-"`
	Professor string       `xmlrpc:"employee_name_init" json:"professor"`
	Subject   SubjectArray `xmlrpc:"lesson" json:"subject"`
	Room      string       `xmlrpc:"place" json:"room"`
	DayWeek   string       `xmlrpc:"day_week" json:"day_week"`
	Time      string       `xmlrpc:"time" json:"time"`
	Date      string       `xmlrpc:"date" json:"date"`
}

var daysWeek = map[string]string{
	"1": "Понедельник",
	"2": "Вторник",
	"3": "Среда",
	"4": "Четверг",
	"5": "Пятница",
	"6": "Суббота",
}

type Exams []Exam
