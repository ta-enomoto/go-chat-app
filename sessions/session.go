package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

//セッション構造体
type Session struct {
	sid          string
	timeAccessed time.Time
	SessionValue map[string]string
	Manager      *MANAGER
}

type MANAGER struct {
	Database    map[interface{}]*Session
	CookieName  string
	maxlifetime int64
}

var Manager *MANAGER

func init() {
	_manager := NewManager("cookieName", 60)
	Manager = _manager
}

func NewManager(cookieName string, maxlifetime int64) *MANAGER {
	database := make(map[interface{}]*Session)
	return &MANAGER{Database: database, CookieName: cookieName, maxlifetime: maxlifetime}
}

//新規sid発行
func (Manager *MANAGER) NewSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//新規セッションの生成
func (Manager *MANAGER) SessionStart(w http.ResponseWriter, r *http.Request, userId string) (session Session) {
	cookie, err := r.Cookie(Manager.CookieName)
	if err != nil || cookie.Value == "" {
		sid := Manager.NewSessionId()
		fmt.Println(sid)
		session := Manager.NewSession(sid, userId)
		cookie := http.Cookie{Name: Manager.CookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(Manager.maxlifetime)}
		http.SetCookie(w, &cookie)
		Manager.Database[sid] = session
		fmt.Println(Manager.Database)
	} //else {  もしsidがdatabaseに登録されていたらokの処理
	//sid, _ := url.QueryUnescape(cookie.Value)
	//session := manager.SessionRead(sid, userId)
	//_ = session
	//}
	return
}

func (Manager *MANAGER) NewSession(sid string, userId string) *Session {
	sv := make(map[string]string)
	sv["ID"] = userId
	newSession := &Session{sid: sid, timeAccessed: time.Now(), SessionValue: sv}
	return newSession
}

func (Manager *MANAGER) SidCheck(w http.ResponseWriter, r *http.Request) bool {
	clientCookie, err := r.Cookie(Manager.CookieName)
	if err != nil {
		return false
	} else {
		clientSid, _ := url.QueryUnescape(clientCookie.Value)
		if _, ok := Manager.Database[clientSid]; ok {
			if (Manager.Database[clientSid].timeAccessed.Unix() + Manager.maxlifetime) > time.Now().Unix() {
				Manager.Database[clientSid].timeAccessed = time.Now()
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}

}

func (Manager *MANAGER) SessionDestroy(w http.ResponseWriter, r *http.Request) error {
	clientCookie, err := r.Cookie(Manager.CookieName)
	if err != nil {
		fmt.Println(err.Error())
	}
	clientSid, _ := url.QueryUnescape(clientCookie.Value)
	delete(Manager.Database, clientSid)
	clientCookie.MaxAge = -1
	http.SetCookie(w, clientCookie)
	return nil
}

//func (manager *Manager) SessionRead(sid string, userId string) *Session {
//	if existingSession, ok := manager.Database[sid]; ok {
//		return existingSession
//	} else {
//		newsess := manager.NewSession(sid, userId)
//		return newsess
//	}
//}
