package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"yinkn.cn/ws-server/model"
	"yinkn.cn/ws-server/service"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Ws(c *gin.Context) {
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	client := &model.WsClient{
		Id:            conn.RemoteAddr().String(),
		User:          "test",
		Socket:        conn,
		HeartbeatTime: 0,
	}

	// 初始化
	client.Init(c)
	// 注册
	err = service.DefaultClientManager.Register(c, client)
	if err != nil {
		return
	}
	client.RegisterReceived(c)
	// 持续接收心跳
	client.ReceiveHeartBeat(c)
}

func Send(c *gin.Context) {
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	clients, err := service.DefaultClientManager.GetClients(c, []string{"test"})
	if err != nil {
		return
	}
	for _, client := range clients {
		fmt.Println("send to:", client.Id)
		client.Socket.WriteMessage(websocket.TextMessage, bs)
	}
}
