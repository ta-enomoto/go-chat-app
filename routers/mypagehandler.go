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
)

func MypageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	/*アクセスがあった時の処理。自身で作成したルームと他人が作成したルームを別々に取得して、
	ルーム一覧のスライスをつくってテンプレに渡す*/
	case "GET":
		if ok := session.Manager.SessionIdCheck(w, r); ok {
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
		} else {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
		}
	/*新しいルーム作成のポストがあった時の処理。ルーム名と相手メンバーを指定する
	同名のルーム名は、相手メンバー異なる場合のみ有効。*/
	case "POST":
		if ok := session.Manager.SessionIdCheck(w, r); ok {

			newchatroom := new(query.Chatroom)
			newchatroom.RoomName = r.FormValue("roomName")
			newchatroom.Member = r.FormValue("memberName")

			if newchatroom.RoomName == "" || newchatroom.Member == "" {
				fmt.Fprintf(w, "メンバーまたはルーム名が入力されていません")
			} else {

				userCookie, _ := r.Cookie(session.Manager.CookieName)
				userSid, _ := url.QueryUnescape(userCookie.Value)
				userSessionVar := session.Manager.SessionStore[userSid].SessionValue["userId"]

				dbUsr, err := sql.Open("mysql", query.ConStrUsr)
				if err != nil {
					fmt.Println(err.Error())
				}
				defer dbUsr.Close()

				users := query.SelectAllUser(dbUsr)
				userIdExist := query.ContainsUserName(users, newchatroom.Member)

				if userIdExist == true {
					dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
					if err != nil {
						fmt.Println(err.Error())
					}
					defer dbChtrm.Close()

					inserted := query.InsertChatroom(userSessionVar, newchatroom.RoomName, newchatroom.Member, dbChtrm)
					if inserted == true {
						fmt.Fprintf(w, "新しいルームが作成されました")
					} else {
						fmt.Fprintf(w, "既に登録されているルームです")
					}
				} else {
					fmt.Fprintf(w, "ユーザーが存在しません")
				}
			}
		} else {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
		}
	}
}
