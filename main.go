package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	ss "scheduleBot/service"
	"time"
)

var (
	grps = []string{"B21-01","B21-02","B21-03","B21-04","B21-05","B21-06","B21-07","B21-08"}
	weekDays = map[string]int{"ПНД" : 0, "ВТР" : 1, "СРД" : 2, "ЧТВ" : 3, "ПТН" : 4}
)


type conn struct{
	bot       *tb.Bot
	btns      map[string]*tb.ReplyMarkup
	btnsText  map[string][]tb.Btn
	userGroup string
	days      []ss.Day
	port 	  string
}
func newConn(botSettings tb.Settings) (*conn, error){
	bot, err := tb.NewBot(botSettings)
	if err != nil{
		return nil, err
	}
	return &conn{
		bot: bot,
		days: ss.CreateWeek("BS1 ", "B21-07"),
		btns: make(map[string]*tb.ReplyMarkup),
		btnsText: make(map[string][]tb.Btn),
		userGroup: "B21-01",
	}, nil
}
func (c *conn) createButtons(buttonName string, needResize bool, buttonTexts ...string){
	c.btns[buttonName] = &tb.ReplyMarkup{ResizeReplyKeyboard: needResize}

	var rows []tb.Row

	for _, buttonText := range buttonTexts{
		curButton := c.btns[buttonName].Text(buttonText)
		c.btnsText[buttonName] = append(c.btnsText[buttonName], curButton)
		rows = append(rows, c.btns[buttonName].Row(curButton))
	}
	c.btns[buttonName].Reply(
		rows...
	)
}
func (c *conn)  buttonsSetup(){
	// creating buttons #1
	c.createButtons("r1", false, grps...)

	// creating buttons #2
	c.createButtons("r2", true,"СЕЙЧАС❗️", "ДЕНЬ⌛️")

	// creating buttons #3
	c.createButtons("r3", true,"ПНД", "ВТР", "СРД", "ЧТВ", "ПТН")


}
func (c *conn) HandleStart(){
	c.bot.Handle("/start", func(m *tb.Message) {
		c.bot.Send(m.Sender, "Hello! choose group below", c.btns["r1"])
	})
}
func (c *conn) HandleEdit(){
	c.bot.Handle("/edit", func(m *tb.Message) {
		c.bot.Send(m.Sender, "Какая группа?", c.btns["r1"])
	})
}
func (c *conn) HandleGroupButton(){
	for _, someBtn := range c.btnsText["r1"]{
		c.bot.Handle(&someBtn, func(m *tb.Message) {
			c.userGroup = m.Text
			c.days = ss.CreateWeek("BS1 ", c.userGroup)
			c.bot.Send(m.Sender, "Ваши данные успешно сохранены "+c.userGroup, c.btns["r2"])
		})
	}
}
func (c *conn) HandleNowDayButton(){
	for _, someBtn := range c.btnsText["r2"]{
		switch someBtn.Text {
		case "СЕЙЧАС❗️":
			c.bot.Handle(&someBtn, func(m *tb.Message) {
				idx := int(time.Now().Weekday()) - 1
				if idx == -1 { // если воскресенье
					c.bot.Send(m.Sender, "🦍чил🦍", c.btns["r2"])
					return
				}
				curday := c.days[idx]
				c.bot.Send(m.Sender, curday.PrettyWithTimer(), c.btns["r2"])
			})
		case "ДЕНЬ⌛️":
			c.bot.Handle(&someBtn, func (m *tb.Message){
				c.bot.Send(m.Sender, "Выберете день", c.btns["r3"])
			})
		}
	}

}
func (c *conn) HandleWeekDayButton(){
	for _, someBtn := range c.btnsText["r3"]{
		c.bot.Handle(&someBtn, func(m *tb.Message){
			curday := c.days[weekDays[m.Text]]
			c.bot.Send(m.Sender, curday.PrettyDay(), c.btns["r2"])
		})
	}
}
func main(){

	conn1, err := newConn(tb.Settings{
		Token: os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: time.Minute*5},
	})

	if err != nil{
		log.Fatal(err)
	}

	conn1.buttonsSetup()

	conn1.HandleStart()
	conn1.HandleEdit()

	conn1.HandleGroupButton()
	conn1.HandleWeekDayButton()
	conn1.HandleNowDayButton()


	log.Println("started running a bot")
	conn1.bot.Start()
}