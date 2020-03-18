package model

type DBConnection struct {
	Database string `json:"database"`
	User     struct {
		ID  string `json:"id"`
		Pwd string `json:"pwd"`
	} `json:"user"`
	Host struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	} `json:"host"`
}
