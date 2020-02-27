package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type GetDataRequest struct {
	Date  string `json:"date"`
	Token string `json:"token"`
}

type GetDataResponse struct {
	Weight int64  `json:"weight"`
	Note   string `json:"note"`
	Money  int64  `json:"money"`
}

func GetData(context *register.HandleContext) (err error) {
	defer errors.Wrapper(&err)

	args := new(GetDataRequest)
	res := new(GetDataResponse)
	data := make([]GetDataResponse, 0)

	context.RequestParams(args)

	count, err := mongo.Find("checkin", "date", bson.M{
		"data": args.Date,
	}, nil, &data)

	if count == 0 {
		res.Money = -1
		res.Weight = 0
		res.Note = ""
	} else {
		*res = data[0]
	}

	err = context.ReturnJSON(res)
	return
}
