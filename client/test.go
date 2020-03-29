package main

import (
	"cnet/cnet/cnet"
)

func main()  {

	server := cnet.ServerInstance()

	server.Server()

}