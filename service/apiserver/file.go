package apiserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/babyphoto/babyphoto_server/service/util"
	"github.com/labstack/echo"
)

func MakeFolder(name string) error {
	err := os.Mkdir(`G:\공유 드라이브\babyphoto\images\`+name, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (s *APIServer) FileList(c echo.Context) error {
	GroupNum := c.FormValue("groupNum")
	if len(GroupNum) == 0 {
		return c.JSON(http.StatusBadRequest, "groupNum가 없습니다.")
	}
	groupNum, err := strconv.Atoi(GroupNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "groupNum형식이 잘못되었습니다")
	}

	FileList, err := s.db.FileList(groupNum)
	util.CheckError("Group_MyGroupList ::: ", err)
	response := map[string]interface{}{}
	response["fileList"] = FileList
	res := util.ReturnMap(response)
	return c.JSON(http.StatusOK, res)
}

func (s *APIServer) UploadFiles(c echo.Context) error {
	UserNum := c.FormValue("userNum")
	if len(UserNum) == 0 {
		return c.JSON(http.StatusBadRequest, "userNum가 없습니다.")
	}
	userNum, err := strconv.Atoi(UserNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "userNum형식이 잘못되었습니다")
	}

	GroupNum := c.FormValue("groupNum")
	if len(UserNum) == 0 {
		return c.JSON(http.StatusBadRequest, "groupNum가 없습니다.")
	}
	groupNum, err := strconv.Atoi(GroupNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "groupNum형식이 잘못되었습니다")
	}

	form, err := c.MultipartForm()
	util.CheckError("file_MultipartForm ::: ", err)
	files := form.File["files"]
	fmt.Println(files)
	isSuccess := false
	fmt.Println(isSuccess)
	UserInfo, err := s.db.GetUserWithUserNum(userNum)
	util.CheckError("file_InsertFile ::: ", err)

	FilePath := `G:\공유 드라이브\babyphoto\images\` + UserInfo.UserType + `.` + UserInfo.UserCode + `\`

	for _, file := range files {
		src, err := file.Open()
		util.CheckError("file_file.Open() ::: ", err)
		defer src.Close()

		fmt.Println(file.Filename)

		dst, err := os.Create(FilePath + file.Filename)
		util.CheckError("file_os.Create ::: ", err)
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			util.CheckError("file_io.Copy ::: ", err)
		}

		ss := strings.Split(file.Filename, ".")
		FileExtention := ss[len(ss)-1]

		result, err := s.db.InsertFile(userNum, groupNum, file.Filename, FilePath+file.Filename, FileExtention)
		util.CheckError("file_os.Create ::: ", err)
		if result < 0 {
			return c.JSON(http.StatusOK, "Upload fail")
		}
	}

	res := map[string]interface{}{}
	res["result"] = "file upload success"
	return c.JSON(http.StatusOK, res)
}

func (s *APIServer) UploadFile(c echo.Context) error {
	log.Println("upload")

	UserNum := c.FormValue("userNum")
	if len(UserNum) == 0 {
		return c.JSON(http.StatusBadRequest, "userNum가 없습니다.")
	}
	userNum, err := strconv.Atoi(UserNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "userNum형식이 잘못되었습니다")
	}

	fileName := c.FormValue("fileName")
	if len(fileName) == 0 {
		return c.JSON(http.StatusBadRequest, "fileName가 없습니다.")
	}

	form, err := c.MultipartForm()
	util.CheckError("file_MultipartForm ::: ", err)
	files := form.File["files"]
	fmt.Println(files)
	isSuccess := false
	fmt.Println(isSuccess)
	UserInfo, err := s.db.GetUserWithUserNum(userNum)
	util.CheckError("file_InsertFile ::: ", err)

	FilePath := `G:\공유 드라이브\babyphoto\images\` + UserInfo.UserType + `.` + UserInfo.UserCode + `\`

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(FilePath + fileName)
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

func (s *APIServer) DownloadFile(c echo.Context) error {
	path := c.QueryParam("path")
	return c.File(path)
}
