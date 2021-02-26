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
		newUser.Password = r.FormValue("password")

		if newUser.UserId == "" || newUser.Password == "" {
			fmt.Fprintf(w, "IDまたはパスワードが入力されていません")
		} else {
			// データベース接続
			dbUsr, err := sql.Open("mysql", query.ConStrUsr)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer dbUsr.Close()
			insertedUser := query.InsertUser(newUser.UserId, newUser.Password, dbUsr)
			if insertedUser == true {
				t := template.Must(template.ParseFiles("./templates/resistrationcompleted.html"))
				t.ExecuteTemplate(w, "resistrationcompleted.html", nil)
			} else {
				fmt.Fprintf(w, "既に登録されているIDです")
			}
		}
	}
}
