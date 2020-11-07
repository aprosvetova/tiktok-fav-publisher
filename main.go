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

	log.Println("Polling...")

	for {
		checkNewVideos()
		time.Sleep(time.Minute)
	}
}

func checkNewVideos() {
	likes, err := getLikes(cfg.TikTokUsername)
	if err != nil {
		log.Println("Likes", err)
		return
	}

	for _, id := range likes {
		if wasAlreadyPosted(id) {
			continue
		}

		link, err := getDownloadURL(id)
		if err != nil {
			log.Println("Download URL", err)
			continue
		}

		log.Println("Posting", id)

		_, err = tg.Send(tb.ChatID(cfg.ChannelID), &tb.Video{File: tb.File{FileURL: link}})
		if err != nil {
			log.Println("Send video", err, link)
		}

		time.Sleep(time.Second * 3)
	}
}
