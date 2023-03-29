package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mole828/gin-websocket-demo/src/chatroom"
	"github.com/mole828/gin-websocket-demo/src/chatroom/user"
	"github.com/streadway/amqp"
)

type Msg struct {
	From    string
	Message string
}

func Port_already_in_use(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	} else {
		ln.Close()
		return false
	}
}

func RunServer(ch *amqp.Channel) {
	app := gin.New()
	app.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"health": true,
		})
	})
	chatroom, _ := chatroom.New(ch)
	app.GET("/ws", func(ctx *gin.Context) {
		upGrande := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
		}
		conn, err := upGrande.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			return
		}
		chatroom.Join(user.New(conn))
	})
	port := 8080
	for Port_already_in_use(port) {
		port += 1
	}
	app.Run(fmt.Sprintf(":%d", port))
}

func main() {
	conn, _ := amqp.Dial("amqp://golang:golangpass@www.moles.top:5672/")
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	RunServer(ch)
}
