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
	sessionValue map[string]string
	manager      *Manager
}

type Manager struct {
	Database    map[interface{}]*Session
	cookieName  string
	maxlifetime int64
}

func NewManager(cookieName string, maxlifetime int64) *Manager {
	database := make(map[interface{}]*Session)
	return &Manager{Database: database, cookieName: cookieName, maxlifetime: maxlifetime}
}

//新規sid発行
func (manager *Manager) NewSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//新規セッションの生成
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request, userId string) (session Session) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.NewSessionId()
		fmt.Println(sid)
		session := manager.NewSession(sid, userId)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w, &cookie)
		manager.Database[sid] = session
		fmt.Println(manager.Database)
	} //else {  もしsidがdatabaseに登録されていたらokの処理
	//sid, _ := url.QueryUnescape(cookie.Value)
	//session := manager.SessionRead(sid, userId)
	//_ = session
	//}
	return
}

func (manager *Manager) NewSession(sid string, userId string) *Session {
	sv := make(map[string]string)
	sv["ID"] = userId
	newSession := &Session{sid: sid, timeAccessed: time.Now(), sessionValue: sv}
	return newSession
}

func (manager *Manager) SidCheck(w http.ResponseWriter, r *http.Request) bool {
	clientCookie, err := r.Cookie("cookieName")
	if err != nil {
		return false
	} else {
		clientSid, _ := url.QueryUnescape(clientCookie.Value)
		if _, ok := manager.Database[clientSid]; ok {
			if (manager.Database[clientSid].timeAccessed.Unix() + manager.maxlifetime) > time.Now().Unix() {
				manager.Database[clientSid].timeAccessed = time.Now()
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}

}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) error {
	clientCookie, err := r.Cookie("cookieName")
	if err != nil {
		fmt.Println(err.Error())
	}
	clientSid, _ := url.QueryUnescape(clientCookie.Value)
	delete(manager.Database, clientSid)
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
