package main

import (
	"goserver/routers"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type MyMux struct {
}

func (mux MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		routers.LoginHandler(w, r)
	case "/mypage":
		routers.MypageHandler(w, r)
	case "/resistration":
		routers.ResistrationHandler(w, r)
	case "/logout":
		routers.LogoutHandler(w, r)
	case "/withdrawal":
		routers.WithdrawalHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	mux := MyMux{}
	http.ListenAndServe(":8080", mux) //監視するポートを設定します。
}

//各ページにセッションチェック追加
