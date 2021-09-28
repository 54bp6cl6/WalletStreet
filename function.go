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

func HandleEvent(event *linebot.Event) (err error) {
	err = HandleFollowEvent(event)
	return
}

func HandleFollowEvent(event *linebot.Event) (err error) {
	switch event.Type {
	case linebot.EventTypeFollow:
		_, err = bot.ReplyMessage(event.ReplyToken, ui.FollowMessage()).Do()
		return
	case linebot.EventTypeUnfollow:
		return
	}

	err = HandleNotInGameEvent(event)
	return
}

func HandleNotInGameEvent(event *linebot.Event) (err error) {
	var inGame bool
	if inGame, err = db.IsUserInGame(event.Source.UserID); err != nil {
		_, err = bot.ReplyMessage(event.ReplyToken, ui.ErrorMessage(err)).Do()
		return
	}

	if inGame {
		// TODO: Go to next middleware
		fmt.Print("跳過NotInGame Middleware")
		return
	}

	fmt.Printf("進入NotInGame Middleware %v", event.Type)
	switch event.Type {
	case linebot.EventTypePostback:
		fmt.Printf("Postback: %v", event.Postback.Data)
		var data map[string]interface{}
		if data, err = postback.ToMap(event.Postback.Data); err != nil {
			_, err = bot.ReplyMessage(event.ReplyToken, ui.ErrorMessage(err)).Do()
			return
		}

		if data[postback.Action] == postback.CreateGame {
			var gameId string
			if gameId, err = db.CreateGame(event.Source.UserID); err != nil {
				_, err = bot.ReplyMessage(event.ReplyToken, ui.ErrorMessage(err)).Do()
				return
			}
			_, err = bot.ReplyMessage(event.ReplyToken, ui.CreateGameMessage(gameId)).Do()
			return
		}
	case linebot.EventTypeMessage:
		return
	}
	return
}
