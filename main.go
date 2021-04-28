package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enjuus/oriebot/models"
	tb "github.com/tucnak/telebot"
)

// Env is the main struct being passed ot all commands
type Env struct {
	db             models.Datastore
	bot            *tb.Bot
	LastFMAPIKey   string
	LastFMSecret   string
	OpenWeatherAPI string
	YandexAPI      string
	ListOfAuth     [2]int64
	Ticker         *time.Ticker
	QuitCall       chan struct{}
	MainChannel    *tb.Chat
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

	ticker := time.NewTicker(60 * time.Minute)
	quit := make(chan struct{})

	env := &Env{db, b, LastFMAPIKey, LastFMSecret, OpenWeatherAPI, YandexAPI, listOfAuth, ticker, quit, &tb.Chat{}}

	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("running")

	b.Handle("/chat", env.HandleChatID)
	b.Handle("/quote", env.HandleQuotes)
	b.Handle("/lastfm", env.HandleLastFM)
	b.Handle("/topalbums", env.HandleLastFMTopAlbums)
	b.Handle("/weather", env.HandleWeather)
	b.Handle("/uwu", env.HandleUWU)
	b.Handle("/spurdo", env.HandleSpurdo)
	b.Handle("/blog", env.HandleBlog)
	b.Handle("/tl", env.HandleTranslate)
	b.Handle("/decide", env.HandleDecide)
	b.Handle("/turnips", env.HandleTurnips)
	b.Handle("/terms", env.HandleTerms)
	b.Handle("/term", env.HandleTerm)
	b.Handle("/helth", env.HandleHelth)
	b.Handle("/unhelth", env.HandleNoMoreHelth)
	b.Handle("/starthelth", env.HandleStartHelth)
	b.Handle("/stophelth", env.HandleStopHelth)
	b.Handle(tb.OnText, env.HandleTermCount)
	b.Start()

}
