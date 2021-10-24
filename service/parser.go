package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"strings"
	"time"
)
var Npm = [...]string{"мм121", "ммр121", "мак121", "му121", "мх121", "мв121",
"мв221", "мв321", "мв421", "ми121", "ми221", "ми321", "мпм121"}
//const sheet = `Лист1`
type Para struct{
	Nomer string
	Room string
	Vid string
	Teach string
	Name string
	Vremya string
}
type Day struct{
	Weekday int
	Pars []Para
}
type TwoWeeks struct{
	WeekEven []Day
	WeekOdd []Day
}
func CreateGroups() (map[string]TwoWeeks){
	res := make(map[string]TwoWeeks)
	for _, v := range Npm{
		wOdd := CreateOddWeek(v)
		wEven := CreateOddWeek(v)
		res[v] = TwoWeeks{
			WeekOdd: wOdd,
			WeekEven: wEven,
		}
	}
	return res
}
func (p *Para) prettyPara(){
	fmt.Println(p.Nomer, p.Vremya, p.Name)
}
func CreateDaysImplement(weeks TwoWeeks) []Day{
	_, v := time.Now().ISOWeek()
	if v % 2 == 0{
		return weeks.WeekEven
	}else{
		return weeks.WeekOdd
	}
}
func CreateOddWeek(sheet string) ([]Day){
	f, err := excelize.OpenFile("service/forParseFull.xlsx")
	if err != nil{
		log.Fatal(err)
	}
	// A день недели
	// B номер аудитории
	// C время пары
	// D аудитория
	// E вид занятия
	// F препод
	// G имя пары
	b := [...]string{"A", "B", "C", "D", "E", "F", "G"}
	days := make([]Day, 6)
	thisDay := &Day{}
	curDay := 0
	for i := 16; i <= 92; i++{
		// если в это время вообще нет пары
		testZapros1 := "D" + strconv.FormatInt(int64(i), 10)
		cell, err := f.GetCellValue(sheet, testZapros1)
		if err != nil{
			log.Fatal(err)
		}
		if len(cell) == 0{
			continue
		}

		// если новый день
		dayName, err := f.GetCellValue(sheet, "A" + strconv.FormatInt(int64(i), 10))
		if err != nil{
			log.Fatal(err)
		}
		if strings.Contains(dayName, newDay){
			days[curDay] = *thisDay
			thisDay = &Day{}
			curDay += 1
			continue
		}
		// читает пару
		newPara := &Para{}
		for _, v := range b{
			axis := v + strconv.FormatInt(int64(i), 10)
			cell, err := f.GetCellValue(sheet, axis)
			if err != nil{
				log.Fatal(err)
			}
			switch v {
			case "B":
				newPara.Nomer = cell
			case "C":
				newPara.Vremya = cell
			case "D":
				newPara.Room = cell
			case "E":
				newPara.Vid = cell
			case "F":
				newPara.Teach = cell
			case "G":
				newPara.Name = cell
			}
		}
		// записывает пару в день
		thisDay.Weekday = curDay
		thisDay.Pars = append(thisDay.Pars, *newPara)
	}
	return days
}
func CreateEvenWeek(sheet string) ([]Day){
	f, err := excelize.OpenFile("service/forParseFull.xlsx")
	if err != nil{
		log.Fatal(err)
	}
	// A день недели
	// B номер аудитории
	// C время пары
	// D аудитория
	// E вид занятия
	// F препод
	// G имя пары
	b := [...]string{"H", "I", "J", "K", "L", "M"}
	days := make([]Day, 6)
	thisDay := &Day{}
	curDay := 0
	for i := 16; i <= 92; i++{
		// если в это время вообще нет пары
		testZapros1 := "D" + strconv.FormatInt(int64(i), 10)
		cell, err := f.GetCellValue(sheet, testZapros1)
		if err != nil{
			log.Fatal(err)
		}
		if len(cell) == 0{
			continue
		}

		// если новый день
		dayName, err := f.GetCellValue(sheet, "A" + strconv.FormatInt(int64(i), 10))
		if err != nil{
			log.Fatal(err)
		}
		if strings.Contains(dayName, newDay){
			days[curDay] = *thisDay
			thisDay = &Day{}
			curDay += 1
			continue
		}
		// читает пару
		newPara := &Para{}
		for _, v := range b{
			axis := v + strconv.FormatInt(int64(i), 10)
			cell, err := f.GetCellValue(sheet, axis)
			if err != nil{
				log.Fatal(err)
			}
			switch v {
			case "L":
				newPara.Nomer = cell
			case "M":
				newPara.Vremya = cell
			case "K":
				newPara.Room = cell
			case "J":
				newPara.Vid = cell
			case "I":
				newPara.Teach = cell
			case "H":
				newPara.Name = cell
			}
		}
		// записывает пару в день
		thisDay.Weekday = curDay
		thisDay.Pars = append(thisDay.Pars, *newPara)
	}
	return days
}
func PrettyDays(d []Day){
	for _, v := range d{
		fmt.Println(v.Weekday)
		for _, v2 := range v.Pars{
			fmt.Println(v2)
		}
	}
}
func noParsWithClock(d *Day) bool{
	curr := time.Now()
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

func (d *Day) PrettyDayWithTimer() string{
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
		timeNow := time.Now()
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
			res += fmt.Sprintf("%s\n👨‍🏫 %s\n🕐 %s\n🚪%s\n🏛%s\n⏸️%dЧ %dM\n\n",
				v.Name, v.Teach, v.Vremya, v.Room, hh, mm)
		}else if timeNow.Hour() <= starts.Hour() && timeNow.Minute() <= starts.Minute(){
			hh, mm := timeRazn(timeNow, starts)
			res += fmt.Sprintf("%s\n👨‍🏫 %s\n🕐 %s\n🚪%s\n🏛%s\n▶️%dh %dm\n\n",
				v.Name, v.Teach, v.Vremya, v.Room, v.Vid, hh, mm)
		}
		if cnt <= 0{
			break
		}
		cnt -= 1
	}
	return res
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
		res += fmt.Sprintf("%s\n👨‍🏫 %s\n🕐 %s\n🚪%s\n🏛%s\n\n",
			v.Name, v.Teach, v.Vremya, v.Room, v.Vid)
	}
	return res
}
