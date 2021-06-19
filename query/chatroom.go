//チャット情報を扱うクエリパッケージ
package query

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goserver/config"
	"time"
)

/*チャットルームごとにtableを動的に生成できないため、
ルーム情報とチャット情報を分け、別々のdbに保存する。
①マイページではルーム情報のみ扱う
②アクセスがあったら、ルーム情報を元にチャットを取得する*/

/*チャットルーム構造体(テーブル「ROOM_STRUCTS_OF_CHAT」に保存)
マイページでルーム情報だけ取得するために使用
チャット情報からルーム情報を抽出しようとすると処理が重くなる
(チャットルームはIDがPKになっているため重複がない)*/
type Chatroom struct {
	Id       int    `json:"id"`
	UserId   string `json:"userId"`
	RoomName string `json:"roomName"`
	Member   string `json:"member"`
}

/*チャット構造体(テーブル「ALL_STRUCTS_OF_CHAT」に保存)
各チャットルーム内で投稿された書き込みの情報
ルームの情報はチャットルーム構造体と共通
*/
type Chat struct {
	Chatroom Chatroom
	Chat     string
	PostDt   time.Time
}

var ConStrChtrm string

func init() {
	confDbChtrm, err := config.ReadConfDbChtrm()
	if err != nil {
		fmt.Println(err.Error())
	}
	ConStrChtrm = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=%s", confDbChtrm.User, confDbChtrm.Pass, confDbChtrm.Host, confDbChtrm.Port, confDbChtrm.DbName, confDbChtrm.Charset)
}

// 新規チャットルーム登録関数
func InsertChatroom(userSessionVal string, roomName string, memberName string, db *sql.DB) bool {

	stmt, err := db.Prepare("INSERT INTO ROOM_STRUCTS_OF_CHAT(USER_ID, ROOM_NAME, MEMBER) VALUES(?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(userSessionVal, roomName, memberName)
	if err != nil {
		return false
	} else {
		return true
	}
}

//特定のユーザーが作成したチャットルームをすべて取得する
func SelectAllChatroomsByUserId(userSessionVal string, db *sql.DB) (chatrooms []Chatroom) {

	rows, err := db.Query("SELECT * FROM ROOM_STRUCTS_OF_CHAT WHERE USER_ID = ?", userSessionVal)
	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		chatroom := Chatroom{}
		err := rows.Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
		if err != nil {
			fmt.Println(err.Error())
		}
		chatrooms = append(chatrooms, chatroom)
	}
	return
}

//特定のユーザーがメンバーとして参加しているチャットルームをすべて取得する
func SelectAllChatroomsByMember(userSessionVal string, db *sql.DB) (chatrooms []Chatroom) {

	rows, err := db.Query("SELECT * FROM ROOM_STRUCTS_OF_CHAT WHERE Member = ?", userSessionVal)
	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		chatroom := Chatroom{}
		err := rows.Scan(&chatroom.Id, &chatroom.Member, &chatroom.RoomName, &chatroom.UserId)
		if err != nil {
			fmt.Println(err.Error())
		}
		chatrooms = append(chatrooms, chatroom)
	}
	return
}

//idでチャットルームを選択
func SelectChatroomById(id int, db *sql.DB) (chatroom Chatroom) {

	err := db.QueryRow("SELECT * FROM ROOM_STRUCTS_OF_CHAT WHERE ID = ?", id).Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

func SelectChatroomByUser(userId string, db *sql.DB) (chatroom Chatroom) {

	err := db.QueryRow("SELECT ID, USER_ID, ROOM_NAME, MEMBER FROM ROOM_STRUCTS_OF_CHAT WHERE USER_ID = ?").Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

func SelectChatroomByUserAndRoomNameAndMember(userId string, roomName string, member string, db *sql.DB) (chatroom Chatroom) {

	err := db.QueryRow("SELECT * FROM ROOM_STRUCTS_OF_CHAT WHERE USER_ID = ? AND ROOM_NAME = ? AND MEMBER = ?", userId, roomName, member).Scan(&chatroom.Id, &chatroom.UserId, &chatroom.RoomName, &chatroom.Member)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

//特定のチャットルームのチャットをすべて取得する
func SelectAllChatsById(id int, db *sql.DB) (chats []Chat) {

	rows, err := db.Query("SELECT * FROM ALL_STRUCTS_OF_CHAT WHERE ID = ?", id)
	if err != nil {
		return chats
	}

	for rows.Next() {
		chat := Chat{}
		err := rows.Scan(&chat.Chatroom.Id, &chat.Chatroom.UserId, &chat.Chatroom.RoomName, &chat.Chatroom.Member, &chat.Chat, &chat.PostDt)
		if err != nil {
			fmt.Println(err.Error())
		}
		chats = append(chats, chat)
	}
	return
}

// 新規チャット投稿関数
func InsertChat(id int, userId string, roomName string, member string, chat string, postDt time.Time, db *sql.DB) bool {

	stmt, err := db.Prepare("INSERT INTO ALL_STRUCTS_OF_CHAT(ID, USER_ID, ROOM_NAME, MEMBER, Chat, POST_DT) VALUES(?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, userId, roomName, member, chat, postDt)
	if err != nil {
		return false
	} else {
		return true
	}
}

//チャットルーム削除関数
func DeleteChatroomById(id int, db *sql.DB) bool {

	stmtDeleteChatroom, err := db.Prepare("DELETE FROM ROOM_STRUCTS_OF_CHAT WHERE ID = ?")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmtDeleteChatroom.Close()

	stmtDeleteChat, err := db.Prepare("DELETE FROM ALL_STRUCTS_OF_CHAT WHERE ID = ?")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmtDeleteChat.Close()

	_, err = stmtDeleteChatroom.Exec(id)
	if err != nil {
		return false
	}

	_, err = stmtDeleteChat.Exec(id)
	if err != nil {
		return false
	} else {
		return true
	}

}
