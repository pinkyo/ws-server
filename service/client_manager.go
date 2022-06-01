package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"yinkn.cn/ws-server/model"
)

// 消息类型
const (
	HeartbeatCheckTime = 10 // 心跳检测几秒检测一次
	HeartbeatTime      = 20 // 心跳距离上一次的最大时间
)

var DefaultClientManager = NewClientManager()

type ClientManager interface {
	Register(ctx context.Context, client *model.WsClient) error
	UnRegister(ctx context.Context, client *model.WsClient) error
	GetClients(ctx context.Context, users []string) ([]*model.WsClient, error)
	Close()
}

type clientManager struct {
	rwMu sync.RWMutex

	clients  map[string]*model.WsClient
	users    map[string]map[string]struct{} // user -> []clientId
	ctx      context.Context
	cancelFn context.CancelFunc
}

func NewClientManager() ClientManager {
	ctx, cancelFn := context.WithCancel(context.Background())
	result := clientManager{
		clients:  make(map[string]*model.WsClient),
		users:    make(map[string]map[string]struct{}),
		ctx:      ctx,
		cancelFn: cancelFn,
	}
	go result.checkHeartbeat()
	return &result
}

func (m *clientManager) Register(ctx context.Context, client *model.WsClient) error {
	m.rwMu.Lock()
	defer m.rwMu.Unlock()

	m.clients[client.Id] = client
	if _, ok := m.users[client.User]; !ok {
		m.users[client.User] = make(map[string]struct{})
	}
	m.users[client.User][client.Id] = struct{}{}
	return nil
}

func (m *clientManager) UnRegister(ctx context.Context, client *model.WsClient) error {
	m.rwMu.Lock()
	defer m.rwMu.Unlock()

	client.Close()
	delete(m.clients, client.Id)
	if idSet, ok := m.users[client.User]; ok {
		delete(idSet, client.Id)
	}
	return nil
}

func (m *clientManager) GetClients(ctx context.Context, users []string) ([]*model.WsClient, error) {
	m.rwMu.RLock()
	defer m.rwMu.RUnlock()

	var clients []*model.WsClient
	for _, user := range users {
		for id := range m.users[user] {
			clients = append(clients, m.clients[id])
		}
	}
	return clients, nil
}

func (m *clientManager) Close() {
	m.cancelFn()
}

// 维持心跳
func (m *clientManager) checkHeartbeat() {
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			// 检测心跳
		}

		// 获取所有的Clients
		m.rwMu.Lock()
		var clients []*model.WsClient
		for _, c := range m.clients {
			clients = append(clients, c)
		}
		m.rwMu.Unlock()

		fmt.Printf("check heartbeat, clients: %v\n", clients)
		for _, c := range clients {
			if time.Now().Unix()-HeartbeatTime > int64(c.HeartbeatTime) {
				m.UnRegister(m.ctx, c)
			}
		}
		time.Sleep(time.Second * HeartbeatCheckTime)
	}
}
