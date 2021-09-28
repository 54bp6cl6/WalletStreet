package function

import (
	"github.com/54bp6cl6/WalletStreet/ui"
	"github.com/line/line-bot-sdk-go/linebot"
)

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
