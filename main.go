package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"time"
	"yuntech-worklog-bot/bot"
	"yuntech-worklog-bot/crawler"
	"yuntech-worklog-bot/util"
	"github.com/robfig/cron/v3"
	"gopkg.in/ini.v1"
)

func main() {
	config, _ := ini.ShadowLoad("config.ini")
	// discord
	discordConfig := config.Section("discord")
	enableBot, _ := discordConfig.Key("enableBot").Bool()
	if enableBot {
		botToken := discordConfig.Key("botToken").String()
		bot.GetDiscordBotInstance().InitSession(botToken)
	}
	// crontab
	logger := cron.VerbosePrintfLogger(log.New(os.Stdout, "", log.LstdFlags))
	jar, _ := cookiejar.New(nil)
	job := cron.New(cron.WithChain(cron.Recover(logger)))
	job.AddFunc("20 12 * * */1", func() {
		task(jar, config)
	})
}

func task(jar *cookiejar.Jar, config *ini.File) {
	var workList []string = config.Section("work").Key("work").ValueWithShadows()
	userConfig := config.Section("user")

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
			result := workLogCrawler.FillOutWorkLog()

			enableBot, _ := config.Section("discord").Key("enableBot").Bool()
			if result && enableBot {
				channelId := config.Section("discord").Key("channelID").String()
				message := workLogCrawler.GetFillSuccessMessage()
				bot.GetDiscordBotInstance().SendMessage(message, channelId)
			}
		}

	}

}
