package model

type UserInfo struct {
	UserNum      int
	UserCode     string
	UserType     string
	UserNickName string
	UserName     string
	UserRegDtm   string
	UserProfile  string
}

type GroupUserList struct {
	UserNum      int
	UserCode     string
	UserType     string
	UserNickName string
	UserName     string
	UserProfile  string
	GroupNum     int
	IsAdmin      string
	AbleUpload   string
	AbleDelete   string
	AbleView     string
	GUJoinDtm    string
}
