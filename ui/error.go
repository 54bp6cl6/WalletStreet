package ui

import "github.com/line/line-bot-sdk-go/linebot"

func ErrorMessage(err error) linebot.SendingMessage {
	return linebot.NewTextMessage(err.Error())
}
