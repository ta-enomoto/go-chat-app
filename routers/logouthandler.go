package routers

import (
	"goserver/sessions"
	"html/template"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		session.Manager.DeleteSessionFromStore(w, r)
		t := template.Must(template.ParseFiles("./templates/logout.html"))
		t.ExecuteTemplate(w, "logout.html", nil)
	}
}
