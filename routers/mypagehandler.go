package routers

import (
	"fmt"
	"goserver/sessions"
	"html/template"
	"net/http"
)

func MypageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if ok := session.Manager.SidCheck(w, r); ok {
			t := template.Must(template.ParseFiles("./templates/mypage.html"))
			t.ExecuteTemplate(w, "mypage.html", nil)
		} else {
			fmt.Fprintf(w, "セッションの有効期限が切れています")
		}
	}
}
