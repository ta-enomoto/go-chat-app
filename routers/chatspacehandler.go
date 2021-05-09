//チャットスペースを読み込むためのハンドラ
package routers

import (
	"database/sql"
	"fmt"
	"goserver/query"
	"goserver/sessions"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func ChatspaceHandler(w http.ResponseWriter, r *http.Request) {
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

	t := template.Must(template.ParseFiles("./templates/chatspace.html"))
	t.ExecuteTemplate(w, "chatspace.html", Chats)

}
