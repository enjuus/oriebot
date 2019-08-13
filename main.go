package main

import (
	"bytes"
	"fmt"
	"github.com/enjuus/oriebot/models"
	"log"
	"strings"
	"time"

	tb "github.com/tucnak/telebot"
)

type Env struct {
	db models.Datastore
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  "",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	db, err := models.NewDB("bot.db")
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db}

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/chatid", func(m *tb.Message) {
		b.Send(m.Chat, fmt.Sprintf("ChatID: %d", m.Chat.ID))
	})

	b.Handle("/quote", func(m *tb.Message) {
		if m.ReplyTo != nil {
			err := env.db.AddQuote(m.ReplyTo.Text, m.ReplyTo.Sender.Username, m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName, m.ReplyTo.Sender.ID)
			if err != nil {
				log.Panic(err)
			}
			b.Send(m.Chat, "Added quote!")
		} else {
			ID := strings.Replace(m.Text, "/quote ", "", 1)
			if ID != "/quote" && ID != "all" {
				quote, err := env.db.GetQuote(ID)
				if err != nil {
					b.Send(m.Chat, "That quote doesn't exist")
				}
				str := fmt.Sprintf("*%s* \n\n- _%s_", quote.Message, quote.Sender)
				b.Send(m.Chat, str, tb.ParseMode("Markdown"))
			} else if ID == "all" {
				quotes,err := env.db.AllQuotes()
				if err != nil {
					b.Send(m.Chat, "There are no quotes")
				}
				var str bytes.Buffer
				for _, qt := range quotes {
					quote := fmt.Sprintf("%d: *%s* - _%s_\n", qt.ID, qt.Message, qt.Sender)
					str.WriteString(quote)
				}
				b.Send(m.Sender, str.String(), tb.ParseMode("Markdown"))
			}
		}
	})

	b.Start()
}
