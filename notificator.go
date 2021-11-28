package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Notificator interface {
	Notify(title string, message string)
}

type IFTTTNotificator struct {
	key string
}

type ArrayNotificator struct {
	notificators []Notificator
}

type ConsoleNotificator struct {
}

func NewIFTTTNotificator(key string) *IFTTTNotificator {
	return &IFTTTNotificator{key: key}
}

func (n IFTTTNotificator) Notify(title string, message string) {
	type Email struct {
		Title   string `json:"value1"`
		Message string `json:"value2"`
	}

	email := Email{title, message}

	json, err := json.Marshal(email)
	if err != nil {
		log.Fatal(err)
	}

	_, err = http.Post("https://maker.ifttt.com/trigger/email/with/key/"+n.key, "application/json", bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("IFTTT Notification OK")
}

func (a ArrayNotificator) Notify(title string, message string) {
	for _, notificator := range a.notificators {
		notificator.Notify(title, message)
	}
}

func NewArrayNotificator(notificators ...Notificator) *ArrayNotificator {
	return &ArrayNotificator{notificators: notificators}
}

func (c ConsoleNotificator) Notify(title string, message string) {
	fmt.Println("Title:", title)
	fmt.Println("Message:\n", message)
}

func NewConsoleNotificator() *ConsoleNotificator {
	return &ConsoleNotificator{}
}
