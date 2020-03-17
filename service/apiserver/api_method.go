package apiserver

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func (s *APIServer) Connect(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (s *APIServer) uploadFile(c echo.Context) error {
	log.Println("upload")
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	log.Println("upload" + file.Filename)

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(`G:\공유 드라이브\babyphoto\images\sherwher\` + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	res := map[string]interface{}{}
	res["result"] = "file upload success"
	return c.JSON(http.StatusOK, res)
}

func (s *APIServer) downloadFile(c echo.Context) error {
	name := c.QueryParam("name")
	return c.File(`G:\공유 드라이브\babyphoto\images\sherwher\` + name)
}
