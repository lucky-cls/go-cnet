package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
)


type config struct {
	//serverConfig
	Name string `json:"name"`
	NetWork string `json:"network"`
	Ip string `json:"ip"`
	Port int`json:"port"`
	AcceptTimeOut int `json:"accept_time_out"`


	//connectConfig
	MaxConnectNum int `json:"max_connect_num"`


	//workManagerConfig
	TaskWorkNum int `json:"task_work_num"`
	WorkBuffer  int  `json:"work_buffer"`

	ConfigPath string
}

var GlobalConfig *config

func pathExist(path string) (bool, error)  {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}

	return false, err
}



func (c *config) Reload()  {

	if isExist, err := pathExist(c.ConfigPath); isExist == false {
		panic(err)
	}

	configString, err := ioutil.ReadFile(c.ConfigPath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configString, c)
	if err != nil {
		panic(err.Error())
	}

}

func init()  {

	GlobalConfig = &config{
		Name: "CNAT",
		NetWork: "tcp4",
		Ip: "0.0.0.0",
		Port: 9005,
		ConfigPath: "conf/config.json",
	}

	GlobalConfig.Reload()
}
