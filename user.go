package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string // 用于接受广播消息
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

	// 启动循环监听广播channel的Go程
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
			// 删除旧用户名，换上新用户名
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMessage("您已经更新用户名：" + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[0:3] == "to|" {
		// 消息格式：to|张三|消息内容

		// 获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMessage("消息格式不正确，正确格式为：\"to|目标用户名|消息内容\"\n")
			return
		}

		// 根据用户名，得到对方User对象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if ok == false {
			u.SendMessage("目标用户名不存在\n")
		}

		// 获取消息内容，发送给对端
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMessage("无消息内容，请重新发送\n")
			return
		}
		remoteUser.SendMessage(u.Name + "对您说：" + content + "\n")

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
