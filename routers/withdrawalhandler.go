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
		deleteUser := new(query.USER)
		deleteUser.Id = r.FormValue("loginID")
		deleteUser.Password = r.FormValue("password")

		// データベース接続
		db, err := sql.Open("mysql", query.ConStr)
		if err != nil {
			fmt.Println(err.Error())
		}
		// deferで処理終了前に必ず接続をクローズする
		defer db.Close()

		user := query.SelectUserById(deleteUser.Id, db)
		if deleteUser.Id == user.Id && deleteUser.Password == user.Password {
			deleted := query.DeleteUserById(deleteUser.Id, db)
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
