package babyphoto

import (
	"github.com/babyphoto/babyphoto_server/service/model"
	"github.com/babyphoto/babyphoto_server/service/util"
	_ "github.com/go-sql-driver/mysql"
)

func (db *BabyPhotoDB) GroupNextSerialNum() (int, error) {
	count := 0
	err := db.DB.QueryRow(`
		SELECT MAX(GroupNum) + 1 as count FROM GroupInfo
	`).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *BabyPhotoDB) FineGroupName(UserNum int, GroupName string) (int, error) {
	count := 0
	err := db.DB.QueryRow(`
		SELECT 
			COUNT(*) as count
		  FROM GroupInfo A, GroupUserInfo B
		 WHERE A.GroupNum = B.GroupNum
		   AND B.UserNum = ?
		   AND B.IsAdmin = 'Y'
		   AND A.GroupName = ?
	`, UserNum, GroupName).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *BabyPhotoDB) GroupUserInfo(UserNum int, GroupNum int) (model.GroupUserInfo, error) {
	rows, err := db.DB.Query("SELECT * FROM GroupUserInfo WHERE UserNum=? AND GroupNum=?", UserNum, GroupNum)
	defer rows.Close()
	if err != nil {
		return model.GroupUserInfo{}, err
	}
	groupUserInfo := model.GroupUserInfo{}
	for rows.Next() {
		rows.Scan(&groupUserInfo.UserNum, &groupUserInfo.GroupNum, &groupUserInfo.IsAdmin, &groupUserInfo.AbleUpload, &groupUserInfo.AbleDelete, &groupUserInfo.AbleView, &groupUserInfo.GUJoinDtm, &groupUserInfo.GUUpdateDtm)
	}
	return groupUserInfo, nil
}

func (db *BabyPhotoDB) CreateGroup(UserNum int, GroupName string) (int64, error) {
	NextID, err := db.GroupNextSerialNum()
	if err != nil {
		return -1, err
	}

	count, err := db.FineGroupName(UserNum, GroupName)
	if err != nil {
		return -1, err
	}
	if count > 0 {
		return 0, nil
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO GroupInfo (GroupNum, GroupName, GroupRegDtm) VALUES (?, ?, ?)
	`, NextID, GroupName, util.CurrentDateTime())
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

func (db *BabyPhotoDB) UpdateGroup(UserNum int, GroupNum int, GroupName string) (int64, error) {
	GI, err := db.GroupUserInfo(UserNum, GroupNum)
	if err != nil {
		return -1, err
	}

	if GI.AbleUpload != "Y" {
		return 0, nil
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE GroupInfo SET GroupName=? WHERE GroupNum=?
	`, GroupName, GroupNum)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return 1, nil
}

func (db *BabyPhotoDB) DeleteGroup(UserNum int, GroupNum int) (int64, error) {
	GI, err := db.GroupUserInfo(UserNum, GroupNum)
	if err != nil {
		return -1, err
	}
	if GI.AbleDelete != "Y" {
		return 0, nil
	}
	tx, err := db.DB.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		DELETE FROM GroupUserInfo WHERE GroupNum=?
	`, GroupNum)
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec(`
		DELETE FROM GroupFileInfo WHERE GroupNum=?
	`, GroupNum)
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec(`
		DELETE FROM GroupInfo WHERE GroupNum=?
	`, GroupNum)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return 1, nil
}

func (db *BabyPhotoDB) LeaveGroup(UserNum int, GroupNum int) (int64, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		DELETE FROM GroupUserInfo WHERE GroupNum=? AND UserNum=?
	`, GroupNum, UserNum)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return 1, nil
}

