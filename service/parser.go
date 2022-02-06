package service

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)
type Para struct{
	Nomer  string `json:"nomer"`
	Room   string `json:"room"`
	Teach  string `json:"teach"`
	Name   string `json:"name"`
	Vremya string `json:"vremya"`
}
type Day struct{
	Weekday int
	Pars    []Para
}
func GetToken() string{
	return os.Getenv("TOKEN")
}
func GetTimezone(timezn string) (*time.Location, error){
	loc, err := time.LoadLocation(timezn)
	if err != nil{
		return nil, err
	}
	return loc, nil
}
func CreateWeek(sheet string, group string) []Day{
	f, err := excelize.OpenFile("internal/inno.xlsx")
	if err != nil{
		log.Fatal(err)
	}
	weekdays := [...]string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY"}

	var days []Day
	var pars []Para
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
			days = append(days, Day{Weekday: curDay, Pars: pars})
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
		pars = append(pars, Para{
			Name: lesson,
			Room: room,
			Teach: teacher,
			Vremya: vremya,
		})
		i+=3
	}
	return days
}
func (d *Day) PrettyDay() string{
	res := ""
	if len(d.Pars) == 0{
		res = "🦍чил🦍"
		return res
	}
	for _, v := range d.Pars{
		if len(v.Name) == 0{
			continue
		}
		res += fmt.Sprintf("%s\n👨‍🏫 %s\n🕐 %s\n🚪%s\n\n",
			v.Name, v.Teach, v.Vremya, v.Room)
	}
	return res
}
func noParsWithClock(d *Day) bool{
	loc, _ := GetTimezone("Europe/Moscow")
	curr := time.Now().In(loc)
	if len(d.Pars) == 0{
		return true
	}
	paraTime := strings.Split(d.Pars[len(d.Pars) - 1].Vremya, "-")

	//starts, err := time.Parse("15-04-05", paraTime[0])
	//if err != nil{
	//	log.Fatal(err)
	//}
	ends, err := time.Parse("15:04", paraTime[1])
	if err != nil{
		log.Fatal(err)
	}
	if curr.Hour() > ends.Hour(){
		return true
	}
	if curr.Minute() > ends.Minute(){
		return true
	}
	return false

}
func timeRazn(t1 time.Time, t2 time.Time) (int, int){
	return t2.Hour() - t1.Hour(), t2.Minute() - t1.Minute()
}
func (d *Day) PrettyWithTimer() string{
	res := ""
	if noParsWithClock(d) == true{
		return "🦍чил🦍"
	}
	// сколько пар(уроков) вывести
	cnt := 4
	for _, v := range d.Pars{
		if len(v.Name) == 0{
			continue
		}
		loc, _ := GetTimezone("Europe/Moscow")

		timeNow := time.Now().In(loc)

		//timeNow, err := time.Parse("15:04", "09:00") //для теста(ставит время 09:00)
		//if err != nil{
		//	log.Fatal(err)
		//}

		paraTime := strings.Split(v.Vremya, "-")

		starts, err := time.Parse("15:04", paraTime[0])
		if err != nil{
			log.Fatal(err)
		}
		ends, err := time.Parse("15:04", paraTime[1])
		if err != nil{
			log.Fatal(err)
		}
		// значит сейчас идет эта пара
		if timeNow.Hour() >= starts.Hour() && timeNow.Minute() >= starts.Minute() && timeNow.Hour() <= ends.Hour() && timeNow.Minute() <= ends.Minute(){
			hh, mm := timeRazn(timeNow, ends)
			res += fmt.Sprintf("%s\n👨‍🏫 %s\n🕐 %s\n🚪%s\n⏸️%dh %dm\n\n",
				v.Name, v.Teach, v.Vremya, v.Room, hh, mm)
		}else if timeNow.Hour() <= starts.Hour() && timeNow.Minute() <= starts.Minute(){
			hh, mm := timeRazn(timeNow, starts)
			res += fmt.Sprintf("%s\n👨‍🏫 %s\n🕐 %s\n🚪%s\n▶️%dh %dm\n\n",
				v.Name, v.Teach, v.Vremya, v.Room, hh, mm)
		}
		if cnt <= 0{
			break
		}
		cnt -= 1
	}
	return res
}

func GetJSON(days []Day) ([]byte){
	data ,_ := json.Marshal(days)
	return data
}