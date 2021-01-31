package routers

import (
	"database/sql"
	"fmt"
	"goserver/query"
	"goserver/sessions"
	"html/template"
	"net/http"
	"net/url"
)

func MypageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if ok := session.Manager.SidCheck(w, r); ok {
			t := template.Must(template.ParseFiles("./templates/mypage.html"))
			// データベース接続
			db, err := sql.Open("mysql", query.ConStrCR)
			if err != nil {
				fmt.Println(err.Error())
			}
			// deferで処理終了前に必ず接続をクローズする
			defer db.Close()
			userCookie, _ := r.Cookie(session.Manager.CookieName)
			userSid, _ := url.QueryUnescape(userCookie.Value)
			userSessionVar := session.Manager.Database[userSid].SessionValue["ID"]
			chatrooms := query.SelectChatroomByUserId(userSessionVar, db)

			var Links []string

			for i := 1; i <= len(chatrooms); i++ {
				roomName := chatrooms[i].RoomName
				Links = append(Links, roomName)
				fmt.Println(Links)
			}

			t.ExecuteTemplate(w, "mypage.html", Links)
		} else {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
		}
	}
}
