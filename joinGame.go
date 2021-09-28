package function

import (
	"regexp"
	"strings"

	"github.com/54bp6cl6/WalletStreet/db"
	"github.com/54bp6cl6/WalletStreet/postback"
	"github.com/54bp6cl6/WalletStreet/ui"
	"github.com/line/line-bot-sdk-go/linebot"
)

// 處理創建與加入遊戲的中間層
func HandleJoinGameEvent(event *linebot.Event) (err error) {
	// 取得使用者是否在遊戲中
	var inGame bool
	if inGame, err = db.IsUserInGame(event.Source.UserID); err != nil {
		_, err = bot.ReplyMessage(event.ReplyToken, ui.ErrorMessage(err)).Do()
		return
	}

	// 在遊戲中: 前往下一層
	if inGame {
		// TODO: Go to next middleware
		return
	}

	switch event.Type {
	// 處理創建遊戲的 Postback
	case linebot.EventTypePostback:
		// 將 Postback Data 轉為字典
		var data map[string]interface{}
		if data, err = postback.ToMap(event.Postback.Data); err != nil {
			_, err = bot.ReplyMessage(event.ReplyToken, ui.ErrorMessage(err)).Do()
			return
		}

		// 創建新遊戲
		if data[postback.Action] == postback.CreateGame {
			var gameId string
			if gameId, err = db.CreateGame(event.Source.UserID); err != nil {
				_, err = bot.ReplyMessage(event.ReplyToken, ui.ErrorMessage(err)).Do()
				return
			}
			_, err = bot.ReplyMessage(event.ReplyToken, ui.CreateGameSuccess(gameId)).Do()
		}

		return // 無效的 Postback 不予回應

	// 處理加入遊戲的文字訊息
	case linebot.EventTypeMessage:
		if message, ok := event.Message.(*linebot.TextMessage); ok {
			// 檢查訊息是否為 4 個數字
			message.Text = strings.Trim(message.Text, " ")
			var match bool
			if match, err = regexp.MatchString("^[0-9]{4}$", message.Text); err != nil {
				_, err = bot.ReplyMessage(event.ReplyToken, ui.ErrorMessage(err)).Do()
				return
			}

			if match {
				// 檢查遊戲是否存在
				var exist bool
				if exist, err = db.IsGameExist(message.Text); err != nil {
					return
				}

				if exist {
					if err = db.JoinGame(message.Text, event.Source.UserID); err != nil {
						return
					}
					_, err = bot.ReplyMessage(event.ReplyToken, ui.JoinGameSuccess()).Do()
				} else {
					_, err = bot.ReplyMessage(event.ReplyToken, ui.JoinGameFail()).Do()
				}
				return
			}
		}
	}

	_, err = bot.ReplyMessage(event.ReplyToken, ui.FollowMessage()).Do()
	return
}
