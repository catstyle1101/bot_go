package telegram

import (
	"bot_go/clients/telegram"
	"bot_go/lib/e"
	"bot_go/storage"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("command: text=%q, username=%q", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendStart(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	sendMsg := chatIDMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err = p.storage.Save(page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: send random", err) }()

	sendMsg := chatIDMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}
func (p *Processor) sendStart(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func chatIDMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
