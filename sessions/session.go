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

//セッションマネージャ構造体
type MANAGER struct {
	SessionStore map[interface{}]*Session
	CookieName   string
	maxlifetime  int64
}

var Manager *MANAGER

//サーバー起動時にセッションマネージャも初期化
func init() {
	Manager = NewManager("cookieName", 3600)

}

//init()でセッションマネージャ初期化時に使う関数
func NewManager(cookieName string, maxlifetime int64) *MANAGER {
	database := make(map[interface{}]*Session)
	return &MANAGER{SessionStore: database, CookieName: cookieName, maxlifetime: maxlifetime}
}

//新規セッションを開始する関数(セッション変数はマネージャのDatabaseマップに保存)
func (Manager *MANAGER) SessionStart(w http.ResponseWriter, r *http.Request, userId string) (session *Session) {
	cookie, err := r.Cookie(Manager.CookieName)
	if err != nil || cookie.Value == "" {
		sid := Manager.NewSessionId()
		session := Manager.NewSession(sid, userId)
		cookie := http.Cookie{Name: Manager.CookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(Manager.maxlifetime)}

		http.SetCookie(w, &cookie)
		Manager.SessionStore[sid] = session
	}
	return
}

//新規セッション”id”を発行する関数(SessionStart関数で使用)
func (Manager *MANAGER) NewSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//新規セッション(マップ)を作成するための関数(SessionStart関数で使用)
func (Manager *MANAGER) NewSession(sid string, userId string) (session *Session) {
	sv := make(map[string]string)

	sv["userId"] = userId
	return &Session{sid: sid, timeAccessed: time.Now(), SessionValue: sv}
}

//クライアントのクッキーid(セッションid)が、マネージャのDatabaseマップに登録されてるかチェックする関数
func (Manager *MANAGER) SessionIdCheck(w http.ResponseWriter, r *http.Request) bool {
	clientCookie, err := r.Cookie(Manager.CookieName)
	if err != nil {
		return false
	} else {
		clientSid, _ := url.QueryUnescape(clientCookie.Value)
		if _, ok := Manager.SessionStore[clientSid]; ok {
			if (Manager.SessionStore[clientSid].timeAccessed.Unix() + Manager.maxlifetime) > time.Now().Unix() {
				Manager.SessionStore[clientSid].timeAccessed = time.Now()
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}

}

//ログアウト時、セッションを削除マネージャのDatabaseマップから削除する関数
func (Manager *MANAGER) DeleteSessionFromStore(w http.ResponseWriter, r *http.Request) error {
	clientCookie, err := r.Cookie(Manager.CookieName)
	if err != nil {
		fmt.Println(err.Error())
	}
	clientSid, _ := url.QueryUnescape(clientCookie.Value)
	delete(Manager.SessionStore, clientSid)

	clientCookie.MaxAge = -1
	http.SetCookie(w, clientCookie)
	return nil
}
