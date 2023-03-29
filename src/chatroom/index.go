package chatroom

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/mole828/gin-websocket-demo/src/chatroom/user"
	"github.com/streadway/amqp"
)

type Chatroom struct {
	users map[string]*user.User
	mqCh  *amqp.Channel
}

type Message struct {
	From  string
	Value string
}

func consumeMessages(ch *amqp.Channel, exchangeName string) (<-chan amqp.Delivery, error) {
	// 创建一个新的队列
	queue, err := ch.QueueDeclare(
		"",    // 队列名称，由服务器随机生成
		false, // 是否持久化
		true,  // 是否自动删除
		false, // 是否排他性
		false, // 是否等待服务器响应
		nil,   // 额外参数
	)
	if err != nil {
		return nil, err
	}

	// 将队列绑定到 fanout 交换机上
	err = ch.QueueBind(
		queue.Name,   // 队列名称
		"",           // 路由键
		exchangeName, // 交换机名称
		false,        // 是否等待服务器响应
		nil,          // 额外参数
	)
	if err != nil {
		return nil, err
	}

	// 消费队列中的消息
	msgs, err := ch.Consume(
		queue.Name, // 队列名称
		"",         // 路由键
		true,       // 是否自动确认消息
		false,      // 是否独占
		false,      // 是否等待服务器响应
		false,      // 额外参数
		nil,        // 额外参数
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

var exchangeName string = "chat_broadcast"

func New(ch *amqp.Channel) (*Chatroom, error) {
	msgs, err := consumeMessages(ch, exchangeName)
	if err != nil {
		return nil, err
	}
	room := &Chatroom{
		users: make(map[string]*user.User),
		mqCh:  ch,
	}
	go func() {
		for msg := range msgs {
			var message *Message = &Message{}
			err := json.Unmarshal(msg.Body, message)
			if err != nil {
				break
			}
			room.Send(message)
		}
	}()
	return room, nil
}

func (it *Chatroom) Send(message *Message) error {
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
			message := Message{
				From:  id,
				Value: string(msg),
			}
			data, err := json.Marshal(message)
			log.Printf("public: %s \n", message)
			if err != nil {
				break
			}
			it.mqCh.Publish(
				exchangeName,
				"",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        data,
				},
			)
		}
	}()
}
