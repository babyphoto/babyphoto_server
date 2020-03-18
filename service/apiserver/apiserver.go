package apiserver

import (
	"database/sql"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type APIServer struct {
	e  *echo.Echo
	db *sql.DB
}

func NewAPIServer(db *sql.DB) *APIServer {
	s := &APIServer{
		e:  echo.New(),
		db: db,
	}
	return s
}

func (s *APIServer) Run(BindAddress string) error {
	e := s.e
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	api := e.Group("/api")
	user := api.Group("/user")
	user.GET("/userList", s.UserList)
	file := api.Group("/files")
	file.GET("/", s.Connect)
	file.POST("/upload", s.uploadFile)
	file.GET("/download", s.downloadFile)

	return s.e.Start(BindAddress)
}
