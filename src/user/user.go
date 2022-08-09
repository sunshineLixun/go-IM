package user

import "net"

type User struct {
	Name string
	Addr string
	// channel
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动goroutine监听当前user channel消息
	go user.ListenMessage()

	return user
}

func (user *User) ListenMessage() {
	for {
		// 从go程中获取到上线消息
		msg := <-user.C

		_, _ = user.conn.Write([]byte(msg + "\n"))

	}
}
