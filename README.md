# ws-server

参考：https://www.cnblogs.com/ourongxin/archive/2022/02/23/15925620.html
参考：https://www.cnblogs.com/qingfj/p/15058528.html

实现了一个websocket服务。

server入口：`main.go`
client入口：`client/client.html`

测试发送消息：
```bash
curl -X POST http://localhost:8080/send -d 'you are not alone'
```