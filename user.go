package main

import "net"

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

// 处理用户消息的业务
func (u *User) DoMessage(msg string) {
	u.server.BroadCast(u, msg)
}

// ListenMessage 监听当前user chan的方法，一但有消息，立刻发送给客户端（向客户端conn写位置）
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		_, _ = u.Conn.Write([]byte(msg + "\n"))
	}
}
