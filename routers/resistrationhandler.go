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
		newUser := new(query.USER)
		newUser.Id = r.FormValue("loginID")
		newUser.Password = r.FormValue("password")

		if newUser.Id == "" || newUser.Password == "" {
			fmt.Fprintf(w, "IDまたはパスワードが入力されていません")
		} else {
			// データベース接続
			db, err := sql.Open("mysql", query.ConStr)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer db.Close()
			insertedUser := query.InsertUser(newUser.Id, newUser.Password, db)
			if insertedUser == true {
				t := template.Must(template.ParseFiles("./templates/resistrationcompleted.html"))
				t.ExecuteTemplate(w, "resistrationcompleted.html", nil)
			} else {
				fmt.Fprintf(w, "既に登録されているIDです")
			}
		}
	}
}
