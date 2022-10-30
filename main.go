package main

import (
	"strconv"
	"strings"
	"time"
	"yuntech-worklog-bot/crawler"
	// "log"
	// "os"
	"net/http"
	"net/http/cookiejar"
	// "github.com/robfig/cron/v3"
	"gopkg.in/ini.v1"
	"yuntech-worklog-bot/util"
)

func main() {
	// logger := cron.VerbosePrintfLogger(log.New(os.Stdout, "", log.LstdFlags))
	jar, _ := cookiejar.New(nil)
	// job := cron.New(cron.WithChain(cron.Recover(logger)))
	task(jar)
	// job.AddFunc("", func() {

	// })
	// 20 12 * * */1
}

func task(jar *cookiejar.Jar) {
	cfg, _ := ini.ShadowLoad("config.ini")
	var workList []string = cfg.Section("work").Key("work").ValueWithShadows()
	userConfig := cfg.Section("user")

	yunTechSSOCrawler := crawler.YunTechSSOCrawler{
		Username: userConfig.Key("username").String(),
		Password: userConfig.Key("password").String(),
		Client:   &http.Client{Jar: jar},
	}

	for i := 0; i < len(workList); i++ {
		workItem := strings.Split(workList[i], ",")
		workWeekday, _ := strconv.Atoi(workItem[1])

		if int(time.Now().Weekday())-1 == workWeekday {
			yunTechSSOCrawler.Login()

			startTimeText := workItem[2]
			endTimeText := workItem[3]

			workDay := time.Now().AddDate(0, 0, -1)

			workLogCrawler := crawler.WorkLogCrawler{
				YunTechSSOCrawler: yunTechSSOCrawler,
				WorkName:          workItem[0],
				WorkContent:       workItem[4],
				StartTime:         util.ApplyTimeByTimeText(workDay, startTimeText),
				EndTime:           util.ApplyTimeByTimeText(workDay, endTimeText),
			}
			workLogCrawler.FillOutWorkLog()
		}

	}

}
