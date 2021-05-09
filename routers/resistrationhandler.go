//ユーザー登録ページにアクセスがあったときのハンドラ
package routers

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goserver/query"
	"html/template"
	"net/http"
)

func ResistrationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t := template.Must(template.ParseFiles("./templates/resistration.html"))
		t.ExecuteTemplate(w, "resistration.html", nil)

	case "POST":
		newUser := new(query.User)
		newUser.UserId = r.FormValue("userId")
		psw_string := r.FormValue("password")

		if newUser.UserId == "" || psw_string == "" {
			fmt.Fprintf(w, "IDまたはパスワードが入力されていません")
			return
		}

		newUser.Password = []byte(psw_string)

		dbUsr, err := sql.Open("mysql", query.ConStrUsr)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbUsr.Close()

		users := query.SelectAllUser(dbUsr)

		userIdAlreadyExists := query.ContainsUserName(users, newUser.UserId)
		if userIdAlreadyExists {
			fmt.Fprintf(w, "既に登録されているIDです")
			return
		}

		insertedUser := query.InsertUser(newUser.UserId, newUser.Password, dbUsr)
		if insertedUser {
			t := template.Must(template.ParseFiles("./templates/resistrationcompleted.html"))
			t.ExecuteTemplate(w, "resistrationcompleted.html", nil)
		}
	}
}
