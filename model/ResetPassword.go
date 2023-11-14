package model

type ResetPassword struct {
	Token    string `form:"token"`
	Password string `form:"password"`
}
