package babyphoto

import (
	"github.com/babyphoto/babyphoto_server/service/model"
	"github.com/babyphoto/babyphoto_server/service/util"
)

func (db *BabyPhotoDB) FileNextSerialNum() (int, error) {
	count := 0
	err := db.DB.QueryRow(`
		SELECT MAX(FileNum) + 1 as count FROM FileInfo
	`).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *BabyPhotoDB) InsertFile(UserNum int, GroupNum int, FileName string, FilePath string, FileThumbnail string, FileExtention string) (int, error) {
	NextID, err := db.FileNextSerialNum()
	if err != nil {
		return -1, err
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO FileInfo (FileNum, UserNum, FileName, FilePath, FileThumbnail, FileExtention, FileRegDtm) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, NextID, UserNum, FileName, FilePath, FileThumbnail, FileExtention, util.CurrentDateTime())
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec(`
		INSERT INTO GroupFileInfo(FileNum, GroupNum, GFDelete, GFJoinDtm, GFUpdateDtm) VALUES(?, ?, 'N', ?, ?)
	`, NextID, GroupNum, util.CurrentDateTime(), util.CurrentDateTime())
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return 1, nil
}

func (db *BabyPhotoDB) FileList(GroupNum int) ([]model.FileInfo, error) {
	rows, err := db.DB.Query(`
		SELECT
			A.*
		FROM FileInfo A, GroupFileInfo B
		WHERE A.FileNum = B.FileNum
		AND B.GroupNum = ?
		AND B.GFDelete <> 'Y'
		ORDER BY FileRegDtm DESC
	`, GroupNum)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	fileList := []model.FileInfo{}
	for rows.Next() {
		fileModel := model.FileInfo{}
		rows.Scan(&fileModel.FileNum, &fileModel.UserNum, &fileModel.FileName, &fileModel.FilePath, &fileModel.FileThumbnail, &fileModel.FileExtention, &fileModel.FileRegDtm)
		fileList = append(fileList, fileModel)
	}
	return fileList, nil
}

func (db *BabyPhotoDB) UpdateFile(UserNum int, FileNum int, GroupNum int) (int, error) {
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
		UPDATE GroupFileInfo SET GFDelete='Y', GFUpdateDtm=? WHERE FileNum=? AND GroupNum=?
	`, util.CurrentDateTime(), FileNum, GroupNum)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return 1, nil
}
