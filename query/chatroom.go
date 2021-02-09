package query

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goserver/config"
)

type CHATROOM struct {
	//Id       int
	UserId   string
	RoomName string
	Member   string
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

//ユーザー別新規テーブル作成関数
//func CreateUsersChatrooms(userId string, db *sql.DB) bool {
//	stmt, err := db.Prepare("create table chatroom_?(id INT AUTO_INCREMENT NOT NULL PRIMARY KEY, USER_ID VARCHAR(10), ROOM_NAME VARCHAR(10));")
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	defer stmt.Close()
//
//	insertedOrNot, err := stmt.Exec(userId)
//	if err != nil {
//		return false
//	} else {
//		_ = insertedOrNot
//		return true
//	}
//}

// 新規チャットルーム登録関数
func InsertChatroom(userSessionVal string, roomName string, memberName string, db *sql.DB) bool {

	newChatroom := CHATROOM{userSessionVal, roomName, memberName}
	chatrooms := SelectAllChatroomsByUserId(userSessionVal, db)
	roomExist := contains(chatrooms, newChatroom)
	if roomExist == true {
		return false
	} else {
		//チャットルームを作成したユーザーからの登録
		stmt, err := db.Prepare("INSERT INTO chatrooms(USER_ID, ROOM_NAME, MEMBER) VALUES(?,?,?)")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer stmt.Close()
		insertedOrNot1, err := stmt.Exec(userSessionVal, roomName, memberName)
		if err != nil {
			return false
		} else {
			_ = insertedOrNot1
			return true
		}
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

func SelectAllChatroomsByUserId(userSessionVal string, db *sql.DB) []CHATROOM {

	// 構造体CHATROOM型の変数chatroomを宣言

	var chatrooms []CHATROOM

	// プリペアードステートメント
	rows, err := db.Query("SELECT * FROM chatrooms WHERE USER_ID = ?", userSessionVal)
	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		chatroom := CHATROOM{}
		err := rows.Scan(&chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
		if err != nil {
			fmt.Println(err.Error())
		}
		chatrooms = append(chatrooms, chatroom)
		fmt.Println(chatrooms)
	}
	return chatrooms
}

func SelectAllChatroomsByMember(userSessionVal string, db *sql.DB) []CHATROOM {
	var chatrooms []CHATROOM

	// プリペアードステートメント
	rows, err := db.Query("SELECT * FROM chatrooms WHERE Member = ?", userSessionVal)
	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		chatroom := CHATROOM{}
		err := rows.Scan(&chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
		if err != nil {
			fmt.Println(err.Error())
		}
		chatrooms = append(chatrooms, chatroom)
		fmt.Println(chatrooms)
	}
	return chatrooms
}

func contains(s []CHATROOM, e CHATROOM) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
