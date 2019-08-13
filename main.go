package main

import (
	"fmt"
	"github.com/enjuus/oriebot/models"
	"log"
	"time"
	"github.com/joshbetz/config"
	tb "github.com/tucnak/telebot"
)

type Env struct {
	db models.Datastore
	bot  *tb.Bot
}

var TGToken string

func main() {
	c := config.New(".config.json")
	c.Get("TGToken", &TGToken)

	b, err := tb.NewBot(tb.Settings{
		Token:  TGToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	db, err := models.NewDB("bot.db")
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db, b}

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("running")

	b.Handle("/chat", env.HandleChatID)
	b.Handle("/quote", env.HandleQuotes)

	b.Start()
}
