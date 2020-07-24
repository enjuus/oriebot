package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/openweathermap"
	"github.com/dafanasev/go-yandex-translate"
	"github.com/enjuus/go-collage"
	"github.com/enjuus/spurdo"
	"github.com/enjuus/uwu"
	"github.com/ndyakov/go-lastfm"
	tb "github.com/tucnak/telebot"
)

// HandleChatID returns the senders ChatID / GroupID
func (env *Env) HandleChatID(m *tb.Message) {
	_, err := env.bot.Send(m.Chat, fmt.Sprintf("ChatID: %d", m.Chat.ID))
	if err != nil {
		log.Panic(err)
	}
}

func (env *Env) HandleCommandAddr(command string, text string) string {
	if strings.Contains(text, "@oriebot") {
		command = command + "@oriebot"
	}
	addr := strings.Replace(text, command, "", 1)
	return addr
}

func (env *Env) CheckOptions(arr [6]string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (env *Env) Reverse(Input string) string {
	n := 0
	rune := make([]rune, len(Input))
	for _, r := range Input {
		rune[n] = r
		n++
	}
	rune = rune[0:n]
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}
	output := string(rune)
	return output
}

// HandleQuotes stores and retrieves quotes from the database
func (env *Env) HandleQuotes(m *tb.Message) {
	if m.ReplyTo != nil {
		err := env.db.AddQuote(m.ReplyTo.Text, m.ReplyTo.Sender.Username, m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName, m.ReplyTo.Sender.ID)
		if err != nil {
			return
		}
		_, err = env.bot.Send(m.Chat, "Added quote!")
		if err != nil {
			return
		}
	} else {
		ID := env.HandleCommandAddr("/quote", m.Text)
		if ID != "/quote" && ID != "all" {
			quote, err := env.db.GetQuote(ID)
			fmt.Printf("%v\n", err)
			if err != nil {
				_, err = env.bot.Send(m.Chat, "That quote doesn't exist")
				if err != nil {
					return
				}
				return
			}
			str := fmt.Sprintf("*%s* \n\n- _%s_", quote.Message, quote.Sender)
			_, err = env.bot.Send(m.Chat, str, tb.ParseMode("Markdown"))
			if err != nil {
				log.Panic(err)
			}
		} else if ID == "all" {
			quotes, err := env.db.AllQuotes()
			if err != nil {
				_, err = env.bot.Send(m.Chat, "There are no quotes")
				if err != nil {
					return
				}
				return
			}
			var str bytes.Buffer
			for _, qt := range quotes {
				quote := fmt.Sprintf("%d: *%s* - _%s_\n", qt.ID, qt.Message, qt.Sender)
				str.WriteString(quote)
			}
			_, err = env.bot.Send(m.Sender, str.String(), tb.ParseMode("Markdown"))
			if err != nil {
				return
			}
		}
	}
}

// HandleLastFMTopAlbums generates a collage of the top albums of the requested last.fm user
func (env *Env) HandleLastFMTopAlbums(m *tb.Message) {
	lf, err := env.db.GetLastFM(m.Sender.ID)
	options := [6]string{"overall", "7day", "1month", "3month", "6month", "12month"}
	period := env.HandleCommandAddr("/topalbums", m.Text)
	if env.CheckOptions(options, period) {
		period = "3month"
	}
	baseDir, err := os.UserHomeDir()
	if err != nil {
		env.bot.Send(m.Chat, fmt.Sprintf("i pooped and shidded"))
	}
	folder := baseDir + "/npimg/"
	if lf == nil {
		env.bot.Send(m.Chat, fmt.Sprintf("No User set, set it with /lastfm"))
		return
	}
	lfm := lastfm.New(env.LastFMAPIKey, env.LastFMSecret)
	response, err := lfm.User.GetTopAlbums(lf.LastfmName, period, 0, 9)
	if err != nil {
		env.bot.Send(m.Chat, fmt.Sprintf("i pooped and shidded"))
	}
	images := make(map[int]string)
	for i, element := range response.TopAlbums {
		resp, err := http.Get(element.Image[3].URL)
		if err == nil {
			defer resp.Body.Close()
			images[i] = folder + path.Base(element.Image[3].URL)
			if _, err := os.Stat(images[i]); os.IsNotExist(err) {
				file, _ := os.Create(images[i])
				defer file.Close()
				io.Copy(file, resp.Body)
			}
		} else {
			images[i] = baseDir + "/standin.jpg"
		}
	}

	files, err := collage.MapImages(images)
	if err != nil {
		env.bot.Send(m.Chat, fmt.Sprintf("i pooped and shidded"))
	}
	err = collage.MakeNewCollage(files, folder+"collage.jpg", 100)
	if err != nil {
		env.bot.Send(m.Chat, fmt.Sprintf("i pooped and shidded"))
	}

	photo := &tb.Photo{File: tb.FromDisk(folder + "collage.jpg")}
	env.bot.Send(m.Chat, photo)
	os.RemoveAll(folder)
	os.MkdirAll(folder, os.ModePerm)

}

