package cnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

type Massage struct {
	MsgId uint32
	DataLength uint32
	Data []byte
}

func NewMassage(msgId uint32, data []byte) Massage {
	return Massage{
		MsgId:      msgId,
		DataLength: uint32(len(data)),
		Data:       data,
	}
}



func (massage Massage) GetData() []byte {

	return massage.Data
}

func (massage Massage) GetMsgId() uint32  {

	return massage.MsgId
}

func (massage Massage) GetDataLength() uint32  {

	return massage.DataLength
}

type Pack struct { }

func NewPack() *Pack {
	return &Pack{}
}

func (pack *Pack) getHeadLength() int  {
	return 8
}

func (pack *Pack) Pack(msg Massage) []byte  {
	dataBuff := bytes.NewBuffer([]byte{})
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err!= nil {
		fmt.Println("write msgId error ", err)
		return nil
	}
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLength())
	if err!= nil {
		fmt.Println("write msgLength error ", err)
		return nil
	}
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err!= nil {
		fmt.Println("write msgData error ", err)
		return nil
	}

	return dataBuff.Bytes()

}

func (pack *Pack) UnPack(tcpConn *net.TCPConn) (*Massage, error) {
	headData := make([]byte, pack.getHeadLength())

	_ = tcpConn.SetReadDeadline(time.Now().Add(3*time.Second))
	_, err := io.ReadFull(tcpConn, headData)
	if err != nil {
		fmt.Println("read msg head error ", err)
		return nil, err
	}
	_ = tcpConn.SetReadDeadline(time.Time{})


	dataBuff := bytes.NewReader(headData)
	msg := &Massage{}
	err = binary.Read(dataBuff, binary.LittleEndian, &msg.MsgId)
	if err!= nil {
		fmt.Println("read msgId error ", err)
		return nil, err
	}

	err = binary.Read(dataBuff, binary.LittleEndian, &msg.DataLength)
	if err!= nil {
		fmt.Println("read msgLength error ", err)
		return nil, err
	}

	msg.Data = make([]byte, msg.GetDataLength())
	_, err = io.ReadFull(tcpConn, msg.Data)
	if err != nil {
		fmt.Println("read msgLength error ", err)
		return nil, err
	}

	return msg, nil
}
