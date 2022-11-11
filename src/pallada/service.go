package pallada

import (
	"log"
	"strings"
	"time"
)

func (e *Exams) SortingByDate() Exams {
	exams := *e
	for i := 0; i < len(exams)-1; i++ {
		for j := 0; j < len(exams)-i-1; j++ {
			currentDate, err := time.Parse("2006-01-02", exams[j].Date)
			if err != nil {
				log.Print(err)
			}
			nextDate, err := time.Parse("2006-01-02", exams[j+1].Date)
			if err != nil {
				log.Print(err)
			}
			if currentDate.After(nextDate) {
				exams[j], exams[j+1] = exams[j+1], exams[j]
			}
		}
	}
	return exams
}

func (exams *Exams) ConvertDayWeek() {
	for _, e := range *exams {
		e.DayWeek = daysWeek[e.DayWeek]
	}
}

func (exams *Exams) MakeValidTime() {
	for _, e := range *exams {
		e.Time = strings.Replace(e.Time, "-", "", -1)
	}
}
