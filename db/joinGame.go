package db

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// database
const (
	Games = "Games"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 資料庫 Games 裡的文件模版
type game struct {
	CreateTime time.Time
	Players    []string
}

// 查詢使用者是否已經加入遊戲
func IsUserInGame(userId string) (inGame bool, err error) {
	query := client.Collection(Games).Where("Players", "array-contains", userId)
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

// 於資料庫中創建新遊戲，創建人(userId)將會自動加入
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

// 查詢遊戲編號是否存在對應的房間
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

// 將使用者加入遊戲
func JoinGame(gameId string, userId string) (err error) {
	doc := client.Collection(Games).Doc(gameId)
	// 使用 Transaction
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// 取得房間資料
		var docsnap *firestore.DocumentSnapshot
		docsnap, err := tx.Get(doc)
		if status.Code(err) == codes.NotFound {
			err := fmt.Errorf("game %v not found", gameId)
			return err
		} else if err != nil {
			return err
		}
		g := game{}
		if err = docsnap.DataTo(&g); err != nil {
			return err
		}

		// 加入使用者並回寫
		g.Players = append(g.Players, userId)
		tx.Set(doc, g)
		return err
	})
	return
}

// 生成不重複的遊戲編號
func generateGameId() (gameId string, err error) {
	for i := 0; i < 10; i++ { // 最多嘗試生成10次
		gameId = fmt.Sprintf("%04d", random.Intn(10000))
		var exist bool
		if exist, err = IsGameExist(gameId); err != nil || !exist {
			return
		}
	}
	err = errors.New("generate game id failed to many times")
	return
}
