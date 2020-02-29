package api

import (
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/checkin-server/api/pkg/user"
	"github.com/OhYee/rainbow/errors"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(context *register.HandleContext) (err error) {
	defer errors.Wrapper(&err)

	args := new(LoginRequest)
	res := new(LoginResponse)

	context.RequestParams(args)

	res.Token = user.Login(args.Username, args.Password)

	err = context.ReturnJSON(res)
	return
}

type LogoutRequest struct {
	Token string `json:"token"`
}

type LogoutResponse Response

func Logout(context *register.HandleContext) (err error) {
	defer errors.Wrapper(&err)

	args := new(LogoutRequest)
	res := new(LogoutResponse)

	context.RequestParams(args)

	res.Success = user.Logout(args.Token)

	err = context.ReturnJSON(res)
	return
}

type InfoRequest struct {
	Token string `json:"token"`
}

type InfoResponse struct {
	Money   int64 `json:"money"`
	Day     int64 `json:"day"`
	Lottery int64 `json:"lottery"`
}

func Info(context *register.HandleContext) (err error) {
	defer errors.Wrapper(&err)

	args := new(InfoRequest)
	res := new(InfoResponse)

	context.RequestParams(args)

	username := user.CheckToken(args.Token)
	if username == "" {
		context.Forbidden()
		return
	}
	// ================================
	if res.Money, res.Day, res.Lottery, err = user.Info(username); err != nil {
		return
	}

	err = context.ReturnJSON(res)
	return
}
