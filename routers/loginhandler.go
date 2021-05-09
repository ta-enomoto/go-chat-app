//ログインページにアクセスがあった時のハンドラ
package routers

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
		return

	//ログインID・パスワードがポストされた時
	case "POST":
		accessingUser := new(query.User)
		accessingUser.UserId = r.FormValue("userId")
		fmt.Println(accessingUser.UserId)
		psw_string := r.FormValue("password")

		if accessingUser.UserId == "" || psw_string == "" {
			fmt.Fprintf(w, "IDまたはパスワードが入力されていません")
			return
		}

		accessingUser.Password = []byte(psw_string)
		fmt.Println(accessingUser.Password)

		dbUsr, err := sql.Open("mysql", query.ConStrUsr)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbUsr.Close()

		user := query.SelectUserById(accessingUser.UserId, dbUsr)
		fmt.Println(user)

		pswMatchOrNot := bcrypt.CompareHashAndPassword(user.Password, accessingUser.Password)

		if accessingUser.UserId == user.UserId && pswMatchOrNot == nil {
			//if文でsessionstartがうまくいった時というふうに(ブラウザで/に戻った時、sid出し直してる)
			session.Manager.SessionStart(w, r, accessingUser.UserId)
			http.Redirect(w, r, "/mypage", 301)
		} else {
			fmt.Fprintf(w, "IDまたはパスワードが間違っています。")
		}
	}
}
