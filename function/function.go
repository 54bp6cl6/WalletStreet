package function

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/54bp6cl6/WalletStreet/db"
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
	db.Connect()
	defer db.Close()
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
