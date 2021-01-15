package main

import (
	"database/sql"
	"fmt"
	"goserver/config"
	"goserver/query"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type MyMux struct {
}

type LOGINUSER struct {
	Id       string
	Password string
}

func (mux MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		t := template.Must(template.ParseFiles("./templates/login.html"))
		t.ExecuteTemplate(w, "login.html", nil)
		return
	case "/mypage":
		//ログイン判定
		accessingUser := new(LOGINUSER)
		accessingUser.Id = r.FormValue("loginID")
		accessingUser.Password = r.FormValue("password")
		fmt.Println(accessingUser.Id, accessingUser.Password)

		// 設定ファイルを読み込む
		confDB, err := config.ReadConfDB()
		if err != nil {
			fmt.Println(err.Error())
		}
		// 設定値から接続文字列を生成
		conStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", confDB.User, confDB.Pass, confDB.Host, confDB.Port, confDB.DbName, confDB.Charset)
		// データベース接続
		db, err := sql.Open("mysql", conStr)
		if err != nil {
			fmt.Println(err.Error())
		}
		// deferで処理終了前に必ず接続をクローズする
		defer db.Close()

		user := query.SelectUserById(accessingUser.Id, db)

		if accessingUser.Id == user.Id && accessingUser.Password == user.Password {
			t := template.Must(template.ParseFiles("./templates/mypage.html"))
			values := map[string]string{
				"loginID": r.FormValue("loginID"),
			}
			t.ExecuteTemplate(w, "mypage.html", values)
			return
		} else {
			fmt.Fprintf(w, "IDまたはパスワードが間違っています。")
		}
	}
	http.NotFound(w, r)
}

func main() {
	mux := MyMux{}
	http.ListenAndServe(":8080", mux) //監視するポートを設定します。
}
