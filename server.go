package main

import (
	"fmt"
	"goserver/routers"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

type MyMux struct {
}

func (mux MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var login = regexp.MustCompile(`/login`)
	var mypage = regexp.MustCompile(`^/mypage$`)
	var resistration = regexp.MustCompile(`/resistration`)
	var logout = regexp.MustCompile(`/logout`)
	var withdrawal = regexp.MustCompile(`/withdrawal`)
	var dirUnderMypage = regexp.MustCompile(`/mypage/.*`)
	url := r.URL.Path

	switch { //r.URL.Path {
	case login.MatchString(url):
		routers.LoginHandler(w, r)
	case mypage.MatchString(url):
		routers.MypageHandler(w, r)
	case resistration.MatchString(url):
		routers.ResistrationHandler(w, r)
	case logout.MatchString(url):
		routers.LogoutHandler(w, r)
	case withdrawal.MatchString(url):
		routers.WithdrawalHandler(w, r)
	case dirUnderMypage.MatchString(url):
		routers.ChatroomHandler(w, r)
	default:
		url := r.URL.Path
		fmt.Println(url)
		http.NotFound(w, r)
	}
}

func main() {
	mux := MyMux{}
	http.ListenAndServe(":8080", mux) //監視するポートを設定します。
}

//各ページにセッションチェック追加
