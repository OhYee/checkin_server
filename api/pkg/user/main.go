package user

import (
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckToken(token string) (username string) {
	defer func() { output.Debug("%+v %+v", token, username) }()
	users := make([]map[string]interface{}, 0)
	count, err := mongo.Find("checkin", "users", bson.M{
		"token": token,
	}, nil, &users)
	if err != nil {
		output.Err(err)
		return
	}
	output.Log("%+v", users)
	if count != 0 {
		username = users[0]["username"].(string)
	}
	return
}

func Login(username string, password string) (token string) {
	password = user.PasswordHash(username + password)
	output.Log("[%s] [%s]", username, password)

	count, err := mongo.Find("checkin", "users", bson.M{
		"username": username,
		"password": password,
	}, nil, nil)

	if count != 0 {
		token = user.GenerateToken()
		if _, err = mongo.Update("checkin", "users", bson.M{
			"username": username,
		}, bson.M{
			"$set": bson.M{"token": token},
		}, nil); err != nil {
			output.Err(err)
			token = ""
		}
	}
	return
}

func Logout(token string) bool {
	result, err := mongo.Update("checkin", "users", bson.M{
		"token": token,
	}, bson.M{
		"token": "",
	}, nil)
	return err != nil && result.ModifiedCount != 0
}

func GetDay(username string) (day int64, lottery int64) {
	users := make([]struct {
		day     int64
		lottery int64
	}, 0)
	count, err := mongo.Find("checkin", "users", bson.M{
		"username": username,
	}, nil, &users)
	if err != nil || count == 0 {
		return 0, 0
	}
	return users[0].day, users[0].lottery
}

func SetDay(username string, day int64, lottery int64) {
	mongo.Update("checkin", "users", bson.M{
		"username": username,
	}, bson.M{
		"$set": bson.M{
			"day":     day,
			"lottery": lottery,
		},
	}, nil)
}

func SetMoney(username string, money int64) {
	mongo.Update("checkin", "users", bson.M{
		"username": username,
	}, bson.M{
		"$inc": bson.M{
			"money": money,
		},
	}, nil)
}

func Info(username string) (money int64, day int64, lottery int64, err error) {
	users := make([]struct {
		Money   int64 `bson:"money"`
		Day     int64 `bson:"day"`
		Lottery int64 `bson:"lottery"`
	}, 0)
	count, err := mongo.Find("checkin", "users", bson.M{
		"username": username,
	}, nil, &users)
	if err != nil || count == 0 {
		output.Debug("%s %+v", username, users)
		err = errors.New("Unexcepted error")
		return
	}
	money = users[0].Money
	day = users[0].Day
	lottery = users[0].Lottery
	return
}
