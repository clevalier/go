package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewSer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("连接建立成功")

}

func (this *Server) Start() {
	//socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Printf("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		go this.Handler(conn)

	}

}
