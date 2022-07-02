package main

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	models "scheduleBot/models"
	service "scheduleBot/service"
	"time"
)

var (
	grps = []string{"B21-01","B21-02","B21-03","B21-04","B21-05","B21-06","B21-07","B21-08"}
	weekDays = map[string]int{"ПНД" : 0, "ВТР" : 1, "СРД" : 2, "ЧТВ" : 3, "ПТН" : 4}
)
type connection struct{
	bot       *tele.Bot
	btns      map[string]*tele.ReplyMarkup
	btnsText  map[string][]tele.Btn
	Group 	  models.Group
}

func newConnection(botSettings tele.Settings) (*connection, error){
	bot, err := tele.NewBot(botSettings)
	if err != nil{
		return nil, err
	}
	grp := models.Group{
		Name: "B21-07",
		Week: service.CreateWeek("BS1 ", "B21-07"),
	}

	return &connection{
		bot: bot,
		Group: grp,
		btns: make(map[string]*tele.ReplyMarkup),
		btnsText: make(map[string][]tele.Btn),
	}, nil
}
func (conn *connection) createButtons(buttonName string, needResize bool, buttonTexts ...string){
	conn.btns[buttonName] = &tele.ReplyMarkup{ResizeKeyboard: needResize}

	var rows []tele.Row

	for _, buttonText := range buttonTexts{
		curButton := conn.btns[buttonName].Text(buttonText)
		conn.btnsText[buttonName] = append(conn.btnsText[buttonName], curButton)
		rows = append(rows, conn.btns[buttonName].Row(curButton))
	}
	conn.btns[buttonName].Reply(
		rows...
		)
}
func (conn *connection)  buttonsSetup(){
	// creating buttons #1
	conn.createButtons("r1", false, grps...)

	// creating buttons #2
	conn.createButtons("r2", true,"СЕЙЧАС❗️", "ДЕНЬ⌛️")

	// creating buttons #3
	conn.createButtons("r3", true,"ПНД", "ВТР", "СРД", "ЧТВ", "ПТН")
}
func (conn *connection) HandleStart(){
	conn.bot.Handle("/start", func(context tele.Context) error {
		return context.Send("Hello! choose group below", conn.btns["r1"])
	})
}
func (conn *connection) HandleEdit(){
	conn.bot.Handle("/edit", func(context tele.Context) error {
		return context.Send("Какая группа?", conn.btns["r1"])
	})
}
func (conn *connection) HandleGroupButton(){
	for _, someBtn := range conn.btnsText["r1"]{
		conn.bot.Handle(&someBtn, func(context tele.Context) error {
			conn.Group.Name = context.Text()
			conn.Group.Week = service.CreateWeek("BS1 ", conn.Group.Name)
			return context.Send("Ваши данные успешно сохранены "+conn.Group.Name, conn.btns["r2"])
		})
	}
}
func (conn *connection) HandleNowDayButton(){
	for _, someBtn := range conn.btnsText["r2"]{
		switch someBtn.Text {
		case "СЕЙЧАС❗️":
			conn.bot.Handle(&someBtn, func(context tele.Context) error {
				idx := int(time.Now().Weekday()) - 1
				if idx == -1 || len(conn.Group.Week) <= idx{ // если воскресенье
					return context.Send("🦍чил🦍", conn.btns["r2"])
				}
				curday := conn.Group.Week[idx]
				return context.Send(curday.PrettyWithTimer(), conn.btns["r2"])
			})
		case "ДЕНЬ⌛️":
			conn.bot.Handle(&someBtn, func(context tele.Context) error {
				return context.Send("Выберете день", conn.btns["r3"])
			})
		}
	}

}
func (conn *connection) HandleWeekDayButton(){
	for _, someBtn := range conn.btnsText["r3"]{
		conn.bot.Handle(&someBtn, func(context tele.Context) error {
			curday := conn.Group.Week[weekDays[context.Text()]]
			return context.Send(curday.PrettyDay(), conn.btns["r2"])
		})
	}
}
func main(){
	token := os.Getenv("TOKEN")
	if token == ""{
		log.Fatal("need to set TOKEN")
	}
	
	conn, err := newConnection(tele.Settings{
		Token: token,
		Poller: &tele.LongPoller{Timeout: 10*time.Second},
	})
	if err != nil{
		log.Fatal(err)
	}
	conn.buttonsSetup()

	conn.HandleStart()
	conn.HandleEdit()

	conn.HandleGroupButton()
	conn.HandleWeekDayButton()
	conn.HandleNowDayButton()

	log.Println("started running a bot")
	conn.bot.Start()
}