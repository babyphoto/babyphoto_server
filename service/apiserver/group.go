package apiserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/babyphoto/babyphoto_server/service/util"
	"github.com/labstack/echo"
)

func (s *APIServer) CreateGroup(c echo.Context) error {
	UserNum := c.FormValue("userNum")
	if len(UserNum) == 0 {
		return c.JSON(http.StatusBadRequest, "userNum가 없습니다.")
	}
	userNum, err := strconv.Atoi(UserNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "userNum형식이 잘못되었습니다")
	}
	GroupName := c.FormValue("groupName")
	if len(GroupName) == 0 {
		GroupName = ""
	}

	result, err := s.db.CreateGroup(userNum, GroupName)
	util.CheckError("Group_CreateGroup ::: ", err)
	if result < 0 {
		return c.JSON(http.StatusOK, "Create fail")
	} else if result == 0 {
		return c.JSON(http.StatusOK, "Create fail - Group Name Exist")
	}
	return c.JSON(http.StatusOK, "Create Success")
}

func (s *APIServer) UpdateGroup(c echo.Context) error {
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

	GroupName := c.FormValue("groupName")
	if len(GroupName) == 0 {
		GroupName = ""
	}

	result, err := s.db.UpdateGroup(userNum, groupNum, GroupName)
	util.CheckError("Group_UpdateGroup ::: ", err)
	if result < 0 {
		return c.JSON(http.StatusOK, "Update fail")
	} else if result == 0 {
		return c.JSON(http.StatusOK, "Update fail - Lack of authority")
	}

	return c.JSON(http.StatusOK, "Update Success")
}

func (s *APIServer) DeleteGroup(c echo.Context) error {
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

	result, err := s.db.DeleteGroup(userNum, groupNum)
	util.CheckError("Group_DeleteGroup ::: ", err)
	if result < 0 {
		return c.JSON(http.StatusOK, "Delete fail")
	} else if result == 0 {
		return c.JSON(http.StatusOK, "Delete fail - Lack of authority")
	}

	return c.JSON(http.StatusOK, "Delete Success")
}

func (s *APIServer) LeaveGroup(c echo.Context) error {
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

	result, err := s.db.LeaveGroup(userNum, groupNum)
	util.CheckError("Group_LeaveGroup ::: ", err)
	if result < 0 {
		return c.JSON(http.StatusOK, "Leave fail")
	} else if result == 0 {
		return c.JSON(http.StatusOK, "Leave fail - Lack of authority")
	}

	return c.JSON(http.StatusOK, "Leave Success")
}

func (s *APIServer) GroupList(c echo.Context) error {
	UserNum := c.FormValue("userNum")
	if len(UserNum) == 0 {
		return c.JSON(http.StatusBadRequest, "userNum가 없습니다.")
	}
	userNum, err := strconv.Atoi(UserNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "userNum형식이 잘못되었습니다")
	}

	fmt.Println(UserNum)

	MyGroupList, err := s.db.MyGroupList(userNum)
	util.CheckError("Group_MyGroupList ::: ", err)

	InvitedGroupList, err := s.db.InvitedGroupList(userNum)
	util.CheckError("Group_InvitedGroupList ::: ", err)

	response := map[string]interface{}{}
	response["myGroupList"] = MyGroupList
	response["invitedGroupList"] = InvitedGroupList
	res := util.ReturnMap(response)
	return c.JSON(http.StatusOK, res)
}

func (s *APIServer) InviteGroup(c echo.Context) error {
	// UserNum int, GroupNum int, InviteUserNum int, AbleUpload string, AbleDelete string, AbleView string
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

	InviteUserNum := c.FormValue("inviteUserNum")
	if len(InviteUserNum) == 0 {
		return c.JSON(http.StatusBadRequest, "inviteUserNum가 없습니다.")
	}
	inviteUserNum, err := strconv.Atoi(InviteUserNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "inviteUserNum형식이 잘못되었습니다")
	}

	AbleUpload := c.FormValue("ableUpload")
	if len(AbleUpload) == 0 {
		AbleUpload = "N"
	}

	AbleDelete := c.FormValue("ableDelete")
	if len(AbleDelete) == 0 {
		AbleDelete = "N"
	}

	AbleView := c.FormValue("ableView")
	if len(AbleView) == 0 {
		AbleView = "Y"
	}

	result, err := s.db.InviteGroup(userNum, groupNum, inviteUserNum, AbleUpload, AbleDelete, AbleView)
	util.CheckError("Group_InviteGroup ::: ", err)

	if result < 0 {
		return c.JSON(http.StatusOK, "Invite fail")
	} else if result == 0 {
		return c.JSON(http.StatusOK, "Invite fail - Lack of authority")
	}

	return c.JSON(http.StatusOK, "Invite Success")
}
