package main

import (
	"database/sql"

	"github.com/babyphoto/babyphoto_server/service/apiserver"
	"github.com/babyphoto/babyphoto_server/service/database/babyphoto"
	"github.com/babyphoto/babyphoto_server/service/util"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "sherwher:lskun@tcp(174.129.71.0:3306)/babyphoto")
	util.CheckError("main ::: db connectoin ::: ", err)
	defer db.Close()

	bybyPhotodb := &babyphoto.BabyPhotoDB{
		DB: db,
	}

	API := apiserver.NewAPIServer(bybyPhotodb)
	go API.Run(":8080")
	select {}
}
