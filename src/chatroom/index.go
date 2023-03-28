package chatroom

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/mole828/gin-websocket-demo/src/chatroom/user"
)

type Chatroom struct {
	users map[string]*user.User
}

func New() *Chatroom {
	return &Chatroom{
		users: make(map[string]*user.User),
	}
}

type Message struct {
	From  string
	Value string
}

func (it *Chatroom) Send(message Message) error {
	log.Printf("Send(%+v)\n", message)
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	var multi error
	for _, user := range it.users {
		err = user.Send(data)
		multi = multierror.Append(multi, err)
	}
	return multi
}

func (it *Chatroom) Join(user *user.User) {
	id := uuid.NewString()
	log.Printf("Join(*) id=%s", id)
	it.users[id] = user
	user.OnLogout(func() {
		delete(it.users, id)
	})
	go func() {
		for msg := range user.Channel {
			it.Send(Message{
				From:  id,
				Value: string(msg),
			})
		}
	}()
}
