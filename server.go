package main

import (
	"database/sql"
	"fmt"
	"goserver/config"
	"goserver/query"
	"goserver/sessions"
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

var manager *session.Manager

func init() {
	_manager := session.NewManager("cookieName", 60)
	manager = _manager
}

func (mux MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		switch r.Method {
		case "GET":
			t := template.Must(template.ParseFiles("./templates/login.html"))
			t.ExecuteTemplate(w, "login.html", nil)
		case "POST":
			//ログイン判定
			accessingUser := new(LOGINUSER)
			accessingUser.Id = r.FormValue("loginID")
			accessingUser.Password = r.FormValue("password")
			fmt.Println(accessingUser.Id, accessingUser.Password)

			if accessingUser.Id == "" || accessingUser.Password == "" {
				fmt.Fprintf(w, "IDまたはパスワードが入力されていません")
			} else {
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
					//if文でsessionstartがうまくいった時というふうに(ブラウザで/に戻った時、sid出し直してる)
					manager.SessionStart(w, r, accessingUser.Id)
					http.Redirect(w, r, "/mypage", 301)
				} else {
					fmt.Fprintf(w, "IDまたはパスワードが間違っています。")
				}
			}
		}
	case "/mypage":
		if ok := manager.SidCheck(w, r); ok {
			t := template.Must(template.ParseFiles("./templates/mypage.html"))
			t.ExecuteTemplate(w, "mypage.html", nil)
		} else {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
		}
	case "/resistration":
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
				inserted := query.InsertUser(newUser.Id, newUser.Password, db)
				if inserted == true {
					//ログイン完了時の処理
					t := template.Must(template.ParseFiles("./templates/resistrationcompleted.html"))
					t.ExecuteTemplate(w, "resistrationcompleted.html", nil)
				} else {
					fmt.Fprintf(w, "既に登録されているIDです")
				}
			}
		}
	case "/logout":
		switch r.Method {
		case "GET":
			manager.SessionDestroy(w, r)
			t := template.Must(template.ParseFiles("./templates/logout.html"))
			t.ExecuteTemplate(w, "logout.html", nil)
		}
	case "/withdrawal":
		switch r.Method {
		case "GET":
			t := template.Must(template.ParseFiles("./templates/withdrawal.html"))
			t.ExecuteTemplate(w, "withdrawal.html", nil)
		case "POST":
			deleteUser := new(query.USER)
			deleteUser.Id = r.FormValue("loginID")
			deleteUser.Password = r.FormValue("password")

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

			user := query.SelectUserById(deleteUser.Id, db)
			if deleteUser.Id == user.Id && deleteUser.Password == user.Password {
				deleted := query.DeleteUserById(deleteUser.Id, db)
				if deleted == true {
					//削除完了時の処理
					manager.SessionDestroy(w, r)
					t := template.Must(template.ParseFiles("./templates/withdrawalcompleted.html"))
					t.ExecuteTemplate(w, "withdrawalcompleted.html", nil)
				}
			} else {
				fmt.Fprintf(w, "IDまたはパスワードが間違っています")
			}
		}
	default:
		http.NotFound(w, r)
	}
}

func main() {
	mux := MyMux{}
	http.ListenAndServe(":8080", mux) //監視するポートを設定します。
}

//各ページにセッションチェック追加
