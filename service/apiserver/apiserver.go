package apiserver

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type APIServer struct {
	e      *echo.Echo
	dbpath string
}

func NewAPIServer(dbpath string) *APIServer {
	s := &APIServer{
		e:      echo.New(),
		dbpath: dbpath,
	}
	return s
}

func (s *APIServer) Run(BindAddress string) error {
	e := s.e
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	api := e.Group("/api")
	file := api.Group("/files")
	file.GET("/", s.Connect)
	file.POST("/upload", s.uploadFile)
	file.GET("/download", s.downloadFile)
	return s.e.Start(BindAddress)
}
