package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Msg struct {
	Message string
}

func f(x int) int {
	if x < 1 {
		return x
	}
	return x + f(x-1)
}

func main() {
	app := gin.New()
	app.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"health": true,
		})
	})

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

		go func() {
			defer conn.Close()
			for {
				//
				messageType, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println(err)
					break
				}
				fmt.Printf("ack: %d, %s \n", messageType, message)
				v := gin.H{}
				err = json.Unmarshal(message, &v)
				if err != nil {
					break
				}
				fmt.Println(v)
			}
		}()
	})
	app.Run(":8080")
}
