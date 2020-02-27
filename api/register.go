package api

import (
	"github.com/OhYee/blotter/register"
)

func Register() {
	register.Register("get", GetData)
	register.Register("login", Login)
}
