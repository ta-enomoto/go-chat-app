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
			dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
			if err != nil {
				fmt.Println(err.Error())
			}
			// deferで処理終了前に必ず接続をクローズする
			defer dbChtrm.Close()
			userCookie, _ := r.Cookie(session.Manager.CookieName)
			userSid, _ := url.QueryUnescape(userCookie.Value)
			userSessionVar := session.Manager.Database[userSid].SessionValue["userId"]
			chatroomsFromUserId := query.SelectAllChatroomsByUserId(userSessionVar, dbChtrm)
			chatroomsFromMember := query.SelectAllChatroomsByMember(userSessionVar, dbChtrm)

			var Links = append(chatroomsFromUserId, chatroomsFromMember...)

			/*for i := 0; i <= len(chatroomsFromUserId)-1; i++ {
				roomName1 := chatroomsFromUserId[i].RoomName + " with " + chatroomsFromUserId[i].Member
				Links = append(Links, roomName1)
			}*/

			/*for i := 0; i <= len(chatroomsFromMember)-1; i++ {
				roomName2 := chatroomsFromMember[i].RoomName + " with " + chatroomsFromMember[i].UserId
				Links = append(Links, roomName2)
			}*/

			t.ExecuteTemplate(w, "mypage.html", Links)
		} else {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
		}
	case "POST":
		if ok := session.Manager.SidCheck(w, r); ok {
			newchatroom := new(query.Chatroom)
			newchatroom.RoomName = r.FormValue("roomName")
			newchatroom.Member = r.FormValue("memberName")
			if newchatroom.RoomName == "" || newchatroom.Member == "" {
				fmt.Fprintf(w, "IDまたはルーム名が入力されていません")
			} else {
				userCookie, _ := r.Cookie(session.Manager.CookieName)
				userSid, _ := url.QueryUnescape(userCookie.Value)
				userSessionVar := session.Manager.Database[userSid].SessionValue["userId"]

				//ユーザー名の存在チェック
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
