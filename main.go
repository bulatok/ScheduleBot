package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	ss "scheduleBot/service"
	"time"
)

func main(){
	groups := ss.CreateGroups()
	grps := ss.Npm // просто массив строк из названий групп

	// данные по умолчанию-----------------------------------
	UserGroup := "мпм121" // тут хранится группа для юзера (по умолчанию мпм121)
	days := ss.CreateDaysImplement(groups["мпм121"]) // по умолчанию неделя для группы мпм121

	// -------------------------
	b, err := tb.NewBot(tb.Settings{
		Token: ss.Token,
		Poller: &tb.LongPoller{Timeout: 10*time.Second},
	})
	if err != nil{
		log.Fatal(err)
	}

	// кнопки для вывода расписания---------------------------
	r2 := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnNow := r2.Text("СЕЙЧАС❗️")
	btnDay := r2.Text("ДЕНЬ⌛️")
	r2.Reply(
		r2.Row(btnNow),
		r2.Row(btnDay),
	)
	// кнопки для выбора дня----------------------------------
	r3 := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	weekDays := map[string]int{"ПНД" : 0, "ВТР" : 1, "СРД" : 2, "ЧТВ" : 3, "ПТН" : 4, "СУБ" : 5}
	btnWeekDays := make([]tb.Btn, len(weekDays))
	for name, num := range weekDays{
		btnWeekDays[num] = r3.Text(name)
	}
	r3.Reply(
		r3.Row(btnWeekDays...),
	)
	// кнопки для выбора групппы------------------------------
	r1 := &tb.ReplyMarkup{}
	btnGroups := make([]tb.Btn, len(grps))
	for i, v := range grps{
		btnGroups[i] = r1.Text(v)
	}
	r1.Reply(
		r1.Row(btnGroups...),
	)



	// Сохранение данных---------------------------------------
	b.Handle("/edit", func (m *tb.Message){
		b.Send(m.Sender, "Какая группа?", r1)
	})
	for _, someBtn := range btnGroups {
		b.Handle(&someBtn, func(m *tb.Message) {
			UserGroup = m.Text
			days = ss.CreateDaysImplement(groups[UserGroup])
			b.Send(m.Sender, "Ваши данные успешно сохранены "+UserGroup, r2)
		})
	}
	// Вывод данных---------------------------------------------
	b.Handle(&btnNow, func (m *tb.Message){
		idx := int(time.Now().Weekday()) - 1
		if idx == -1{ // если воскресенье
			b.Send(m.Sender, "🦍чил🦍", r2)
			return
		}
		curday := days[idx]
		b.Send(m.Sender, curday.PrettyDayWithTimer(), r2)
	})
	b.Handle(&btnDay, func (m *tb.Message){
		b.Send(m.Sender, "Выберете день", r3)
	})
	for _, someBtn := range btnWeekDays{
		b.Handle(&someBtn, func(m *tb.Message){
			curday := days[weekDays[m.Text]]
			b.Send(m.Sender, curday.PrettyDay(), r2)
		})
	}
	b.Start()
}