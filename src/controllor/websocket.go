package controllor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var clients = make(map[*websocket.Conn]bool)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CreateWs(c *gin.Context) {

	client, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	clients[client] = true
	defer func() {
		client.Close()
		delete(clients, client)
	}()
	for {
		// 监听接受信息
		mt, msg, err := client.ReadMessage()
		if err == nil {

			logrus.Debugln(mt, msg)

			BroadcastMsg([]byte("hello"))
			//sendMsg(mt, message)
		} else {
			break
		}
	}
}

func BroadcastMsg(message []byte) {
	for client := range clients {
		client.WriteMessage(1, message)
	}
}
