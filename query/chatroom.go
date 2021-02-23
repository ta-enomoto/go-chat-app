package query

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goserver/config"
	"time"
)

type CHATROOM struct {
	Id       int
	UserId   string
	RoomName string
	Member   string
}

type CHAT struct {
	Id       int
	UserId   string
	RoomName string
	Member   string
	Chat     string
	PostDt   time.Time
}

var confDBCR *config.ConfigCR
var ConStrCR string

func init() {
	_confDBCR, err := config.ReadConfDBCR()
	if err != nil {
		fmt.Println(err.Error())
	}
	confDBCR = _confDBCR
	_conStrCR := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=%s", confDBCR.User, confDBCR.Pass, confDBCR.Host, confDBCR.Port, confDBCR.DbName, confDBCR.Charset)
	ConStrCR = _conStrCR
}

// 新規チャットルーム登録関数
func InsertChatroom(userSessionVal string, roomName string, memberName string, db *sql.DB) bool {

	newChatroom := CHATROOM{UserId: userSessionVal, RoomName: roomName, Member: memberName}
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

//特定のユーザーが作成したチャットルームをすべて取得する
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
		err := rows.Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
		if err != nil {
			fmt.Println(err.Error())
		}
		chatrooms = append(chatrooms, chatroom)
		fmt.Println(chatrooms)
	}
	return chatrooms
}

//特定のユーザーがメンバーとして参加しているチャットルームをすべて取得する
func SelectAllChatroomsByMember(userSessionVal string, db *sql.DB) []CHATROOM {
	var chatrooms []CHATROOM

	// プリペアードステートメント
	rows, err := db.Query("SELECT * FROM chatrooms WHERE Member = ?", userSessionVal)
	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		chatroom := CHATROOM{}
		err := rows.Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
		if err != nil {
			fmt.Println(err.Error())
		}
		chatrooms = append(chatrooms, chatroom)
		fmt.Println(chatrooms)
	}
	return chatrooms
}

// 特定のルームを取得する
/*func SelectChatroomByUser(userId string, db *sql.DB) CHATROOM {

	// 構造体CHATROOM型の変数chatroomを宣言
	chatroom := CHATROOM{}

	// プリペアードステートメント
	err := db.QueryRow("SELECT ID, USER_ID, ROOM_NAME, MEMBER FROM chatroom WHERE USER_ID = ?", userId).Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
	if err != nil {
		fmt.Println(err.Error())
	}
	return chatroom
}*/

func SelectChatroomById(id int, db *sql.DB) CHATROOM {

	// 構造体CHATROOM型の変数chatroomを宣言
	chatroom := CHATROOM{}

	// プリペアードステートメント
	err := db.QueryRow("SELECT * FROM chatrooms WHERE ID = ?", id).Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
	if err != nil {
		fmt.Println(err.Error())
	}
	return chatroom
}

func SelectChatroomByUser(userId string, db *sql.DB) CHATROOM {

	// 構造体CHATROOM型の変数chatroomを宣言
	chatroom := CHATROOM{}

	// プリペアードステートメント
	err := db.QueryRow("SELECT ID, USER_ID, ROOM_NAME, MEMBER FROM chatrooms WHERE USER_ID = ?").Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
	if err != nil {
		fmt.Println(err.Error())
	}
	return chatroom
}

//チャットルームの重複をチェックする
func contains(s []CHATROOM, e CHATROOM) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

//特定のチャットルームのチャットをすべて取得する
func SelectAllChatById(id int, db *sql.DB) []CHAT {

	var chats []CHAT

	// プリペアードステートメント
	rows, err := db.Query("SELECT * FROM chat WHERE ID = ?", id)
	if err != nil {
		return chats
	}

	for rows.Next() {
		chat := CHAT{}
		err := rows.Scan(&chat.Id, &chat.UserId, &chat.RoomName, &chat.Member, &chat.Chat, &chat.PostDt)
		if err != nil {
			fmt.Println(err.Error())
		}
		chats = append(chats, chat)
		fmt.Println(chats)
	}
	return chats
}

// 新規チャット投稿関数
func InsertChat(id int, userId string, roomName string, member string, chat string, postDt time.Time, db *sql.DB) bool {

	//newChat := CHAT{UserId: userId, RoomName: roomName, Member: member, Chat: chat, PostDt: postDt}
	stmt, err := db.Prepare("INSERT INTO chat(ID, USER_ID, ROOM_NAME, MEMBER, CHAT, POST_DT) VALUES(?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()
	fmt.Println(id, userId, roomName, member, chat, postDt)
	insertedOrNot1, err := stmt.Exec(id, userId, roomName, member, chat, postDt)
	if err != nil {
		return false
	} else {
		_ = insertedOrNot1
		return true
	}
}
