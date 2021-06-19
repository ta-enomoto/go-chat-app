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
	//"time"
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
		userCookie, _ := r.Cookie(session.Manager.CookieName)
		userSid, _ := url.QueryUnescape(userCookie.Value)
		userSessionVar := session.Manager.SessionStore[userSid].SessionValue["userId"]

		//本番環境で
		roomUrl := r.URL.Path
		_roomId := strings.TrimPrefix(roomUrl, "/mypage/chatroom")
		roomId, _ := strconv.Atoi(_roomId)

		dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbChtrm.Close()

		selectedChatroom := query.SelectChatroomById(roomId, dbChtrm)
		userId := selectedChatroom.UserId
		member := selectedChatroom.Member

		if userId != userSessionVar && member != userSessionVar {
			fmt.Fprintf(w, "ルームにアクセスする権限がありません")
			return
		}

		Chats := query.SelectAllChatsById(selectedChatroom.Id, dbChtrm)

		t := template.Must(template.ParseFiles("./templates/mypage/chatroom.html"))
		t.ExecuteTemplate(w, "chatroom.html", Chats)

	case "POST":
		if ok := session.Manager.SessionIdCheck(w, r); !ok {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
			return
		}

		if r.FormValue("delete-room") == "このルームを削除する" {
			roomUrl := r.URL.Path
			_roomId := strings.TrimPrefix(roomUrl, "/mypage/chatroom")
			roomId, _ := strconv.Atoi(_roomId)

			dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbChtrm.Close()

			query.DeleteChatroomById(roomId, dbChtrm)

			t := template.Must(template.ParseFiles("./templates/mypage/chatroomdeleted.html"))
			t.ExecuteTemplate(w, "chatroomdeleted.html", nil)
		}
	}
}
