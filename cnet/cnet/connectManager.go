package cnet

import (
	"cnet/cnet/conf"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type connectManager struct {
	connectPool map[string]*Connect
}

var	managerInstance *connectManager

func (connectManager *connectManager) Add(connect *net.TCPConn) (connId string)  {
	if len(connectManager.connectPool) <= conf.GlobalConfig.MaxConnectNum {

		connId = connectManager.createConnId()

		connectManager.connectPool[connId] = NewConnect(connect)
		connectManager.connectPool[connId].AddAttribute("connId", connId)
		go connectManager.connectPool[connId].Start()
	}

	return
}

func (connectManager *connectManager) SendMessage(connId string, msgId uint32,data []byte)  {
	if conn, ok := connectManager.connectPool[connId]; ok {
		conn.SendMsg(msgId, data)
	}
}

func (connectManager *connectManager) createConnId() (connId string)  {
	timeStamp := time.Now().Unix()
	rand.Seed(timeStamp)
	randInt64, _ := strconv.ParseInt(strconv.Itoa(rand.Intn(9000)+1000), 10, 64)
	connId = strconv.FormatInt(timeStamp*100000 + randInt64, 32)
	return
}

func (connectManager *connectManager) StopConnect(connId string)  {

	if _, ok := connectManager.connectPool[connId]; ok {
		_ = connectManager.connectPool[connId].Stop()
		delete(connectManager.connectPool, connId)
	}

	return
}
func (connectManager *connectManager) Stop() {
	for connId, _ := range connectManager.connectPool {
		connectManager.StopConnect(connId)
	}
}

func ConnectManagerInstance() *connectManager  {

	if managerInstance == nil {
		managerInstance = &connectManager{
			connectPool:make(map[string]*Connect),
		}
	}

	return managerInstance
}


