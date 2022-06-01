package model

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const MessageTypeRegister = "register"
const MessageTypeHeartbeat = "heartbeat" // 心跳
const MessageTypeClose = "close"         // 关闭

type WsClient struct {
	Id            string
	User          string
	Socket        *websocket.Conn
	HeartbeatTime int //心跳
	ctx           context.Context
	cancelFn      context.CancelFunc
}

type wsReq struct {
	Type string `json:"type"`
}

type wsReply struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func (c *WsClient) Init(ctx context.Context) {
	ctx, cancelFn := context.WithCancel(ctx)
	c.ctx = ctx
	c.cancelFn = cancelFn
	c.HeartbeatTime = 0
}

func (c *WsClient) Close() {
	c.cancelFn()
}

func (c *WsClient) RegisterReceived(ctx context.Context) {
	reply, _ := json.Marshal(wsReply{
		Type: MessageTypeRegister,
		Msg:  "ok",
	})
	c.Socket.WriteMessage(websocket.TextMessage, reply)
}

func (c *WsClient) ReceiveHeartBeat(ctx context.Context) {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 接收消息
		}
		// 读取消息
		_, body, err := c.Socket.ReadMessage()
		if err != nil {
			break
		}

		var msg wsReq
		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Println(err)
			continue
		}

		switch msg.Type { // 消息类型
		case MessageTypeHeartbeat:
			// 刷新连接时间
			c.HeartbeatTime = int(time.Now().Unix())

			reply, _ := json.Marshal(wsReply{
				Type: MessageTypeHeartbeat,
				Msg:  "ok",
			})
			// 回复心跳
			err = c.Socket.WriteMessage(websocket.TextMessage, reply)
			if err != nil {
				log.Println(err)
			}
		case MessageTypeClose:
			c.Close()
		}
	}
}
