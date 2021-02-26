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

func WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t := template.Must(template.ParseFiles("./templates/withdrawal.html"))
		t.ExecuteTemplate(w, "withdrawal.html", nil)
	case "POST":
		deleteUser := new(query.User)
		deleteUser.UserId = r.FormValue("userId")
		deleteUser.Password = r.FormValue("password")

		// データベース接続
		dbUsr, err := sql.Open("mysql", query.ConStrUsr)
		if err != nil {
			fmt.Println(err.Error())
		}
		// deferで処理終了前に必ず接続をクローズする
		defer dbUsr.Close()

		user := query.SelectUserById(deleteUser.UserId, dbUsr)
		if deleteUser.UserId == user.UserId && deleteUser.Password == user.Password {
			deleted := query.DeleteUserById(deleteUser.UserId, dbUsr)
			if deleted == true {
				//削除完了時の処理
				session.Manager.SessionDestroy(w, r)
				t := template.Must(template.ParseFiles("./templates/withdrawalcompleted.html"))
				t.ExecuteTemplate(w, "withdrawalcompleted.html", nil)
			}
		} else {
			fmt.Fprintf(w, "IDまたはパスワードが間違っています")
		}
	}
}
