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
	err = os.Mkdir(`G:\공유 드라이브\babyphoto\images\`+name+`\thumbnail`, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (s *APIServer) FileList(c echo.Context) error {
	UserNum := c.FormValue("userNum")
	if len(UserNum) == 0 {
		return c.JSON(http.StatusBadRequest, "userNum가 없습니다.")
	}
	userNum, err := strconv.Atoi(UserNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "userNum형식이 잘못되었습니다")
	}

	GroupNum := c.FormValue("groupNum")
	if len(GroupNum) == 0 {
		return c.JSON(http.StatusBadRequest, "groupNum가 없습니다.")
	}
	groupNum, err := strconv.Atoi(GroupNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "groupNum형식이 잘못되었습니다")
	}

	FileList, err := s.db.FileList(userNum, groupNum)
	util.CheckError("Group_MyGroupList ::: ", err)
	if FileList == nil && err == nil {
		return c.JSON(http.StatusOK, "Viewing fail - Lack of authority")
	} else {
		response := map[string]interface{}{}
		response["fileList"] = FileList
		res := util.ReturnMap(response)
		return c.JSON(http.StatusOK, res)
	}
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
	FileThumbnail := `G:\공유 드라이브\babyphoto\images\` + UserInfo.UserType + `.` + UserInfo.UserCode + `\thumbnail\`
	for _, file := range files {
		src, err := file.Open()
		util.CheckError("file_file.Open() ::: ", err)
		defer src.Close()

		log.Println(file.Filename)

		dst, err := os.Create(FilePath + file.Filename)
		util.CheckError("file_os.Create_Create ::: ", err)
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			util.CheckError("file_io.Copy ::: ", err)
		}

		ss := strings.Split(file.Filename, ".")
		FileExtention := ss[len(ss)-1]

		result, err := s.db.InsertFile(userNum, groupNum, file.Filename, FilePath+file.Filename, FileThumbnail+file.Filename, FileExtention)
		util.CheckError("file_os.Create_InsertFile ::: ", err)
		if result < 0 {
			return c.JSON(http.StatusOK, "Upload fail")
		}
	}

	res := map[string]interface{}{}
	res["result"] = "file upload success"
	return c.JSON(http.StatusOK, res)
}

func (s *APIServer) UploadFile(c echo.Context) error {
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

	fileType := c.FormValue(("fileType"))
	if len(fileType) == 0 {
		fileType = "image"
	}

	UserInfo, err := s.db.GetUserWithUserNum(userNum)
	util.CheckError("file_InsertFile ::: ", err)

	FilePath := `G:\공유 드라이브\babyphoto\images\` + UserInfo.UserType + `.` + UserInfo.UserCode + `\thumbnail\`
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	InsertFileName := ""
	if fileType == "video" {
		InsertFileName = fileName
	} else {
		InsertFileName = file.Filename
	}

	log.Println(InsertFileName)

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(FilePath + InsertFileName)
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
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	// Type, err := util.GetFileContentType(f)
	// if err != nil {
	// 	return err
	// }

	return c.Inline(path, f.Name())
}

func (s *APIServer) DownloadVideo(c echo.Context) error {
	path := c.QueryParam("path")
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	Type, err := util.GetFileContentType(f)
	if err != nil {
		return err
	}

	return c.Stream(http.StatusOK, Type, f)
}

func (s *APIServer) DeleteFile(c echo.Context) error {

	UserNum := c.FormValue("userNum")
	if len(UserNum) == 0 {
		return c.JSON(http.StatusBadRequest, "userNum가 없습니다.")
	}
	userNum, err := strconv.Atoi(UserNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "userNum형식이 잘못되었습니다")
	}

	GroupNum := c.FormValue("groupNum")
	if len(GroupNum) == 0 {
		return c.JSON(http.StatusBadRequest, "groupNum가 없습니다.")
	}
	groupNum, err := strconv.Atoi(GroupNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "groupNum형식이 잘못되었습니다")
	}

	FileNum := c.FormValue("fileNum")
	if len(FileNum) == 0 {
		return c.JSON(http.StatusBadRequest, "fileNum가 없습니다.")
	}
	fileNum, err := strconv.Atoi(FileNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "fileNum형식이 잘못되었습니다")
	}

	result, err := s.db.UpdateFile(userNum, fileNum, groupNum)
	util.CheckError("Group_InviteGroup ::: ", err)

	if result < 0 {
		return c.JSON(http.StatusOK, "Delete fail")
	} else if result == 0 {
		return c.JSON(http.StatusOK, "Delete fail - Lack of authority")
	}

	return c.JSON(http.StatusOK, "Delete Success")
}
