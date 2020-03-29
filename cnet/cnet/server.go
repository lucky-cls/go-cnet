package cnet

import (
	"cnet/cnet/conf"
	"fmt"
	"net"
	"time"
)

type Server struct {
	Name string
	IpVersion string
	Ip string
	Port int
	isRunning bool
	stopChannel chan interface{}
	connChannel chan *net.Conn
}

var instance *Server

func (s *Server) loadConfig() {
	s.Name = conf.GlobalConfig.Name
	s.IpVersion = conf.GlobalConfig.NetWork
	s.Ip = conf.GlobalConfig.Ip
	s.Port = conf.GlobalConfig.Port
}

func (s *Server) start() {

	fmt.Println("starting Tcp server...")

	tcpAddr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		panic(err)
	}

	tcpListener, err := net.ListenTCP(s.IpVersion, tcpAddr)
	if err != nil {
		panic(err)
	}



	connectManager := ConnectManagerInstance()
	workManager := WorkManagerInstance()
	workManager.Start()

	fmt.Println("starting Tcp Listening... ")
	for {
		select {
		case _, ok := <-s.stopChannel:
			if !ok {
				_ = tcpListener.Close()
				connectManager.Stop()
				workManager.Stop()
				return
			}
		default:
			_ = tcpListener.SetDeadline(time.Now().Add(time.Duration(conf.GlobalConfig.AcceptTimeOut) * time.Second))
			conn, err := tcpListener.AcceptTCP()
			if err != nil {
				fmt.Println("accept one tcp connect error..", err)
			} else  {
				fmt.Println(conn.RemoteAddr().String(), "accept one tcp connect success..")

				connId := connectManager.Add(conn)
				if connId != "" {
					fmt.Println("connId: ", connId, " add connect success..")
				} else {
					fmt.Println("connId: ", connId, " add connect faild..")
				}
			}
		}
	}
}

func (s *Server) Stop()  {
	if s.isRunning == false {
		return
	}

	fmt.Println("send stop singnal...")
	close(s.stopChannel)
	return
}

func ServerInstance() *Server  {
	if instance == nil {
		instance = &Server{
			isRunning:false,
			stopChannel:make(chan interface{}),
		}
	}

	return instance
}

func (s *Server) Server()  {

	if s.isRunning {
		return
	}

	s.loadConfig()

	go s.start()

	s.isRunning = true
	fmt.Println("starting Tcp server complete ...")

	select {
	case _, ok := <-s.stopChannel:
		if !ok {
			fmt.Println("recive stop singnal...")
			s.isRunning = false
		}
		break
	}
}