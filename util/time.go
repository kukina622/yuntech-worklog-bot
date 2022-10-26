package util

import (
	"strconv"
	"strings"
	"time"
)

func ApplyTimeByTimeText(day time.Time, text string) time.Time {
	timeZone, _ := time.LoadLocation("Asia/Taipei")
	timeTextList := strings.Split(text, ":")
	hour, _ := strconv.Atoi(timeTextList[0])
	minute, _ := strconv.Atoi(timeTextList[1])
	return time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, timeZone)
}
