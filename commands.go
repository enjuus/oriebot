package main

import (
	"bytes"
	"fmt"
	"github.com/ndyakov/go-lastfm"
	tb "github.com/tucnak/telebot"
	"log"
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
				fmt.Println(err)
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
