package ui

import (
	"fmt"

	"github.com/54bp6cl6/WalletStreet/postback"
	"github.com/line/line-bot-sdk-go/linebot"
)

func FollowMessage() linebot.SendingMessage {
	altText := "嗨~歡迎使用 Wallet Street 瓦雷街~"
	buttonLabel := "創建遊戲"
	buttonData := postback.CreateGameData()
	greeting := fmt.Sprintf("%v\n您可以輸入遊戲編號來加入別人的遊戲，或是點擊下方「%v」按鈕自己創建一局。", altText, buttonLabel)

	contents := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: greeting,
					Wrap: true,
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Action: linebot.NewPostbackAction(buttonLabel, buttonData, "", ""),
				},
			},
		},
	}

	return linebot.NewFlexMessage(altText, contents)
}
