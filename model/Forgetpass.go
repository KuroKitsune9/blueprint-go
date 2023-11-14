package model

import "time"

type ForgotPassword struct {
	UserId    int       `db:"user_id"`
	Token     string    `db:"token"`
	ExpiredAt time.Time `db:"expired_at"`
}
