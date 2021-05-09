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

func WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t := template.Must(template.ParseFiles("./templates/withdrawal.html"))
		t.ExecuteTemplate(w, "withdrawal.html", nil)

	//削除するユーザーのID・パスワードがポストされたときの処理
	case "POST":
		deleteUser := new(query.User)
		deleteUser.UserId = r.FormValue("userId")
		psw_string := r.FormValue("password")
		deleteUser.Password = []byte(psw_string)

		dbUsr, err := sql.Open("mysql", query.ConStrUsr)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbUsr.Close()

		user := query.SelectUserById(deleteUser.UserId, dbUsr)

		pswMatchOrNot := bcrypt.CompareHashAndPassword(user.Password, deleteUser.Password)

		if deleteUser.UserId == user.UserId && pswMatchOrNot == nil {
			fmt.Fprintf(w, "IDまたはパスワードが間違っています")
			return
		}
		userDeletedFromDb := query.DeleteUserById(deleteUser.UserId, dbUsr)
		if userDeletedFromDb {
			session.Manager.DeleteSessionFromStore(w, r)
			t := template.Must(template.ParseFiles("./templates/withdrawalcompleted.html"))
			t.ExecuteTemplate(w, "withdrawalcompleted.html", nil)
		}
	}
}
