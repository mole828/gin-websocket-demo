package user

import (
	"github.com/gorilla/websocket"
)

type User struct {
	conn     *websocket.Conn
	Channel  chan []byte
	onLogout []func()
}

func New(conn *websocket.Conn) *User {
	user := &User{
		conn:     conn,
		Channel:  make(chan []byte),
		onLogout: make([]func(), 0),
	}
	go func() {
		for {
			if message_type, message, err := conn.ReadMessage(); err == nil {
				if message_type == websocket.TextMessage {
					user.Channel <- message
				}
			} else {
				break
			}
		}
		for _, callback := range user.onLogout {
			callback()
		}
		conn.Close()
	}()
	return user
}

func (it *User) OnLogout(callback func()) {
	it.onLogout = append(it.onLogout, callback)
}

func (it *User) Send(data []byte) error {
	return it.conn.WriteMessage(websocket.TextMessage, data)
}
