package cnet

import (
	"cnet/cnet/conf"
	"errors"
	"fmt"
	"time"
)

type workManager struct {
	taskWordPool  []chan Request
	router map[uint32]func(request Request) []byte
}

var	workInstance *workManager

func WorkManagerInstance() *workManager  {
	if workInstance == nil {
		workInstance = &workManager{
			taskWordPool:make([]chan Request, conf.GlobalConfig.TaskWorkNum ),
		}
	}
	return workInstance
}

func (workManager *workManager) startOneTask(i int, taskWork chan Request)  {
	fmt.Println("task work ", i, " starting...")
	ConnectManager := ConnectManagerInstance()
	for {
		select {
		case request, ok := <-taskWork:
			if !ok {
				fmt.Println("task work ", i, " quit...")
				return
			}
			if function ,ok:= workManager.router[request.GetMsg().GetMsgId()]; !ok {
				//发送消息 给 路由
				data := function(request)
				ConnectManager.SendMessage(request.GetConn().GetAttribute("connId").(string), request.GetMsg().GetMsgId(), data)
			}
		}
	}
}

func (workManager *workManager) AddRouter(routerId uint32, method func(request Request) []byte) error  {
	if _, ok := workManager.router[routerId]; ok {
		return errors.New("routerId is exist")
	}

	workManager.router[routerId] = method
	return nil
}

func (workManager *workManager) SendRequest(request Request)  {
	i := time.Now().Unix() % int64(conf.GlobalConfig.TaskWorkNum)

	workManager.taskWordPool[i] <- request
}

func (workManager *workManager) startAllTask()  {

	for i:=0; i < conf.GlobalConfig.TaskWorkNum; i++ {
		workManager.taskWordPool[i] = make(chan Request, conf.GlobalConfig.WorkBuffer)

		go workManager.startOneTask(i, workManager.taskWordPool[i])
	}
}

func (workManager *workManager) Stop()  {
	for i, _ := range workManager.taskWordPool{
		close(workManager.taskWordPool[i])
	}

	return
}

func (workManager *workManager) Start()  {
	workManager.startAllTask()

}

