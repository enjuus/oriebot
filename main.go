package main

import (
	"fmt"
	"github.com/enjuus/oriebot/models"
	tb "github.com/tucnak/telebot"
	"log"
	"time"
)

type Env struct {
	db             models.Datastore
	bot            *tb.Bot
	LastFMAPIKey   string
	LastFMSecret   string
	OpenWeatherAPI string
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  TGToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	db, err := models.NewDB("bot.db")
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db, b, LastFMAPIKey, LastFMSecret, OpenWeatherAPI}

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("running")

	b.Handle("/chat", env.HandleChatID)
	b.Handle("/quote", env.HandleQuotes)
	b.Handle("/lastfm", env.HandleLastFM)
	b.Handle("/weather", env.HandleWeather)
	b.Handle("/uwu", env.HandleUWU)
	b.Handle("/spurdo", env.HandleSpurdo)
	b.Start()
}
