package main

import (
	"bytes"
	"fmt"
	"github.com/briandowns/openweathermap"
	"github.com/dafanasev/go-yandex-translate"
	"github.com/enjuus/go-collage"
	"github.com/enjuus/spurdo"
	"github.com/enjuus/uwu"
	"github.com/ndyakov/go-lastfm"
	tb "github.com/tucnak/telebot"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func (env *Env) HandleChatID(m *tb.Message) {
	_, err := env.bot.Send(m.Chat, fmt.Sprintf("ChatID: %d", m.Chat.ID))
	if err != nil {
		log.Panic(err)
	}
}

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
		ID := strings.Replace(m.Text, "/quote ", "", 1)
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

func (env *Env) HandleLastFMTopAlbums(m *tb.Message) {
	lf, err := env.db.GetLastFM(m.Sender.ID)
	if lf == nil {
		env.bot.Send(m.Chat, fmt.Sprintf("No User set, set it with /lastfm"))
		return
	}
	lfm := lastfm.New(env.LastFMAPIKey, env.LastFMSecret)
	response, err := lfm.User.GetTopAlbums(lf.LastfmName, "3month", 0, 9)
	if err != nil {
		env.bot.Send(m.Chat, fmt.Sprintf("i pooped and shidded"))
	}
	images := make(map[int]string)
	for i, element := range response.TopAlbums {
		resp, err := http.Get(element.Image[3].URL)
		if err != nil {
			env.bot.Send(m.Chat, fmt.Sprintf("i pooped and shidded"))
		}
		defer resp.Body.Close()
		path := "/tmp/" + path.Base(element.Image[3].URL)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			file, _ := os.Create(path)
			defer file.Close()
			io.Copy(file, resp.Body)
		}
		images[i] = path
	}

	files, err := collage.MapImages(images)
	if err != nil {
		env.bot.Send(m.Chat, fmt.Sprintf("i pooped and shidded"))
	}
	err = collage.MakeNewCollage(files, "/tmp/collage.jpg", 100)
	if err != nil {
		env.bot.Send(m.Chat, fmt.Sprintf("i pooped and shidded"))
	}

	photo := &tb.Photo{File: tb.FromDisk("/tmp/collage.jpg")}
	env.bot.Send(m.Chat, photo)

}

func (env *Env) HandleLastFM(m *tb.Message) {
	user := strings.Replace(m.Text, "/lastfm ", "", 1)
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

func (env *Env) HandleWeather(m *tb.Message) {
	addr := strings.Replace(m.Text, "/weather ", "", 1)
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

func (env *Env) HandleUWU(m *tb.Message) {
	var text string
	if m.ReplyTo != nil {
		text = m.ReplyTo.Text
	} else {
		text = strings.Replace(m.Text, "/uwu ", "", 1)
	}
	str, err := uwu.Translate(text)
	if err != nil {
		str = "Message can't be empty"
	}
	_, err = env.bot.Send(m.Chat, str)
}

func (env *Env) HandleSpurdo(m *tb.Message) {
	var text string
	if m.ReplyTo != nil {
		text = m.ReplyTo.Text
	} else {
		text = strings.Replace(m.Text, "/spurdo ", "", 1)
	}
	str, err := spurdo.Translate(text)
	if err != nil {
		str = "Message can't be empty"
	}
	_, err = env.bot.Send(m.Chat, str)
}

func (env *Env) HandleBlog(m *tb.Message) {
	_, err := env.bot.Send(m.Chat, "Nobody fucking cares, dude")
	if err != nil {
		return
	}
}

func (env *Env) HandleTranslate(m *tb.Message) {
	var text string
	tr := translate.New(env.YandexAPI)
	if m.ReplyTo != nil {
		text = m.ReplyTo.Text
	} else {
		text = strings.Replace(m.Text, "/tl ", "", 1)
	}
	translation, err := tr.Translate("en", text)
	if err != nil {
		fmt.Println(err)
	} else {
		_, _ = env.bot.Send(m.Chat, translation.Result())
	}
}
