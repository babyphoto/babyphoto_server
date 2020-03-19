package apiserver

import (
	"github.com/babyphoto/babyphoto_server/service/database/babyphoto"
	"github.com/labstack/echo"
)

type APIServer struct {
	e  *echo.Echo
	db *babyphoto.BabyPhotoDB
}

func NewAPIServer(DB *babyphoto.BabyPhotoDB) *APIServer {
	s := &APIServer{
		e:  echo.New(),
		db: DB,
	}
	return s
}
