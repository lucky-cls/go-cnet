package cnet

type Request struct {
	conn Connect
	msg  Massage
}

func NewRequest(conn Connect, msg Massage) Request  {
	return Request{
		conn: conn,
		msg: msg,
	}
}

func (request Request) GetConn() Connect {

	return request.conn
}

func (request Request) GetMsg() Massage {

	return request.msg
}



