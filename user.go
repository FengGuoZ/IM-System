package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn
}

// NewUser 创建1个用户的API
func NewUser(conn net.Conn) *User {
	// 获取对端地址
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,
	}

	// 启动监听当前user channel循环的goroutine
	go user.ListenMessage()

	return user
}

// ListenMessage 监听当前user chan的方法，一但有消息，立刻发送给客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.Conn.Write([]byte(msg + "\n"))
	}
}
