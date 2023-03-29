package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

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

type AmqpConfig struct {
	Uri string
}

type Config struct {
	Amqp AmqpConfig
}

func ReadConfig() (*Config, error) {
	fmt.Println(os.Getwd())
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件大小
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := fi.Size()

	// 创建足够大的字节数组来存储文件内容
	data := make([]byte, size)

	// 读取文件内容
	n, err := file.Read(data)
	if err != nil {
		return nil, err
	}
	fmt.Println("read:", n)

	s := string(data)
	fmt.Println(s)
	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	fmt.Println(config.Amqp.Uri)
	return config, nil
}

func main() {
	config, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	conn, _ := amqp.Dial(config.Amqp.Uri)
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	RunServer(ch)
}
