package babyphoto

import (
	"github.com/babyphoto/babyphoto_server/service/model"
)

func (db *BabyPhotoDB) IsExistUser(UserType string, UserCode string) (int, error) {
	count := 0
	err := db.DB.QueryRow(`
		SELECT COUNT(*) as count FROM UserInfo WHERE UserType=? AND UserCode=?
	`, UserType, UserCode).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *BabyPhotoDB) IsExistNickName(UserNickName string) (int, error) {
	count := 0
	err := db.DB.QueryRow(`
		SELECT COUNT(*) as count FROM UserInfo WHERE UserNickName=?
	`, UserNickName).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *BabyPhotoDB) GetUser(m model.UserInfo) (model.UserInfo, error) {
	rows, err := db.DB.Query("SELECT * FROM UserInfo WHERE UserType=? AND UserCode=?", m.UserType, m.UserCode)
	defer rows.Close()
	if err != nil {
		return model.UserInfo{}, err
	}
	userinfo := model.UserInfo{}
	for rows.Next() {
		rows.Scan(&userinfo.UserNum, &userinfo.UserCode, &userinfo.UserType, &userinfo.UserNickName, &userinfo.UserName, &userinfo.UserRegDtm)
	}
	return userinfo, nil
}

func (db *BabyPhotoDB) InsertUser(m model.UserInfo) (int64, error) {
	result, err := db.DB.Exec(`
		INSERT INTO UserInfo(UserCode, UserType, UserNickName, UserName, UserRegDtm) VALUES (?, ?, ?, ?, ?)
	`, m.UserCode, m.UserType, m.UserNickName, m.UserName, m.UserRegDtm)
	if err != nil {
		return -1, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return n, nil
}

func (db *BabyPhotoDB) AllUserList() ([]model.UserInfo, error) {
	rows, err := db.DB.Query("SELECT * from UserInfo")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	userinfos := []model.UserInfo{}
	for rows.Next() {
		userinfo := model.UserInfo{}
		rows.Scan(&userinfo.UserNum, &userinfo.UserCode, &userinfo.UserType, &userinfo.UserNickName, &userinfo.UserName, &userinfo.UserRegDtm)
		userinfos = append(userinfos, userinfo)
	}
	return userinfos, nil
}
