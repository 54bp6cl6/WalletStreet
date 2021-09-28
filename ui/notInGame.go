package ui

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

func CreateGameMessage(gameId string) linebot.SendingMessage {
	return linebot.NewTextMessage(fmt.Sprintf("創建遊戲成功!!遊戲編號: %v", gameId))
}
