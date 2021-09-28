package postback

import "fmt"

const (
	// fields
	Action = "action"

	// Actions
	CreateGame = "create_game"
)

func CreateGameData() string {
	return fmt.Sprintf(`{"%v": "%v"}`, Action, CreateGame)
}
