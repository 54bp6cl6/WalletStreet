package function

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	// TODO: Go to next middleware
	return
}
