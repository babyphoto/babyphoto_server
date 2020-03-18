package main

import (
	"github.com/babyphoto/babyphoto_server/service/apiserver"
	"github.com/babyphoto/babyphoto_server/service/database/mariadb"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := mariadb.Connect()

	API := apiserver.NewAPIServer(db)
	go API.Run(":8080")

	select {}
}
