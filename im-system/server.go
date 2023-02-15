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
	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//消息广播的channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("连接建立成功")
	user := NewUser(conn, this)
	user.Online()
	//用户上线，将用户加入到onlineMap中
	//this.mapLock.Lock()
	//this.OnlineMap[user.Name] = user
	//this.mapLock.Unlock()
	//广播当前用户上线消息
	//this.BroadCast(user, "已上线")
	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				//this.BroadCast(user, "已经下线")
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//提取用户的消息（去除"\n"）
			msg := string(buf[:n-1])
			//this.BroadCast(user, msg)
			user.DoMessage(msg)
		}
	}()
	//当前handler阻塞
	select {}
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
	//启动监听Message的goroutine
	go this.ListenMessager()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		go this.Handler(conn)
	}
}

//新增广播消息方法

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

//新增监听广播消息channel方法,一旦有消息就发送给全部在线的user

func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		//将msg发给全部在线的user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()

	}
}
