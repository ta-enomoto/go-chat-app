package main

import (
		"fmt"
		"net/http"
		"html/template"
		"goserver/conf"
		"goserver/query"
		"database/sql"
		_ "github.com/go-sql-driver/mysql"
)

type MyMux struct {
}

type LOGINUSER struct {
	Id string
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
			accessing_user := new(LOGINUSER)
			accessing_user.Id = r.FormValue("login_id")
			accessing_user.Password = r.FormValue("password")
			fmt.Println(accessing_user.Id, accessing_user.Password)

			// 設定ファイルを読み込む
			confDB, err := conf.ReadConfDB()
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

			user := query.SelectUserById(accessing_user.Id, db)

			if accessing_user.Id == user.Id && accessing_user.Password == user.Password {
					t := template.Must(template.ParseFiles("./templates/mypage.html"))
					values := map[string]string{
							"login_id": r.FormValue("login_id"),
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