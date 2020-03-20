package apiserver

import (
	"net/http"
	"strconv"

	"github.com/babyphoto/babyphoto_server/service/util"
	"github.com/labstack/echo"
)

func (s *APIServer) CreateGroup(c echo.Context) error {
	UserNum := c.FormValue("userNum")
	if len(UserNum) == 0 {
		UserNum = ""
	}
	userNum, err := strconv.Atoi(UserNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "UserNum형식이 잘못되었습니다")
	}
	GroupName := c.FormValue("groupName")
	if len(GroupName) == 0 {
		GroupName = ""
	}

	result, err := s.db.CreateGroup(userNum, GroupName)
	util.CheckError("Group_CreateGroup ::: ", err)
	if result < 0 {
		c.JSON(http.StatusOK, "Create fail")
	}

	return c.JSON(http.StatusOK, "Create Success")
}
