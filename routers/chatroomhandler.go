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
	case "GET":
		if ok := session.Manager.SidCheck(w, r); ok {
			roomUrl := r.URL.Path
			//userCookie, _ := r.Cookie(session.Manager.CookieName)
			//userSid, _ := url.QueryUnescape(userCookie.Value)
			//userId := session.Manager.Database[userSid].SessionValue["ID"]
			_roomId := strings.TrimPrefix(roomUrl, "/mypage/chatroom")
			roomId, _ := strconv.Atoi(_roomId)

			dbCR, err := sql.Open("mysql", query.ConStrCR)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbCR.Close()

			selectedChatroom := query.SelectChatroomById(roomId, dbCR)
			fmt.Println(selectedChatroom.Id)
			Chat := query.SelectAllChatById(selectedChatroom.Id, dbCR)
			fmt.Println(Chat)

			t := template.Must(template.ParseFiles("./templates/mypage/chatroom.html"))
			t.ExecuteTemplate(w, "chatroom.html", Chat)
		} else {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
		}
	case "POST":
		newChat := new(query.CHAT)
		newChat.Chat = r.FormValue("chat")
		if newChat.Chat == "" {
			fmt.Fprintf(w, "何も入力されていません")
		} else {
			roomUrl := r.URL.Path
			_roomId := strings.TrimPrefix(roomUrl, "/mypage/chatroom")
			roomId, _ := strconv.Atoi(_roomId)

			dbCR, err := sql.Open("mysql", query.ConStrCR)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbCR.Close()

			selectedChatroom := query.SelectChatroomById(roomId, dbCR)

			userCookie, _ := r.Cookie(session.Manager.CookieName)
			userSid, _ := url.QueryUnescape(userCookie.Value)
			userId := session.Manager.Database[userSid].SessionValue["ID"]
			newChat.UserId = userId
			newChat.Id = selectedChatroom.Id
			newChat.RoomName = selectedChatroom.RoomName
			if userId == selectedChatroom.UserId {
				//投稿主と部屋作成者が同じ場合
				newChat.UserId = userId
				newChat.Member = selectedChatroom.Member
			} else {
				//投稿主と部屋作成者が違う場合
				newChat.UserId = selectedChatroom.Member
				newChat.Member = selectedChatroom.UserId
			}
			newChat.PostDt = time.Now().UTC().Round(time.Second)

			//dbCR, err := sql.Open("mysql", query.ConStrCR)
			//if err != nil {
			//	fmt.Println(err.Error())
			//}
			//defer dbCR.Close()

			posted := query.InsertChat(newChat.Id, newChat.UserId, newChat.RoomName, newChat.Member, newChat.Chat, newChat.PostDt, dbCR)
			if posted == true {
				fmt.Fprintf(w, "投稿されました")
			} else {
				fmt.Fprintf(w, "投稿できませんでした")
			}
		}
	}
}
