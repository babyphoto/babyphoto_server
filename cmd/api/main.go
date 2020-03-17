package main

import (
	"github.com/babyphoto/babyphoto_server/service/apiserver"
)

func main() {

	API := apiserver.NewAPIServer("/")
	go API.Run(":8080")

	select {}
}
