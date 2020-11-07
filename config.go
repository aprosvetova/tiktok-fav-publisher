package main

type config struct {
	DBAddr string `env:"DB_ADDR" envDefault:"localhost:6379"`

	TelegramToken string `env:"TG_TOKEN,required"`
	ChannelID     int64  `env:"CHANNEL_ID,required"`

	LikesServiceURL    string `env:"LIKES_SERVICE,required"`
	DownloadServiceURL string `env:"DOWNLOAD_SERVICE,required"`

	TikTokUsername string `env:"TIKTOK_USERNAME,required"`
}
