package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer 创建1个server对象
func NewServer(ip string, port int) *Server {
	return &Server{ip, port}
}

// Handler conn业务处理逻辑
func (s *Server) Handler(conn net.Conn) {
	// 当前连接业务
	fmt.Println("连接建立成功")
}

// Start 启动服务接口
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

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
