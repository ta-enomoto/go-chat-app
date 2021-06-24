package wsserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"goserver/query"
	//"goserver/sessions"
	"strconv"
	//"strings"
	"time"
)

type WsChat struct {
	Id       string    `json:"id"`
	UserId   string    `json:"userid"`
	RoomName string    `json:"roomname"`
	Member   string    `json:"member"`
	Chat     string    `json:"chat"`
	PostDt   time.Time `json:"postdt"`
	Cookie   string    `json:"cookie"`
}

type WsChatroom struct {
	id      string
	forward chan string
	Join    chan *WsClient
	Leave   chan *WsClient
	clients map[*WsClient]bool
}

type Message struct {
	Message string
}

type WsClient struct {
	Send chan string
	Room *WsChatroom
}

type MANAGER struct {
	WsRooms map[string]*WsChatroom
}

var Manager *MANAGER

//サーバー起動時にセッションマネージャも初期化
func init() {
	Manager = NewManager()
}

//init()でセッションマネージャ初期化時に使う関数
func NewManager() *MANAGER {
	database := make(map[string]*WsChatroom)
	return &MANAGER{WsRooms: database}
}

func WebSocketHandler(ws *websocket.Conn) {
	defer ws.Close()

	var chatroomJson string
	if err := websocket.Message.Receive(ws, &chatroomJson); err == nil {
		var chatroom WsChat
		json.Unmarshal([]byte(chatroomJson), &chatroom)
		roomId := chatroom.Id
		if _, exist := Manager.WsRooms[roomId]; exist {

			WsClient := &WsClient{
				Send: make(chan string),
				Room: Manager.WsRooms[roomId],
			}

			Manager.WsRooms[roomId].Join <- WsClient
			defer func() {
				Manager.WsRooms[roomId].Leave <- WsClient
				if len(Manager.WsRooms[roomId].clients) == 0 {
					delete(Manager.WsRooms, roomId)
					fmt.Println("WebSocket用ルーム削除")
				}
			}()

			go WsClient.Write(ws)
			WsClient.Read(ws)
		} else {
			WsChatroom := &WsChatroom{
				forward: make(chan string),
				Join:    make(chan *WsClient),
				Leave:   make(chan *WsClient),
				clients: make(map[*WsClient]bool),
			}

			go WsChatroom.ChatroomRun()

			var chatroom WsChat
			json.Unmarshal([]byte(chatroomJson), &chatroom)
			WsChatroom.id = chatroom.Id

			Manager.WsRooms[WsChatroom.id] = WsChatroom

			WsClient := &WsClient{
				Send: make(chan string),
				Room: WsChatroom,
			}

			WsChatroom.Join <- WsClient
			defer func() {
				WsChatroom.Leave <- WsClient
				Manager.WsRooms[WsChatroom.id].Leave <- WsClient
				if len(Manager.WsRooms[WsChatroom.id].clients) == 0 {
					delete(Manager.WsRooms, WsChatroom.id)
					fmt.Println("WebSocket用ルーム削除")
				}
			}()

			go WsClient.Write(ws)
			WsClient.Read(ws)
		}
	}
}

func (WsChatroom *WsChatroom) ChatroomRun() {
	for {
		select {
		case WsClient := <-WsChatroom.Join:
			WsChatroom.clients[WsClient] = true
			fmt.Printf("クライアントが入室しました。現在 %x 人のクライアントが存在します。\n", len(WsChatroom.clients))

		case WsClient := <-WsChatroom.Leave:
			delete(WsChatroom.clients, WsClient)
			fmt.Printf("クライアントが退出しました。現在 %x 人のクライアントが存在します。\n", len(WsChatroom.clients))

		case msg := <-WsChatroom.forward:
			fmt.Println("メッセージ受信")
			for target := range WsChatroom.clients {
				select {
				case target.Send <- msg:
					fmt.Println("メッセージ送信成功")
				default:
					fmt.Println("メッセージ送信失敗")
					delete(WsChatroom.clients, target)
				}
			}
		}
	}
}

func (wc *WsClient) Read(ws *websocket.Conn) {

	for {
		var msg string
		var Chat WsChat
		if err := websocket.Message.Receive(ws, &msg); err == nil {
			json.Unmarshal([]byte(msg), &Chat)
		} else {
			fmt.Println("json受信失敗")
			return
		}

		dbChtrm, err := sql.Open("mysql", query.ConStrChtrm)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer dbChtrm.Close()

		postedChat := new(query.Chat)
		roomId, _ := strconv.Atoi(Chat.Id)
		postedChat.Chatroom.Id = roomId
		currentChatroom := query.SelectChatroomById(roomId, dbChtrm)

		postedChat.Chatroom.Id = currentChatroom.Id
		postedChat.Chatroom.RoomName = currentChatroom.RoomName

		//本番環境ではブラウザからcookie取得→sessionからuser特定
		//cookie := Chat.Cookie
		//ここに複数cookieがある場合の処理書いておく
		//cookieValue := strings.TrimPrefix(cookie, "cookieName=")
		//userSessionVar := session.Manager.SessionStore[cookieValue].SessionValue["userId"]

		//if userSessionVar == currentChatroom.UserId {
		//投稿主と部屋作成者が同じ場合
		postedChat.Chatroom.UserId = currentChatroom.UserId
		postedChat.Chatroom.Member = currentChatroom.Member
		//} else {
		//投稿主と部屋作成者が違う場合
		//	postedChat.Chatroom.UserId = currentChatroom.Member
		//	postedChat.Chatroom.Member = currentChatroom.UserId
		//}
		postedChat.Chat = Chat.Chat
		postedChat.PostDt = time.Now().UTC().Round(time.Second)

		posted := query.InsertChat(postedChat.Chatroom.Id, postedChat.Chatroom.UserId, postedChat.Chatroom.RoomName, postedChat.Chatroom.Member, postedChat.Chat, postedChat.PostDt, dbChtrm)
		if posted {
			fmt.Println("投稿成功")
			wc.Room.forward <- msg
			continue
		}
	}
}

func (wc *WsClient) Write(ws *websocket.Conn) {
	for msg := range wc.Send {
		if err := websocket.Message.Send(ws, msg); err != nil {
			break
		}
	}
}
