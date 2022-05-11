package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn

	server *Server
}

// NewUser 创建1个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	// 获取对端地址
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,

		server: server,
	}

	// 启动监听当前user channel循环的goroutine
	go user.ListenMessage()

	return user
}

// Online 用户上线业务
func (u *User) Online() {
	// 用户上线，将新用户加入OnlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// 广播用户上线消息
	u.server.BroadCast(u, "已上线")
}

// Offline 用户下线业务
func (u *User) Offline() {
	// 用户下线，将用户从OnlineMap中删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// 广播用户上线消息
	u.server.BroadCast(u, "已下线")
}

// SendMessage 给当前用户发消息
func (u *User) SendMessage(msg string) {
	u.Conn.Write([]byte(msg))
}

// DoMessage 处理用户消息的业务
func (u *User) DoMessage(msg string) {
	if msg == "who" { // 在线用户查询功能
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			sendMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			u.SendMessage(sendMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[0:7] == "rename|" { // 用户名修改功能
		// 消息格式：rename|张三
		newName := strings.Split(msg, "|")[1]

		// 判断name是否存在(纯查询不用加锁)
		_, OK := u.server.OnlineMap[newName]
		if OK {
			u.SendMessage("当前用户名被占用\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMessage("您已经更新用户名：" + u.Name + "\n")
		}
	} else { // 默认广播功能
		u.server.BroadCast(u, msg)
	}
}

// ListenMessage 监听当前user chan的方法，一但有消息，立刻发送给客户端（向客户端conn写位置）
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		_, _ = u.Conn.Write([]byte(msg + "\n"))
	}
}
