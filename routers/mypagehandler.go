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

		//自分を相手メンバーに選んだ時の禁止処理書く

		newchatroom := new(query.Chatroom)
		newchatroom.RoomName = r.FormValue("roomName")
		newchatroom.Member = r.FormValue("memberName")

		if newchatroom.RoomName == "" || newchatroom.Member == "" {
			fmt.Fprintf(w, "メンバーまたはルーム名が入力されていません")
			return
		}

		userCookie, _ := r.Cookie(session.Manager.CookieName)
		userSid, _ := url.QueryUnescape(userCookie.Value)
		userSessionVar := session.Manager.SessionStore[userSid].SessionValue["userId"]

		if newchatroom.Member == userSessionVar {
			fmt.Fprintf(w, "自分自身をメンバーに加えることはできません。")
			return
		}

		dbUsr, err := sql.Open("mysql", query.ConStrUsr)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbUsr.Close()

		users := query.SelectAllUser(dbUsr)
		userIdExist := query.ContainsUserName(users, newchatroom.Member)

		if !userIdExist {
			fmt.Fprintf(w, "相手ユーザーが存在しません")
			return
		}

		dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbChtrm.Close()

		newChtrmInsertedToDb := query.InsertChatroom(userSessionVar, newchatroom.RoomName, newchatroom.Member, dbChtrm)
		if newChtrmInsertedToDb {
			fmt.Fprintf(w, "新しいルームが作成されました")
		} else {
			fmt.Fprintf(w, "既に登録されているルームです")
		}
	}
}
