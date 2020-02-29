package date

import (
	"math/rand"

	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/checkin-server/api/pkg/user"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func set(username string, date int64, insert bson.M) (err error) {
	defer errors.Wrapper(&err)

	datetime := timeToString(date)

	defaultValue := bson.M{
		"date":     datetime,
		"username": username,
		"money":    -1,
		"weight":   0,
		"note":     "",
	}

	dv := bson.M{}
	for key, value := range defaultValue {
		if _, exist := insert[key]; exist != true {
			dv[key] = value
		}
	}

	_, err = mongo.Update(
		"checkin",
		"date",
		bson.M{
			"date":     datetime,
			"username": username,
		},
		bson.M{
			"$setOnInsert": dv,
			"$set":         insert,
		},
		options.Update().SetUpsert(true),
	)
	return
}

func Set(username string, date int64, weight int64, note string) (err error) {
	err = set(
		username,
		date,
		bson.M{
			"weight": weight,
			"note":   note,
		},
	)
	return
}

func Get(username string, date int64) (data Data, err error) {
	datetime := timeToString(date)
	output.Debug("%s %s", username, datetime)
	_data := make([]Data, 0)
	count, err := mongo.Find("checkin", "date", bson.M{
		"date":     datetime,
		"username": username,
	}, nil, &_data)
	if err != nil {
		return
	}

	if count == 0 {
		data.Money = -1
		data.Weight = 0
		data.Note = ""
	} else {
		data = _data[0]
	}
	return
}

func CheckIn(username string, date int64) (money int64, day int64, lottery int64, err error) {
	defer errors.Wrapper(&err)

	todayData, err := Get(username, date)
	if err != nil {
		return
	}
	if todayData.Money != -1 {
		errors.New("Today has checked")
		return
	}

	yesterday := getYesterday(date)

	data, err := Get(username, yesterday)
	if err != nil {
		return
	}

	if data.Money != -1 {
		// 签到数增加
		_, err = mongo.Update("checkin", "users", bson.M{
			"username": username,
		}, bson.M{
			"$inc": bson.M{"day": 1},
		}, nil)
	} else {
		// 断签
		_, err = mongo.Update("checkin", "users", bson.M{
			"username": username,
		}, bson.M{
			"$set": bson.M{"day": 1},
		}, nil)
	}
	if err != nil {
		return
	}

	day, lottery = user.GetDay(username)
	if day >= 15 {
		day -= 15
		lottery++
		user.SetDay(username, day, lottery)
	}

	money = rand.Int63n(10 * 100)
	user.SetMoney(username, money)

	err = set(
		username,
		date,
		bson.M{
			"money": money,
		},
	)
	return
}
