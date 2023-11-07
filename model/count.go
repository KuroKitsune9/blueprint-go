package model

type Count struct {
	Pending  int `json:"pending"`
	Progress int `json:"progress"`
	Done     int `json:"done"`
}
