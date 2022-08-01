package server

import (
	"fmt"
	"go-IM/src/user"
	"net"
	"sync"
)

type Server struct {
	// ip地址
	IP string
	// 端口
	Port int

	OnlineMap map[string]*user.User
	mapLock   sync.RWMutex

	// 消息广播channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*user.User),
		Message:   make(chan string),
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

	go server.ListenMessager()

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
	// fmt.Println("链接成功")

	user := user.NewUser(conn)

	// 防止并发读写错误
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	server.BroadCast(user, "已上线")

	// 阻塞当前handler
	select {}
}

// BroadCast 广播上线消息
func (server *Server) BroadCast(user *user.User, msg string) {
	sendMsg := "[" + user.Addr + "]" + ":" + msg

	// 往channel写数据
	server.Message <- sendMsg
}

// ListenMessager goroutine 监听在线用户，当有消息过来时发送给所有用户
func (server *Server) ListenMessager() {
	for {
		msg := <-server.Message

		server.mapLock.Lock()
		// 遍历OnlineMap 中搜集的user，向在线的user的channel发消息
		for _, user := range server.OnlineMap {
			user.C <- msg
		}
		server.mapLock.Unlock()
	}
}
