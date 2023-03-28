package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mole828/gin-websocket-demo/src/chatroom"
	"github.com/mole828/gin-websocket-demo/src/chatroom/user"
)

type Msg struct {
	From    string
	Message string
}

func main() {
	app := gin.New()
	app.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"health": true,
		})
	})
	chatroom := chatroom.New()
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
	app.Run(":8080")
}
