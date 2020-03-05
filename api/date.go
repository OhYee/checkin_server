package api

import (
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/checkin-server/api/pkg/date"
	"github.com/OhYee/checkin-server/api/pkg/user"
	"github.com/OhYee/rainbow/errors"
)

type GetDataRequest struct {
	Date  int64  `json:"date"`
	Token string `json:"token"`
}

type GetDataResponse = date.Data

func GetData(context *register.HandleContext) (err error) {
	defer errors.Wrapper(&err)

	args := new(GetDataRequest)
	res := new(GetDataResponse)

	context.RequestParams(args)

	username := user.CheckToken(args.Token)
	if username == "" {
		context.Forbidden()
		return
	}
	// ================================

	if *res, err = date.Get(username, args.Date); err != nil {
		return
	}

	err = context.ReturnJSON(res)
	return
}

type SetDataRequest struct {
	Date  int64  `json:"date"`
	Token string `json:"token"`
	date.Basic
}

type SetDataResponse Response

func SetData(context *register.HandleContext) (err error) {
	defer errors.Wrapper(&err)

	args := new(SetDataRequest)
	res := new(SetDataResponse)

	context.RequestData(args)

	username := user.CheckToken(args.Token)
	if username == "" {
		context.Forbidden()
		return
	}
	// ================================

	if err = date.Set(username, args.Date, args.Weight, args.Note); err != nil {
		return
	}
	res.Success = true
	err = context.ReturnJSON(res)
	return
}

type CheckInRequest struct {
	Date  int64  `json:"date"`
	Token string `json:"token"`
}

type CheckInResponse struct {
	Money   int64 `json:"money"`
	Day     int64 `json:"day"`
	Lottery int64 `json:"lottery"`
}

func CheckIn(context *register.HandleContext) (err error) {
	defer errors.Wrapper(&err)

	args := new(CheckInRequest)
	res := new(CheckInResponse)

	context.RequestParams(args)

	username := user.CheckToken(args.Token)
	if username == "" {
		context.Forbidden()
		return
	}
	// ================================

	res.Money, res.Day, res.Lottery, err = date.CheckIn(username, args.Date)
	if err != nil {
		return
	}
	output.Debug("%+v", err)

	err = context.ReturnJSON(res)
	return
}

type MenstruationRequest struct {
	Date  int64  `json:"date"`
	Token string `json:"token"`
}

type MenstruationResponse Response

func Menstruation(context *register.HandleContext) (err error) {
	defer errors.Wrapper(&err)

	args := new(MenstruationRequest)
	res := new(MenstruationResponse)

	context.RequestParams(args)

	username := user.CheckToken(args.Token)
	if username == "" {
		context.Forbidden()
		return
	}
	// ================================

	err = date.Menstruation(username, args.Date)
	if err != nil {
		return
	}
	output.Debug("%+v", err)

	res.Success = true
	err = context.ReturnJSON(res)
	return
}
