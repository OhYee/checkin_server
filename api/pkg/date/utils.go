package date

import (
	"time"
)

func timeToString(date int64) string {
	return time.Unix(date, 0).Format("2006-01-02")
}

func getYesterday(t int64) int64 {
	return time.Unix(t, 0).AddDate(0, 0, -1).Unix()
}
