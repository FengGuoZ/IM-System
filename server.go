package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表(全局表)
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 用于广播的channel
	Message chan string
}

// NewServer 创建1个server对象
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string), // 服务器的广播通道
	}
}

// ListenMessage 监听Message广播消息通道的goroutine，一旦有消息就发送给所有在线用户
func (s *Server) ListenMessage() {
	for {
		// 阻塞读取Message通道
		msg := <-s.Message

		// 将消息发送给所有在线User
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// BroadCast 广播用户上线的方法
func (s *Server) BroadCast(u *User, msg string) {
	// 拼接出完整消息
	sendMsg := "[" + u.Addr + "]" + u.Name + ":" + msg
	s.Message <- sendMsg
}

// Handler conn业务处理逻辑的go程(连接建立时的处理 + 循环从客户端conn读数据)
func (s *Server) Handler(conn net.Conn) {
	// 当前连接业务
	// 新建1个用户,默认以userAddr为用户名
	user := NewUser(conn, s)

	user.Online()

	// 接受客户端传递的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				// 对端关闭socket，用户下线
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn read error")
				return
			}

			// 提取用户输入的信息（去除末尾\n）
			msg := string(buf[:n-1])

			// 用户针对msg进行处理
			user.DoMessage(msg)
		}
	}()

	// 当前handler阻塞
	select {}
}

// Start 启动服务接口
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket later
	defer listener.Close()

	// 启动监听Message的goroutine
	go s.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			return
		}

		// do handler
		go s.Handler(conn)
	}

}
