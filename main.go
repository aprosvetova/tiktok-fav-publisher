package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

var cfg config
var r *redis.Client
var tg *tb.Bot

func main() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Config", err)
	}

	r = redis.NewClient(&redis.Options{
		Addr: cfg.DBAddr,
	})

	var err error

	tg, err = tb.NewBot(tb.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatalln("Telegram", err)
	}

	cfg.TikTokSecUserID, err = getSecUserID(cfg.TikTokUsername)
	if err != nil {
		log.Fatalln("SecUID", err)
	}

	log.Println("Polling...")

	for {
		checkNewVideos()
		time.Sleep(time.Minute)
	}
}

func checkNewVideos() {
	likes, err := getLikedVideos(cfg.TikTokSecUserID, 10)
	if err != nil {
		log.Println("Likes", err)
		return
	}

	for _, v := range likes {
		if wasAlreadyPosted(v.ID) {
			continue
		}

		log.Println("Posting", v.ID)

		menu := &tb.ReplyMarkup{}
		menu.Inline(
			menu.Row(menu.URL("Оригинал", v.ShareURL)),
		)

		_, err = tg.Send(tb.ChatID(cfg.ChannelID), &tb.Video{
			File: tb.File{FileURL: v.DownloadURL},
		}, menu)
		if err != nil {
			log.Println("Send video", err, v.DownloadURL)
		}

		time.Sleep(time.Second * 3)
	}
}
