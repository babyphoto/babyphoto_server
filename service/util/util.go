package util

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
