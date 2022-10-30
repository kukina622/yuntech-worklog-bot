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


func GetHourDiffer(start_time, end_time time.Time) float64 {
	var hour float64
	if start_time.Before(end_time) {
		diff := end_time.Unix() - start_time.Unix()
		hour = float64(diff) / 3600.
	} 
	return hour
}
