/*個別のチャットルームにアクセスがあったときのハンドラ*/
package routers

import (
	"database/sql"
	"fmt"
	"goserver/query"
	"goserver/sessions"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ChatroomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	/*アクセスあった際、ルームIDが一致するすべての書き込みをスライスで取得し、テンプレに渡す*/
	case "GET":
		if ok := session.Manager.SessionIdCheck(w, r); ok {
			roomUrl := r.URL.Path
			_roomId := strings.TrimPrefix(roomUrl, "/mypage/chatroom")
			roomId, _ := strconv.Atoi(_roomId)

			dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbChtrm.Close()

			selectedChatroom := query.SelectChatroomById(roomId, dbChtrm)
			Chat := query.SelectAllChatById(selectedChatroom.Id, dbChtrm)

			t := template.Must(template.ParseFiles("./templates/mypage/chatroom.html"))
			t.ExecuteTemplate(w, "chatroom.html", Chat)
		} else {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
		}

	/*新しい書き込みがあった時の処理。
	書き込み主はセッション変数から判別。誰宛かは部屋の作成者と書き込み主を比較し判断
	*/
	case "POST":
		//sessionチェックすべき
		newChat := new(query.Chat)
		newChat.Chat = r.FormValue("chat")
		if newChat.Chat == "" {
			fmt.Fprintf(w, "何も入力されていません")
		} else {
			roomUrl := r.URL.Path
			_roomId := strings.TrimPrefix(roomUrl, "/mypage/chatroom")
			roomId, _ := strconv.Atoi(_roomId)

			userCookie, _ := r.Cookie(session.Manager.CookieName)
			userSid, _ := url.QueryUnescape(userCookie.Value)
			userId := session.Manager.SessionStore[userSid].SessionValue["userId"]

			dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbChtrm.Close()

			selectedChatroom := query.SelectChatroomById(roomId, dbChtrm)

			//newChat.Chatroom.UserId   = userId
			newChat.Chatroom.Id = selectedChatroom.Id
			newChat.Chatroom.RoomName = selectedChatroom.RoomName

			if userId == selectedChatroom.UserId {
				//投稿主と部屋作成者が同じ場合
				newChat.Chatroom.UserId = userId
				newChat.Chatroom.Member = selectedChatroom.Member
			} else {
				//投稿主と部屋作成者が違う場合
				newChat.Chatroom.UserId = selectedChatroom.Member
				newChat.Chatroom.Member = selectedChatroom.UserId
			}

			newChat.PostDt = time.Now().UTC().Round(time.Second)

			posted := query.InsertChat(newChat.Chatroom.Id,
				newChat.Chatroom.UserId, newChat.Chatroom.RoomName, newChat.Chatroom.Member, newChat.Chat, newChat.PostDt, dbChtrm)
			if posted == true {
				fmt.Fprintf(w, "投稿されました")
			} else {
				fmt.Fprintf(w, "投稿できませんでした")
			}
		}
	}
}
