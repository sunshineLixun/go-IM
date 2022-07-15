package server

import (
	"fmt"
	"net"
)

type Server struct {
	// ip地址
	IP string
	// 端口
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:   ip,
		Port: port,
	}
	return server
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// 程序结束前关闭服务
	defer listener.Close()

	// 死循环获取链接状态
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener Accept err:", err)
			// 如果这次失败了，跳过这次，继续监听其他链接
			continue
		}

		// 开启go程 处理链接成功结果
		go server.Handler(conn)
	}
}

func (server *Server) Handler(conn net.Conn) {
	fmt.Println("链接成功")
}
