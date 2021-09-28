package function

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/54bp6cl6/WalletStreet/db"
	"github.com/54bp6cl6/WalletStreet/postback"
	"github.com/54bp6cl6/WalletStreet/ui"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	bot *linebot.Client
)

func init() {
	var err error
	bot, err = linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		log.Fatal(err)
		return
	}
}

// Cloud Function 進入點
func Webhook(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		fmt.Fprint(w, err)
		return
	}

	for _, event := range events {
		err = HandleEvent(event)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprint(w, err)
			return
		}
	}

	fmt.Fprint(w, "OK")
}

// Bot 邏輯的起始點
func HandleEvent(event *linebot.Event) (err error) {
	err = HandleFollowEvent(event)
	return
}

// 處理加入好友與封鎖事件的中間層
func HandleFollowEvent(event *linebot.Event) (err error) {
	switch event.Type {
	case linebot.EventTypeFollow:
		_, err = bot.ReplyMessage(event.ReplyToken, ui.FollowMessage()).Do()
		return
	case linebot.EventTypeUnfollow:
		// TODO: 刪除所有使用者資料
		return
	}

	// 進入下一層
	err = HandleJoinGameEvent(event)
	return
}

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
			// 生成遊戲編號
			var gameId string
			if gameId, err = db.CreateGame(event.Source.UserID); err != nil {
				_, err = bot.ReplyMessage(event.ReplyToken, ui.ErrorMessage(err)).Do()
				return
			}
			_, err = bot.ReplyMessage(event.ReplyToken, ui.CreateGameMessage(gameId)).Do()
		}

		return // 無效的 Postback 不予回應

	// 處理加入遊戲的文字訊息
	case linebot.EventTypeMessage:
		return
	}

	_, err = bot.ReplyMessage(event.ReplyToken, ui.FollowMessage()).Do()
	return
}
