package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Para struct {
	Nomer  string `json:"number"`
	Room   string `json:"room"`
	Teach  string `json:"teacher"`
	Name   string `json:"name"`
	Vremya string `json:"time"`
}
type Day struct {
	Weekday int    `json:"weekday"`
	Pars    []Para `json:"pars"`
}
type Group struct {
	Name string `json:"group_name"`
	Week []Day  `json:"week"`
}

func (d Day) PrettyWithTimer() string{
	res := ""
	if d.noParsWithClock() == true{
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
			res += fmt.Sprintf("%s\n👨‍🏫 %s\n🕐 %s\n🚪%s\n⏸️%dЧ %dM\n\n",
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

func GetTimezone(timezn string) (*time.Location, error){
	loc, err := time.LoadLocation(timezn)
	if err != nil{
		return nil, err
	}
	return loc, nil
}
func (d Day) PrettyDay() string{
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
func (d Day) noParsWithClock() bool{
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

// JSON возвращает объект m.Group в формате JSON
func (dd Group) JSON() ([]byte, error){
	data, err := json.Marshal(dd)
	if err != nil{
		return nil, err
	}
	return data, nil
}