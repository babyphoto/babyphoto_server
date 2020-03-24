package model

type GroupNextVal int

type GroupInfo struct {
	GroupNum    int
	GroupName   string
	GroupRegDtm string
}

type GroupUserInfo struct {
	UserNum     int
	GroupNum    int
	IsAdmin     string
	AbleUpload  string
	AbleDelete  string
	AbleView    string
	GUJoinDtm   string
	GUUpdateDtm string
}

type GroupList struct {
	GroupNum         int
	GroupName        string
	IsAdmin          string
	AbleUpload       string
	AbleDelete       string
	AbleView         string
	GroupPeopleCount int
	GroupFileCount   int
	FilePath         string
}
