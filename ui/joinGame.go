package ui

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

func CreateGameSuccess(gameId string) linebot.SendingMessage {
	return linebot.NewTextMessage(fmt.Sprintf("創建遊戲成功!!遊戲編號: %v", gameId))
}

func JoinGameSuccess() linebot.SendingMessage {
	return linebot.NewTextMessage("加入遊戲成功~")
}

func JoinGameFail() linebot.SendingMessage {
	return linebot.NewTextMessage("加入遊戲失敗，請再次確認想加入的遊戲編號~")
}
