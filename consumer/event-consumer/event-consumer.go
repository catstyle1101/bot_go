package event_consumer

import (
	"bot_go/events"
	"log"
	"sync"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

	}
}

/*
TODO:
1. Потеря событий: ретраи, возвращение в хранилище, фоллбек, подтверждение для fetcher
2. обработка всей пачки: останавливаться после первой ошибки, счетчик ошибок
3. Параллельная обработка (sync.WaitGroup{}) ✅
*/

func (c *Consumer) handleEvents(events []events.Event) error {
	var wg sync.WaitGroup

	for _, event := range events {
		log.Printf("[INFO] consumer received event: %s", event)
		wg.Add(1)

		go func() {
			defer wg.Done()
			if err := c.processor.Process(event); err != nil {
				log.Printf("[ERR] consumer: %s", err.Error())
			}
		}()
	}
	wg.Wait()
	return nil
}
