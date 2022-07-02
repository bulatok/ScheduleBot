package internal

import (
	"github.com/xuri/excelize/v2"
	"log"
	m "scheduleBot/models"
	"strconv"
	"strings"
)
func CreateWeek(sheet string, group string) []m.Day{
	f, err := excelize.OpenFile("service/inno.xlsx")
	if err != nil{
		log.Fatal(err)
	}
	weekdays := [...]string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY"}

	var days []m.Day
	var pars []m.Para
	curDay := 0
	b := map[string]string{
		"B21-01":"B",
		"B21-02":"C",
		"B21-03":"D",
		"B21-04":"E",
		"B21-05":"F",
		"B21-06":"G",
		"B21-07":"H",
		"B21-08":"I",
	}

	for i := 4; i <= 98;{
		// chechking either it is new weekday
		newDay := false
		test, err := f.GetCellValue(sheet, "A"+strconv.FormatInt(int64(i), 10))
		if err != nil{
			log.Fatal(err)
		}
		for _, weekDay := range weekdays{
			if weekDay == test{
				newDay = true
				break
			}
		}
		if newDay{
			i += 1
			curDay+=1
			days = append(days, m.Day{Weekday: curDay, Pars: pars})
			pars = nil
		}

		// Implementation...
		lesson, _ := f.GetCellValue(sheet, b[group]+strconv.FormatInt(int64(i), 10))

		teacher, _ := f.GetCellValue(sheet, b[group]+strconv.FormatInt(int64(i+1), 10))

		room, _ := f.GetCellValue(sheet, b[group]+strconv.FormatInt(int64(i+2), 10))
		if strings.HasSuffix(room, ".0"){
			room = room[:len(room)-2]
		}

		vremya, _ := f.GetCellValue(sheet, "A"+strconv.FormatInt(int64(i+2), 10))
		if lesson == ""{
			i += 3
			continue
		}
		pars = append(pars, m.Para{
			Name: lesson,
			Room: room,
			Teach: teacher,
			Vremya: vremya,
		})
		i+=3
	}
	return days
}