func (db *BabyPhotoDB) InviteGroup(UserNum int, GroupNum int, InviteUserNum int, AbleUpload string, AbleDelete string, AbleView string) (int64, error) {
	GI, err := db.GroupUserInfo(UserNum, GroupNum)
	if err != nil {
		return -1, err
	}
	if GI.IsAdmin != "Y" {
		return 0, nil
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO GroupUserInfo (UserNum, GroupNum, IsAdmin, AbleUpload, AbleDelete, AbleView, GUJoinDtm, GUUpdateDtm) VALUES (?, ?, 'N', ?, ?, ?, ?, ?) 
	`, InviteUserNum, GroupNum, AbleUpload, AbleDelete, AbleView, util.CurrentDateTime(), util.CurrentDateTime())
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return 1, nil
}

func (db *BabyPhotoDB) MyGroupList(UserNum int) ([]model.GroupList, error) {
	rows, err := db.DB.Query(`
		SELECT
			A.GroupNum as GroupNum,
			A.GroupName as GroupName,
			B.IsAdmin as IsAdmin,
			B.AbleUpload as AbleUpload,
			B.AbleDelete as AbleDelete,
			B.AbleView as AbleView,
			C.count as GroupPeoPleCount,
			D.count as GroupFileCount,
			E.FilePath as FilePath
			FROM GroupInfo A, GroupUserInfo B, (
			SELECT 
				MAX(count) as count,
				GroupNum
			FROM (
				SELECT 
					count(*) as count,
					GroupNum  as GroupNum
				FROM GroupUserInfo
				WHERE GroupNum IN (
					SELECT
						B.GroupNum
					FROM GroupInfo A, GroupUserInfo B
					WHERE B.GroupNum = A.GroupNum
					AND B.GroupNum
				)
				GROUP BY GroupNum      
				UNION ALL
				SELECT
					0 as count,
					B.GroupNum as GroupNum
				FROM GroupInfo A, GroupUserInfo B
				WHERE B.GroupNum = A.GroupNum     
			) A
			GROUP BY A.GroupNum
			) C, (
			SELECT 
				MAX(count) as count,
				GroupNum
			FROM (
				SELECT 
					count(*) as count,
					GroupNum  as GroupNum
				FROM GroupFileInfo
				WHERE GFDelete <> 'Y'
				  AND GroupNum IN (
					SELECT
						B.GroupNum
					FROM GroupInfo A, GroupUserInfo B
					WHERE B.GroupNum = A.GroupNum
					AND B.GroupNum
				)
				GROUP BY GroupNum      
				UNION ALL
				SELECT
					0 as count,
					B.GroupNum as GroupNum
				FROM GroupInfo A, GroupUserInfo B
				WHERE B.GroupNum = A.GroupNum
			) A
			GROUP BY A.GroupNum
			) D, (
				SELECT 
				MAX(FilePath) as FilePath,
				GroupNum
				FROM (
					SELECT
					A.FilePath as FilePath,
					B.GroupNum as GroupNum
					FROM FileInfo A, GroupFileInfo B
					WHERE A.FileNum = B.FileNum
					AND B.GFDelete <> 'Y'
					AND A.FileNum IN (
						SELECT 
							MAX(B.FileNum) AS FileNum
						FROM GroupFileInfo A, FileInfo B
						WHERE GroupNum IN (
							SELECT
								B.GroupNum
							FROM GroupInfo A, GroupUserInfo B
							WHERE B.GroupNum = A.GroupNum
							AND B.GroupNum
						)
						AND A.GFDelete = 'N'
						AND A.FileNum = B.FileNum
						GROUP BY GroupNum
					)
					UNION ALL
					SELECT
						'' as FilePath,
						B.GroupNum as GroupNum
					FROM GroupInfo A, GroupUserInfo B
					WHERE B.GroupNum = A.GroupNum
				) A
				GROUP BY A.GroupNum
			) E
			WHERE B.GroupNum = A.GroupNum
			AND B.GroupNum = C.GroupNum
			AND B.GroupNum = D.GroupNum
			AND B.GroupNum = E.GroupNum
			AND B.UserNum = ?
			AND B.IsAdmin = 'Y'
			ORDER BY A.GroupRegDtm
	`, UserNum)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	groupList := []model.GroupList{}
	for rows.Next() {
		groupModel := model.GroupList{}
		rows.Scan(&groupModel.GroupNum, &groupModel.GroupName, &groupModel.IsAdmin, &groupModel.AbleUpload, &groupModel.AbleDelete, &groupModel.AbleView, &groupModel.GroupPeopleCount, &groupModel.GroupFileCount, &groupModel.FilePath)
		groupList = append(groupList, groupModel)
	}
	return groupList, nil
}

func (db *BabyPhotoDB) InvitedGroupList(UserNum int) ([]model.GroupList, error) {
	rows, err := db.DB.Query(`
		SELECT
			A.GroupNum as GroupNum,
			A.GroupName as GroupName,
			B.IsAdmin as IsAdmin,
			B.AbleUpload as AbleUpload,
			B.AbleDelete as AbleDelete,
			B.AbleView as AbleView,
			C.count as GroupPeoPleCount,
			D.count as GroupFileCount,
			E.FilePath as FilePath
			FROM GroupInfo A, GroupUserInfo B, (
			SELECT 
				MAX(count) as count,
				GroupNum
			FROM (
				SELECT 
					count(*) as count,
					GroupNum  as GroupNum
				FROM GroupUserInfo
				WHERE GroupNum IN (
					SELECT
						B.GroupNum
					FROM GroupInfo A, GroupUserInfo B
					WHERE B.GroupNum = A.GroupNum
					AND B.GroupNum
				)
				GROUP BY GroupNum      
				UNION ALL
				SELECT
					0 as count,
					B.GroupNum as GroupNum
				FROM GroupInfo A, GroupUserInfo B
				WHERE B.GroupNum = A.GroupNum     
			) A
			GROUP BY A.GroupNum
			) C, (
			SELECT 
				MAX(count) as count,
				GroupNum
			FROM (
				SELECT 
					count(*) as count,
					GroupNum  as GroupNum
				FROM GroupFileInfo
				WHERE GFDelete <> 'Y'
				  AND GroupNum IN (
					SELECT
						B.GroupNum
					FROM GroupInfo A, GroupUserInfo B
					WHERE B.GroupNum = A.GroupNum
					AND B.GroupNum
				)
				GROUP BY GroupNum      
				UNION ALL
				SELECT
					0 as count,
					B.GroupNum as GroupNum
				FROM GroupInfo A, GroupUserInfo B
				WHERE B.GroupNum = A.GroupNum
			) A
			GROUP BY A.GroupNum
			) D, (
				SELECT 
				MAX(FilePath) as FilePath,
				GroupNum
				FROM (
					SELECT
					A.FilePath as FilePath,
					B.GroupNum as GroupNum
					FROM FileInfo A, GroupFileInfo B
					WHERE A.FileNum = B.FileNum
					AND B.GFDelete <> 'Y'
					AND A.FileNum IN (
						SELECT 
							MAX(B.FileNum) AS FileNum
						FROM GroupFileInfo A, FileInfo B
						WHERE GroupNum IN (
							SELECT
								B.GroupNum
							FROM GroupInfo A, GroupUserInfo B
							WHERE B.GroupNum = A.GroupNum
							AND B.GroupNum
						)
						AND A.GFDelete = 'N'
						AND A.FileNum = B.FileNum
						GROUP BY GroupNum
					)
					UNION ALL
					SELECT
						'' as FilePath,
						B.GroupNum as GroupNum
					FROM GroupInfo A, GroupUserInfo B
					WHERE B.GroupNum = A.GroupNum
				) A
				GROUP BY A.GroupNum
			) E
			WHERE B.GroupNum = A.GroupNum
			AND B.GroupNum = C.GroupNum
			AND B.GroupNum = D.GroupNum
			AND B.GroupNum = E.GroupNum
			AND B.UserNum = ?
			AND B.IsAdmin = 'N'
			ORDER BY A.GroupRegDtm
	`, UserNum)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	groupList := []model.GroupList{}
	for rows.Next() {
		groupModel := model.GroupList{}
		rows.Scan(&groupModel.GroupNum, &groupModel.GroupName, &groupModel.IsAdmin, &groupModel.AbleUpload, &groupModel.AbleDelete, &groupModel.AbleView, &groupModel.GroupPeopleCount, &groupModel.GroupFileCount, &groupModel.FilePath)
		groupList = append(groupList, groupModel)
	}
	return groupList, nil
}
