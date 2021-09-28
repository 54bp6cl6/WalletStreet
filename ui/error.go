package ui

import "github.com/line/line-bot-sdk-go/linebot"

// 將錯誤訊息傳送至聊天室
func ErrorMessage(err error) linebot.SendingMessage {
	return linebot.NewTextMessage(err.Error())
}
