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

			dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbChtrm.Close()

			selectedChatroom := query.SelectChatroomById(roomId, dbChtrm)
			fmt.Println(selectedChatroom.Id)
			Chat := query.SelectAllChatById(selectedChatroom.Id, dbChtrm)
			fmt.Println(Chat)

			t := template.Must(template.ParseFiles("./templates/mypage/chatroom.html"))
			t.ExecuteTemplate(w, "chatroom.html", Chat)
		} else {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
		}
	case "POST":
		newChat := new(query.Chat)
		newChat.Chat = r.FormValue("chat")
		if newChat.Chat == "" {
			fmt.Fprintf(w, "何も入力されていません")
		} else {
			roomUrl := r.URL.Path
			_roomId := strings.TrimPrefix(roomUrl, "/mypage/chatroom")
			roomId, _ := strconv.Atoi(_roomId)

			dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbChtrm.Close()

			selectedChatroom := query.SelectChatroomById(roomId, dbChtrm)

			userCookie, _ := r.Cookie(session.Manager.CookieName)
			userSid, _ := url.QueryUnescape(userCookie.Value)
			userId := session.Manager.Database[userSid].SessionValue["userId"]
			newChat.Chatroom.UserId = userId
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

			//dbChtrm, err := sql.Open("mysql", query.ConStrCR)
			//if err != nil {
			//	fmt.Println(err.Error())
			//}
			//defer dbChtrm.Close()

			posted := query.InsertChat(newChat.Chatroom.Id, newChat.Chatroom.UserId, newChat.Chatroom.RoomName, newChat.Chatroom.Member, newChat.Chat, newChat.PostDt, dbChtrm)
			if posted == true {
				fmt.Fprintf(w, "投稿されました")
			} else {
				fmt.Fprintf(w, "投稿できませんでした")
			}
		}
	}
}
