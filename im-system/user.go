package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

//创建一个用户的api

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

//用户的上线业务

func (this *User) Online() {
	//用户上线，将用户加入到Onlinemap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//广播当前用户上线消息
	this.server.BroadCast(this, "已经上线")
}

//用户的下线业务

func (this *User) Offline() {
	//用户下线，将用户从OnlineMap剔除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "下线")
}

//用户处理消息的业务

func (this *User) DoMessage(msg string) {
	//this.server.BroadCast(this, msg)
	//查询当前在线的用户有哪些
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			this.SendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式：rename|张三
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("当前用户名被使用\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMessage("您已经更新用户名：" + this.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		//1.获取对方的用户名
		if remoteName == "" {
			this.SendMessage("消息格式不正确，请使用\"to |csj|hello\"格式\n")
			return
		}
		//根据用户名，得到对方user对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMessage("该用户名不存在\n")
			return
		}
		//获取消息内容，通过对方的user对象将消息内容发送过去
		content := strings.Split(msg, "|")[1]
		if content == "" {
			this.SendMessage("无消息内容，请重发")
			return
		}
		remoteUser.SendMessage(this.Name + "对您说：" + content)

	} else {
		this.server.BroadCast(this, msg)
	}

}

//监听当前user channel的方法，一旦有消息，就直接发送给对端客户端

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

//给当前user对应的客户端发送消息

func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}
