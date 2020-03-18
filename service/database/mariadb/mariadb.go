package mariadb

import (
	"database/sql"

	"github.com/babyphoto/babyphoto_server/service/util"
)

func Connect() *sql.DB {

	db, err := sql.Open("mysql", "sherwher:lskun@tcp(174.129.71.0:3306)/babyphoto")
	util.CheckError("main ::: db connectoin ::: ", err)

	return db
}
