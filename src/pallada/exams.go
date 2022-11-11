package pallada

import (
	"log"
	"strings"
	"time"
)

type Exams []Exam

var daysWeek = map[string]string{
	"1": "Понедельник",
	"2": "Вторник",
	"3": "Среда",
	"4": "Четверг",
	"5": "Пятница",
	"6": "Суббота",
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

func (exams *Exams) SortingByDate() Exams {
	examsArray := *exams
	for i := 0; i < len(examsArray)-1; i++ {
		for j := 0; j < len(examsArray)-i-1; j++ {
			currentDate, err := time.Parse("2006-01-02", examsArray[j].Date)
			if err != nil {
				log.Print(err)
			}
			nextDate, err := time.Parse("2006-01-02", examsArray[j+1].Date)
			if err != nil {
				log.Print(err)
			}
			if currentDate.After(nextDate) {
				examsArray[j], examsArray[j+1] = examsArray[j+1], examsArray[j]
			}
		}
	}
	return examsArray
}

func (exams *Exams) ConvertDayWeek() {
	for i, e := range *exams {
		(*exams)[i].DayWeek = daysWeek[e.DayWeek]
	}
}

func (exams *Exams) MakeValidTime() {
	for i, e := range *exams {
		(*exams)[i].Time = strings.Replace(e.Time, "-", "", -1)
	}
}
