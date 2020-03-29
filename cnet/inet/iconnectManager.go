package inet

import "net"

type IConnectManager interface {
	Iserver
	Add(*net.TCPConn) string
	StopConnect(string)
}
