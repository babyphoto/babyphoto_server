package util

import (
	"fmt"
	"log"
	"time"
)

// CheckError 에러 공통 처리
func CheckError(comment string, err error) {
	if err != nil {
		log.Println("err : ", comment, err)
	}
}

func CurrentDateTime() string {
	t := time.Now()
	formatted := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return formatted
}

func ReturnMap(result interface{}) interface{} {
	response := map[string]interface{}{}
	response["result"] = result
	return response
}
