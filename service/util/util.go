package util

import "log"

// CheckError 에러 공통 처리
func CheckError(comment string, err error) {
	if err != nil {
		log.Println("err : ", comment, err)
	}
}
