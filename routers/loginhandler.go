//ログインページにアクセスがあった時のハンドラ
package routers

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goserver/query"
	"goserver/sessions"
	"html/template"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t := template.Must(template.ParseFiles("./templates/login.html"))
		t.ExecuteTemplate(w, "login.html", nil)

	//ログインID・パスワードがポストされた時
	case "POST":
		accessingUser := new(query.User)
		accessingUser.UserId = r.FormValue("userId")
		accessingUser.Password = r.FormValue("password")

		if accessingUser.UserId == "" || accessingUser.Password == "" {
			fmt.Fprintf(w, "IDまたはパスワードが入力されていません")
			return
		}
		dbUsr, err := sql.Open("mysql", query.ConStrUsr)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbUsr.Close()

		user := query.SelectUserById(accessingUser.UserId, dbUsr)
		if accessingUser.UserId == user.UserId && accessingUser.Password == user.Password {
			//if文でsessionstartがうまくいった時というふうに(ブラウザで/に戻った時、sid出し直してる)
			session.Manager.SessionStart(w, r, accessingUser.UserId)
			http.Redirect(w, r, "/mypage", 301)
		} else {
			fmt.Fprintf(w, "IDまたはパスワードが間違っています。")
		}
		/*
			if accessingUser.UserId == "" || accessingUser.Password == "" {
				fmt.Fprintf(w, "IDまたはパスワードが入力されていません")
			} else {
				dbUsr, err := sql.Open("mysql", query.ConStrUsr)
				if err != nil {
					fmt.Println(err.Error())
				}
				defer dbUsr.Close()

				user := query.SelectUserById(accessingUser.UserId, dbUsr)
				if accessingUser.UserId == user.UserId && accessingUser.Password == user.Password {
					//if文でsessionstartがうまくいった時というふうに(ブラウザで/に戻った時、sid出し直してる)
					session.Manager.SessionStart(w, r, accessingUser.UserId)
					http.Redirect(w, r, "/mypage", 301)
				} else {
					fmt.Fprintf(w, "IDまたはパスワードが間違っています。")
				}
		*/
	}
}
