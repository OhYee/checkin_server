package api

import (
	"github.com/OhYee/blotter/register"
)

func Register() {
	register.Register("get", GetData)
	register.Register("set", SetData)
	register.Register("login", Login)
	register.Register("logout", Logout)
	register.Register("info", Info)
	register.Register("checkin", CheckIn)
}
