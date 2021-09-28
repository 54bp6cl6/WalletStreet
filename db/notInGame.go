package db

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// database
	Games = "Games"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type game struct {
	CreateTime time.Time
	Players    []string
}

// 查詢使用者是否已經加入遊戲
func IsUserInGame(userId string) (inGame bool, err error) {
	query := client.Collection(Games).Where("players", "array-contains", userId)
	iter := query.Documents(ctx)
	defer iter.Stop()

	if _, err = iter.Next(); err == iterator.Done {
		inGame = false
		err = nil
		return
	} else if err != nil {
		return
	}

	inGame = true
	return
}

func CreateGame(userId string) (gameId string, err error) {
	if gameId, err = generateGameId(); err != nil {
		return
	}
	_, err = client.Collection(Games).Doc(gameId).Set(ctx,
		game{
			CreateTime: time.Now(),
			Players:    []string{userId},
		},
	)
	return
}

func IsGameExist(gameId string) (exist bool, err error) {
	_, err = client.Collection(Games).Doc(gameId).Get(ctx)

	if status.Code(err) == codes.NotFound {
		err = nil
		exist = false
		return
	} else if err != nil {
		return
	} else {
		exist = true
		return
	}
}

func generateGameId() (gameId string, err error) {
	for i := 0; i < 10; i++ {
		gameId = fmt.Sprintf("%04d", rand.Intn(10000))
		var exist bool
		if exist, err = IsGameExist(gameId); err != nil || !exist {
			return
		}
	}
	err = errors.New("generate game id failed to many times")
	return
}
