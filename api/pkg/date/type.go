package date

type Type struct {
	Date     string `json:"date" bson:"date"`
	Username string `json:"username" bson:"username`
	Data     `bson:",inline"`
}

type Data struct {
	Money int64 `json:"money" bson:"money"`
	Basic `bson:",inline"`
}

type Basic struct {
	Weight int64  `json:"weight" bson:"weight"`
	Note   string `json:"note" bson:"note"`
}
