package postback

import "fmt"

// fields
const (
	Action = "action"
)

// Actions
const (
	CreateGame = "create_game"
)

// 取得 Create Game Postback button 要夾帶的資料
func CreateGameData() string {
	return fmt.Sprintf(`{"%v": "%v"}`, Action, CreateGame)
}
