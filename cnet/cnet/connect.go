package cnet

import (
	"fmt"
	"net"
	"time"
)

type Connect struct {
	conn *net.TCPConn
	attributes map[string]interface{}

	stopT *time.Timer

	msgChannel chan []byte
	exitChannel chan interface{}

}

func (connect *Connect) AddAttribute(k string, v interface{})  {
	connect.attributes[k] = v
}

func (connect Connect) GetAttribute(k string) interface{}  {
	if _, ok := connect.attributes[k]; ok {
		return connect.attributes[k]
	}
	return nil
}

func (connect *Connect) startWrite()  {
	defer func() {
		fmt.Println("connet write close")
	}()

	for{
		fmt.Println("write")


		select {
		case data, ok := <-connect.msgChannel:
			if ok {
				_ ,err := connect.conn.Write(data)
				if err != nil {

				}
			} else {
				return
			}
		}
	}
}

func (connect *Connect) SendMsg(msgId uint32, data []byte)  {

	pack := NewPack()

	connect.msgChannel <- pack.Pack(NewMassage(msgId, data))
}

func (connect *Connect) startRead()  {
	defer func() {
		fmt.Println("connet reader close")
	}()
	work:= WorkManagerInstance()
	pack := NewPack()
	for{

		fmt.Println(" reader ")
		select {
		case _, ok := <-connect.exitChannel:
			fmt.Println("read ", ok)
			return
		default:
			_, err := pack.UnPack(connect.conn)
			if err != nil {
				continue
			}
			msg := &Massage{}
			request := NewRequest(*connect, *msg)
			work.SendRequest(request)
		}
	}
}

func NewConnect(conn *net.TCPConn) *Connect {
	return &Connect{
		conn: conn,
		stopT: time.NewTimer(10 * time.Second),
		attributes: make(map[string]interface{}),
		msgChannel: make(chan []byte),
		exitChannel: make(chan interface{}),
	}
}

func (connect *Connect) Stop() error {

	close(connect.exitChannel)
	close(connect.msgChannel)
	connect.stopT.Stop()

	return connect.conn.Close()
}



func (connect *Connect) Start()  {
	defer func() {
		fmt.Println("connect close")

		// todo 发送消息给 connectManager 来关闭connect相关进程
		_ = connect.Stop()
	}()

	go connect.startWrite()

	go connect.startRead()

	fmt.Println("connect start")

	<- connect.stopT.C

	fmt.Println((*connect).GetAttribute("connId") , " close")

}
