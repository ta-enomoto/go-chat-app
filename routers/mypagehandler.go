//マイページへアクセスがあった時のハンドラ
package routers

import (
	"database/sql"
	"fmt"
	"goserver/query"
	"goserver/sessions"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

func MypageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	/*アクセスがあった時の処理。自身で作成したルームと他人が作成したルームを別々に取得して、
	ルーム一覧のスライスをつくってテンプレに渡す*/
	case "GET":
		if ok := session.Manager.SessionIdCheck(w, r); !ok {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
			return
		}

		t := template.Must(template.ParseFiles("./templates/mypage.html"))

		dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbChtrm.Close()

		userCookie, _ := r.Cookie(session.Manager.CookieName)
		userSid, _ := url.QueryUnescape(userCookie.Value)
		userSessionVar := session.Manager.SessionStore[userSid].SessionValue["userId"]

		chatroomsFromUserId := query.SelectAllChatroomsByUserId(userSessionVar, dbChtrm)
		chatroomsFromMember := query.SelectAllChatroomsByMember(userSessionVar, dbChtrm)
		var Links = append(chatroomsFromUserId, chatroomsFromMember...)

		t.ExecuteTemplate(w, "mypage.html", Links)

		/*新しいルーム作成のポストがあった時の処理。ルーム名と相手メンバーを指定する
		同名のルーム名は、相手メンバー異なる場合のみ有効。*/
	case "POST":
		if ok := session.Manager.SessionIdCheck(w, r); !ok {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
			return
		}

		newchatroom := new(query.Chatroom)
		newchatroom.RoomName = r.FormValue("roomName")
		newchatroom.Member = r.FormValue("memberName")

		//メンバーまたはルーム名が入力されていない
		if newchatroom.RoomName == "" || newchatroom.Member == "" {
			return
		}

		userCookie, _ := r.Cookie(session.Manager.CookieName)
		userSid, _ := url.QueryUnescape(userCookie.Value)
		userSessionVar := session.Manager.SessionStore[userSid].SessionValue["userId"]

		//自分自身をメンバーに加えることはできない
		if newchatroom.Member == userSessionVar {
			return
		}

		dbUsr, err := sql.Open("mysql", query.ConStrUsr)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbUsr.Close()

		users := query.SelectAllUser(dbUsr)
		userIdExist := query.ContainsUserName(users, newchatroom.Member)

		//相手ユーザーが存在しない
		if !userIdExist {
			return
		}

		dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbChtrm.Close()

		newChtrmInsertedToDb := query.InsertChatroom(userSessionVar, newchatroom.RoomName, newchatroom.Member, dbChtrm)
		if newChtrmInsertedToDb {

			dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbChtrm.Close()

			newChat := new(query.Chat)
			newChat.Chat = "NEW ROOM CREATED"

			createdChatroom := query.SelectChatroomByUserAndRoomNameAndMember(userSessionVar, newchatroom.RoomName, newchatroom.Member, dbChtrm)

			newChat.Chatroom.Id = createdChatroom.Id
			newChat.Chatroom.RoomName = createdChatroom.RoomName

			if userSessionVar == createdChatroom.UserId {
				//投稿主と部屋作成者が同じ場合
				newChat.Chatroom.UserId = userSessionVar
				newChat.Chatroom.Member = createdChatroom.Member
			} else {
				//投稿主と部屋作成者が違う場合
				newChat.Chatroom.UserId = createdChatroom.Member
				newChat.Chatroom.Member = createdChatroom.UserId
			}
			newChat.PostDt = time.Now().UTC().Round(time.Second)
			posted := query.InsertChat(newChat.Chatroom.Id, newChat.Chatroom.UserId, newChat.Chatroom.RoomName, newChat.Chatroom.Member, newChat.Chat, newChat.PostDt, dbChtrm)
			if posted == true {
				fmt.Println("初期投稿成功")
				return
			} else {
				fmt.Println("初期投稿失敗")
				return
			}
		} else {
			fmt.Println("ルーム作成失敗")
		}
	}
}
