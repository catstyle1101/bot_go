package main

import (
	"bot_go/clients/telegram"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(mustToken(), tgBotHost)

	// fetcher = fetcher.New()

	// processor = processor.New()

	//consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String("token", "", "Telegram bot token")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is required")
	}
	return *token
}
