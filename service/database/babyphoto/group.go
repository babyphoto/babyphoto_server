package babyphoto

import (
	"github.com/babyphoto/babyphoto_server/service/util"
	_ "github.com/go-sql-driver/mysql"
)

func (db *BabyPhotoDB) NextSerialNum() (int, error) {
	count := 0
	err := db.DB.QueryRow(`
		SELECT COUNT(*) as count FROM GroupInfo
	`).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *BabyPhotoDB) CreateGroup(UserNum int, GroupName string) (int64, error) {
	NextID, err := db.NextSerialNum()
	if err != nil {
		return -1, err
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO GroupInfo (GroupName, GroupRegDtm) VALUES (?, ?)
	`, GroupName, util.CurrentDateTime())
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec(`
		INSERT INTO GroupUserInfo (UserNum, GroupNum, IsAdmin, AbleUpload, AbleDelete, AbleView, GUJoinDtm, GUUpdateDtm) VALUES (?, ?, 'Y', 'Y', 'Y', 'Y', ?, ?) 
	`, UserNum, NextID, util.CurrentDateTime(), util.CurrentDateTime())
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return 1, nil
}