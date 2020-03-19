package apiserver

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (s *APIServer) Run(BindAddress string) error {
	e := s.e
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	api := e.Group("/api")
	user := api.Group("/user")
	user.GET("/regist", s.RegistUser)
	user.GET("/userList", s.UserList)
	file := api.Group("/files")
	file.GET("/", s.Connect)
	file.POST("/upload", s.UploadFile)
	file.GET("/download", s.DownloadFile)

	return s.e.Start(BindAddress)
}

func (s *APIServer) Connect(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