// HandleLastFM stores a new last.fm name for a TG ID and returns the current/last played song
func (env *Env) HandleLastFM(m *tb.Message) {
	user := env.HandleCommandAddr("/lastfm", m.Text)
	if user != "/lastfm" && user != "" {
		lf, err := env.db.GetLastFM(m.Sender.ID)
		if err == nil {
			fmt.Println(err)
		}
		if lf == nil {
			err := env.db.AddLastFM(m.Sender.ID, user)
			if err != nil {
				env.bot.Send(m.Chat, fmt.Sprintf("No User set, set it with /lastfm"))
				return
			}
		} else if lf.LastfmName != user {
			err := env.db.UpdateLastFM(m.Sender.ID, user)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	lf, _ := env.db.GetLastFM(m.Sender.ID)
	if lf != nil {
		var output bytes.Buffer
		lfm := lastfm.New(env.LastFMAPIKey, env.LastFMSecret)
		response, err := lfm.User.GetRecentTracks(lf.LastfmName, 0, 0, 0, 0)
		if err != nil {
			_, err = env.bot.Send(m.Chat, fmt.Sprintf("Error: %s", err))
			return
		}
		track := response.RecentTracks[0]
		if track.NowPlaying != "" {
			string := fmt.Sprintf("*Now playing*\n")
			output.WriteString(string)
		}
		string := fmt.Sprintf("%s - _%s_\n", track.Artist.Name, track.Name)
		output.WriteString(string)

		string = fmt.Sprintf("[​](%s)", track.Image[2].URL)
		output.WriteString(string)
		_, err = env.bot.Send(m.Chat, output.String(), tb.ParseMode("Markdown"))
		if err != nil {
			return
		}
	} else {
		_, err := env.bot.Send(m.Chat, "Please specify a lastfm user")
		if err != nil {
			return
		}
	}
}

// HandleWeather pulls current weather data from the OpenWeatherAPI
// And outputs it back ot the chat
func (env *Env) HandleWeather(m *tb.Message) {
	addr := env.HandleCommandAddr("/weather", m.Text)
	if addr == "/weather" || addr == "" {
		_, err := env.bot.Send(m.Chat, "Please specify a city/country/address")
		if err != nil {
			return
		}
	}
	w, err := openweathermap.NewCurrent("C", "EN", env.OpenWeatherAPI)
	if err != nil {
		_, err = env.bot.Send(m.Chat, fmt.Sprintf("Error: %s", err))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = w.CurrentByName(addr)
	if err != nil {
		_, err = env.bot.Send(m.Chat, fmt.Sprintf("Error: %s", err))
	}
	wtr := fmt.Sprintf("*%s, %s*\n%.2f°C, %s", w.Name, w.Sys.Country, w.Main.Temp, w.Weather[0].Main)
	_, err = env.bot.Send(m.Chat, wtr, tb.ParseMode("Markdown"))
	if err != nil {
		fmt.Println(err)
		return
	}

}

// HandleUWU "translates" a text into uwu
func (env *Env) HandleUWU(m *tb.Message) {
	var text string
	if m.ReplyTo != nil {
		text = m.ReplyTo.Text
	} else {
		text = env.HandleCommandAddr("/uwu", m.Text)
	}
	str, err := uwu.Translate(text)
	if err != nil {
		str = "Message can't be empty"
	}
	_, err = env.bot.Send(m.Chat, str)
}

// HandleSpurdo "translates" a text into spurdo
func (env *Env) HandleSpurdo(m *tb.Message) {
	var text string
	if m.ReplyTo != nil {
		text = m.ReplyTo.Text
	} else {
		text = env.HandleCommandAddr("/spurdo", m.Text)
	}
	str, err := spurdo.Translate(text)
	if err != nil {
		str = "Message can't be empty"
	}
	_, err = env.bot.Send(m.Chat, str)
}

// HandleBlog is being rude
func (env *Env) HandleBlog(m *tb.Message) {
	_, err := env.bot.Send(m.Chat, "Nobody fucking cares, dude")
	if err != nil {
		return
	}
}

// HandleTranslate uses YandexAPI to translate text to english
func (env *Env) HandleTranslate(m *tb.Message) {
	var text string
	tr := translate.New(env.YandexAPI)
	if m.ReplyTo != nil {
		text = m.ReplyTo.Text
	} else {
		text = env.HandleCommandAddr("/tl", m.Text)
	}
	translation, err := tr.Translate("en", text)
	if err != nil {
		fmt.Println(err)
	} else {
		_, _ = env.bot.Send(m.Chat, translation.Result())
	}
}

// HandleDecide takes a string input and helps you decide
func (env *Env) HandleDecide(m *tb.Message) {
	var text string
	if m.ReplyTo != nil {
		text = m.ReplyTo.Text
	} else {
		text = env.HandleCommandAddr("/decide", m.Text)
	}
	split := strings.Split(text, " or ")
	rand.Seed(time.Now().Unix())
	str := fmt.Sprint("", split[rand.Intn(len(split))])
	env.bot.Send(m.Chat, str)
}

func (env *Env) HandleTurnips(m *tb.Message) {
	price := env.HandleCommandAddr("/turnips", m.Text)
	date := time.Now().Format("2006-01-02")
	if price != "/turnips" && price != "" {
		iprice, _ := strconv.Atoi(strings.TrimSpace(price))
		err := env.db.AddTurnip(int(iprice), m.Sender.ID, m.Sender.Username, date)
		if err != nil {
			env.bot.Send(m.Chat, fmt.Sprintf("The price could not be set, yell at enju"))
			return
		}
		env.bot.Send(m.Chat, fmt.Sprintf("The price has been added for today [%s]", date))
	} else {
		turnips, err := env.db.GetTodaysPrices(date)
		if err != nil || len(turnips) == 0 {
			_, err := env.bot.Send(m.Chat, "No prices yet for today. ["+date+"]")
			if err != nil {
				return
			}
		}
		var str bytes.Buffer
		for _, t := range turnips {
			turnip := fmt.Sprintf("*%s* - _%d_ \n", t.Name, t.Price)
			str.WriteString(turnip)
		}
		_, err = env.bot.Send(m.Chat, str.String(), tb.ParseMode("Markdown"))
		if err != nil {
			return
		}
	}
}

func (env *Env) HandleTermCount(m *tb.Message) {
	terms, err := env.db.GetTerms()
	if err != nil {
		_, err = env.bot.Send(m.Chat, "Hm.")
	}
	text := strings.Fields(m.Text)
	for _, word := range text {
		for _, t := range terms {
			if strings.ToLower(word) == t.Name {
				_ = env.db.AddCounter(t.Name)
				term, err := env.db.GetTerm(t.Name)
				fmt.Println(t)
				err = env.db.CountForUser(t.ID, m.Sender.Username)
				if err != nil {
					fmt.Println("count for user", err)
				}
				if err != nil {
					fmt.Println("hm")
				}
				var prev string
				for _, c := range env.Reverse(fmt.Sprintf("%d", term.Count)) {
					if string(c) == prev {
						_, err = env.bot.Send(m.Chat, fmt.Sprintf("%s: %d", term.Name, term.Count))
						return
					}
					prev = string(c)
				}
			}
		}
	}
}

func (env *Env) HandleTerms(m *tb.Message) {
	term := env.HandleCommandAddr("/terms", m.Text)
	if term != "/terms" && term != "" {
		err := env.db.AddTerm(strings.ToLower(strings.TrimSpace(term)))
		if err != nil {
			return
		}
		_, err = env.bot.Send(m.Chat, "Added term")
	} else {
		terms, err := env.db.GetTerms()
		var str bytes.Buffer
		for _, t := range terms {
			term := fmt.Sprintf("*%s* - _%d_ \n", t.Name, t.Count)
			str.WriteString(term)
		}
		_, err = env.bot.Send(m.Chat, str.String(), tb.ParseMode("Markdown"))
		if err != nil {
			return
		}
	}
}

func (env *Env) HandleTerm(m *tb.Message) {
	term := env.HandleCommandAddr("/term", m.Text)
	if term != "/term" && term != "" {
		te, err := env.db.GetTerm(strings.ToLower(strings.TrimSpace(term)))
		if err != nil {
			fmt.Println("err empty")
			return
		}
		userTerms, err := env.db.GetForUsers(te.ID)
		if err != nil {
			fmt.Println(err)
			_, err = env.bot.Send(m.Chat, "I don't know this term.")
		}
		var str bytes.Buffer
		for _, t := range userTerms {
			term := fmt.Sprintf("*%s* - _%d_ \n", t.UserID, t.Count)
			str.WriteString(term)
		}
		_, err = env.bot.Send(m.Chat, str.String(), tb.ParseMode("Markdown"))
		if err != nil {
			return
		}
		env.db.StoreChatLogID(m.Chat.ID)
		for _, x := range env.ListOfAuth {
			if m.Chat.ID != x {
				env.bot.Leave(m.Chat)
			}
		}
	}
}

func (env *Env) DisconnectUnauthorized(m *tb.Message) {
}
