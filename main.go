package main

import (
	tgClient "bot_go/clients/telegram"
	"bot_go/consumer/event-consumer"
	"bot_go/events/telegram"
	"bot_go/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Println("[INFO] service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("[FATAL] service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("token", "", "Telegram bot token")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is required")
	}
	return *token
}
