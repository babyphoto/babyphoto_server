package apiserver

import (
	"log"
	"net/http"
	"strconv"

	"github.com/babyphoto/babyphoto_server/service/model"
	"github.com/babyphoto/babyphoto_server/service/util"
	"github.com/labstack/echo"
)

func (s *APIServer) RegistUser(c echo.Context) error {
	UserCode := c.FormValue("userCode")
	if len(UserCode) == 0 {
		return c.JSON(http.StatusBadRequest, "UserCode가 없습니다.")
	}
	UserType := c.FormValue("userType")
	if len(UserType) == 0 {
		return c.JSON(http.StatusBadRequest, "UserType이 없습니다.")
	}
	UserNickName := c.FormValue("userNickName")
	UserName := c.FormValue("userName")
	UserRegDtm := util.CurrentDateTime()
	UserProfile := c.FormValue("userProfile")

	log.Println(UserCode, UserType, UserNickName, UserName, UserRegDtm, UserProfile)

	userinfo := model.UserInfo{
		UserCode:     UserCode,
		UserType:     UserType,
		UserNickName: UserNickName,
		UserName:     UserName,
		UserRegDtm:   UserRegDtm,
		UserProfile:  UserProfile,
	}
	response := map[string]interface{}{}
	count, err := s.db.IsExistUser(userinfo.UserType, userinfo.UserCode)
	util.CheckError("isExistUser ::: ", err)
	if count > 0 {
		userinfo, err := s.db.GetUser(userinfo)
		util.CheckError("Regist_GetUser :::", err)
		response["userinfo"] = userinfo
	} else {
		count, err := s.db.IsExistNickName(userinfo.UserNickName)
		util.CheckError("isExistNickName ::: ", err)
		if count > 0 {
			return c.JSON(http.StatusOK, "is Exist NickName")
		} else {
			result, err := s.db.InsertUser(userinfo)
			util.CheckError("Regist_InsertUser :::", err)
			if result > 0 {
				folderName := UserType + "." + UserCode
				err = MakeFolder(folderName)
				util.CheckError("MKDir ::::: ", err)
				userinfo, err := s.db.GetUser(userinfo)
				util.CheckError("Regist_GetUser :::", err)
				response["userinfo"] = userinfo
			} else {
				return c.JSON(http.StatusOK, "Insert Failed")
			}
		}
	}
	res := util.ReturnMap(response)
	return c.JSON(http.StatusOK, res)
}

func (s *APIServer) UpdateUserNickName(c echo.Context) error {
	UserNum := c.FormValue("UserNum")
	if len(UserNum) == 0 {
		return c.JSON(http.StatusBadRequest, "UserNum가 없습니다.")
	}
	userNum, err := strconv.Atoi(UserNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "UserNum형식이 잘못되었습니다")
	}
	UserNickName := c.FormValue("UserNickName")
	if len(UserNickName) == 0 {
		return c.JSON(http.StatusBadRequest, "UserNickName이 없습니다.")
	}
	count, err := s.db.IsExistNickName(UserNickName)
	util.CheckError("isExistNickName ::: ", err)
	if count > 0 {
		return c.JSON(http.StatusOK, "is Exist NickName")
	} else {
		result, err := s.db.UpdateUserNickName(UserNickName, userNum)
		util.CheckError("UpdateUserNickName_UpdateUserNickName :::", err)
		if result > 0 {
			return c.JSON(http.StatusOK, "Update Success")
		} else {
			return c.JSON(http.StatusOK, "Update Failed")
		}
	}
}

func (s *APIServer) UserSearchWithNickName(c echo.Context) error {
	UserNickName := c.FormValue("userNickName")
	if len(UserNickName) == 0 {
		UserNickName = ""
	}
	userinfos, err := s.db.SearchUserList(UserNickName)
	util.CheckError("UserSearchWithNickName.SearchUserList :::", err)
	response := map[string]interface{}{}
	response["userList"] = userinfos
	res := util.ReturnMap(response)
	return c.JSON(http.StatusOK, res)
}

func (s *APIServer) UserList(c echo.Context) error {
	userinfos, err := s.db.AllUserList()
	util.CheckError("UserList ::: ", err)
	response := map[string]interface{}{}
	response["userList"] = userinfos
	res := util.ReturnMap(response)
	return c.JSON(http.StatusOK, res)
}
