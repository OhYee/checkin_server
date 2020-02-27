package api

import (
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
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

	password := user.PasswordHash(args.Username + args.Password)
	output.Log("%s", password)

	count, err := mongo.Find("checkin", "users", bson.M{
		"username": args.Username,
		"password": password,
	}, nil, nil)

	if count != 0 {
		res.Token = user.GenerateToken()
		if _, err = mongo.Update("checkin", "users", bson.M{
			"username": args.Username,
		}, bson.M{
			"$set": bson.M{"token": res.Token},
		}, nil); err != nil {
			return err
		}
	}

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

	result, err := mongo.Update("checkin", "users", bson.M{
		"token": args.Token,
	}, bson.M{
		"token": "",
	}, nil)

	res.Success = result.ModifiedCount != 0

	err = context.ReturnJSON(res)
	return
}
