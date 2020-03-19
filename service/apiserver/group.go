package apiserver

import (
	"github.com/babyphoto/babyphoto_server/service/model"
)

func (s *APIServer) GroupList(UserNum int) []model.GroupInfo {
	// db := mariadb.Connect()
	// defer db.Close()
	// rows, err := db.Query(`
	// 	INSERT INTO UserInfo(UserCode, UserType) VALUES (?, ?)
	// `, UserNum)
	// util.CheckError("GroupList : ", err)
	return []model.GroupInfo{}
}
