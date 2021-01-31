package query

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goserver/config"
)

type CHATROOM struct {
	Id       int
	UserId   string
	RoomName string
}

var confDBCR *config.ConfigCR
var ConStrCR string

func init() {
	_confDBCR, err := config.ReadConfDBCR()
	if err != nil {
		fmt.Println(err.Error())
	}
	confDBCR = _confDBCR
	_conStrCR := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", confDBCR.User, confDBCR.Pass, confDBCR.Host, confDBCR.Port, confDBCR.DbName, confDBCR.Charset)
	ConStrCR = _conStrCR
}

// 新規チャットルーム登録関数
func InsertChatroom(userId string, roomName string, db *sql.DB) bool {

	// プリペアードステートメント
	stmt, err := db.Prepare("INSERT INTO chatroom(USER_ID, ROOM_NAME) VALUES(?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	// クエリ実行
	insertedOrNot, err := stmt.Exec(userId, roomName)
	if err != nil {
		return false
	} else {
		_ = insertedOrNot
		return true
	}
}

// 単一行データ取得関数
/*func SelectChatroomById(id int64, db *sql.DB) (chatroominfo CHATROOM) {

	// 構造体CHATROOM型の変数chatroomを宣言
	chatroom := CHATROOM{}

	// プリペアードステートメント
	err := db.QueryRow("SELECT ID, USER_ID, ROOM_NAME FROM chatroom WHERE ID = ?", id).Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName)
	if err != nil {
		return
	}
	return chatroom
}*/

func SelectChatroomByUserId(userSessionVal string, db *sql.DB) map[int]CHATROOM {

	// 構造体CHATROOM型の変数chatroomを宣言

	chatrooms := make(map[int]CHATROOM)

	// プリペアードステートメント
	rows, err := db.Query("SELECT * FROM chatroom WHERE USER_ID = ?", userSessionVal)
	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		chatroom := CHATROOM{}
		err := rows.Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName)
		if err != nil {
			fmt.Println(err.Error())
		}
		chatrooms[chatroom.Id] = chatroom
		fmt.Println(chatrooms)
	}
	return chatrooms
}
