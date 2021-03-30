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
		if ok := session.Manager.SessionIdCheck(w, r); !ok {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
			return
		}
		//適当にルームIDを変えると、他の人のルームが覗けるので、メンバのルームしかアクセスできないよう処理
		//userCookie, _ := r.Cookie(session.Manager.CookieName)
		//userSid, _ := url.QueryUnescape(userCookie.Value)
		//userSessionVar := session.Manager.SessionStore[userSid].SessionValue["userId"]

		roomUrl := r.URL.Path
		_roomId := strings.TrimPrefix(roomUrl, "/mypage/chatroom")
		roomId, _ := strconv.Atoi(_roomId)

		dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbChtrm.Close()

		selectedChatroom := query.SelectChatroomById(roomId, dbChtrm)
		//userId := selectedChatroom.UserId
		//member := selectedChatroom.Member

		//if userId != userSessionVar && member != userSessionVar {
		//	fmt.Fprintf(w, "ルームにアクセスする権限がありません")
		//	return
		//}

		Chats := query.SelectAllChatsById(selectedChatroom.Id, dbChtrm)
		fmt.Println(Chats)

		t := template.Must(template.ParseFiles("./templates/mypage/chatroom.html"))
		t.ExecuteTemplate(w, "chatroom.html", Chats)

	/*新しい書き込みがあった時の処理。
	書き込み主はセッション変数から判別。誰宛かは部屋の作成者と書き込み主を比較し判断
	*/
	case "POST":
		if ok := session.Manager.SessionIdCheck(w, r); !ok {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
			return
		}

		err := r.ParseForm()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(r.Form)
		if r.FormValue("chat") != "" {
			newChat := new(query.Chat)
			newChat.Chat = r.FormValue("chat")
			if newChat.Chat == "" {
				fmt.Fprintf(w, "何も入力されていません")
				return
			}

			roomUrl := r.URL.Path
			_roomId := strings.TrimPrefix(roomUrl, "/mypage/chatroom")
			roomId, _ := strconv.Atoi(_roomId)

			dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbChtrm.Close()

			currentChatroom := query.SelectChatroomById(roomId, dbChtrm)

			newChat.Chatroom.Id = currentChatroom.Id
			newChat.Chatroom.RoomName = currentChatroom.RoomName

			userCookie, _ := r.Cookie(session.Manager.CookieName)
			userSid, _ := url.QueryUnescape(userCookie.Value)
			userSessionVar := session.Manager.SessionStore[userSid].SessionValue["userId"]

			if userSessionVar == currentChatroom.UserId {
				//投稿主と部屋作成者が同じ場合
				newChat.Chatroom.UserId = userSessionVar
				newChat.Chatroom.Member = currentChatroom.Member
			} else {
				//投稿主と部屋作成者が違う場合
				newChat.Chatroom.UserId = currentChatroom.Member
				newChat.Chatroom.Member = currentChatroom.UserId
			}
			newChat.PostDt = time.Now().UTC().Round(time.Second)

			posted := query.InsertChat(newChat.Chatroom.Id, newChat.Chatroom.UserId, newChat.Chatroom.RoomName, newChat.Chatroom.Member, newChat.Chat, newChat.PostDt, dbChtrm)
			if posted == true {
				return
			} else {
				fmt.Fprintf(w, "投稿できませんでした")
				return
			}
		}
		if r.FormValue("delete-room") != "" {
			roomUrl := r.URL.Path
			_roomId := strings.TrimPrefix(roomUrl, "/mypage/chatroom")
			roomId, _ := strconv.Atoi(_roomId)

			dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbChtrm.Close()

			query.DeleteChatroomById(roomId, dbChtrm)
		}
	}
}
