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
	api.GET("/connect", s.Connect)

	user := api.Group("/user")
	user.GET("/userList", s.UserList)
	user.GET("/userSearchWithNickName", s.UserSearchWithNickName)
	user.POST("/updateNickName", s.UpdateUserNickName)
	user.POST("/regist", s.RegistUser)
	user.GET("/groupUserList", s.GroupUserList)

	group := api.Group("/group")
	group.POST("/createGroup", s.CreateGroup)
	group.POST("/deleteGroup", s.DeleteGroup)
	group.POST("/updateGroup", s.UpdateGroup)
	group.POST("/leaveGroup", s.LeaveGroup)
	group.GET("/groupList", s.GroupList)
	group.POST("/inviteGroup", s.InviteGroup)

	file := api.Group("/files")
	file.GET("/download", s.DownloadFile)
	file.POST("/upload", s.UploadFiles)
	file.GET("/fileList", s.FileList)
	file.POST("/thumnail", s.UploadFile)
	file.POST("/delete", s.DeleteFile)

	return s.e.Start(BindAddress)
}

func (s *APIServer) Connect(c echo.Context) error {
	return c.JSON(http.StatusOK, "Connect Success")
}
