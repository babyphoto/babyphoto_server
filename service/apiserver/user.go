package apiserver

import (
	"net/http"

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

	userinfo := model.UserInfo{
		UserCode:     UserCode,
		UserType:     UserType,
		UserNickName: UserNickName,
		UserName:     UserName,
		UserRegDtm:   UserRegDtm,
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
			response["result"] = "is Exist NickName"
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
				response["result"] = "Insert Failed"
			}
		}
	}
	res := map[string]interface{}{}
	res["result"] = response
	return c.JSON(http.StatusOK, res)
}

func (s *APIServer) UserList(c echo.Context) error {
	userinfos, err := s.db.AllUserList()
	util.CheckError("UserList ::: ", err)
	res := map[string]interface{}{}
	response := map[string]interface{}{}
	response["userList"] = userinfos
	res["result"] = response
	return c.JSON(http.StatusOK, res)
}
