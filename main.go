package main

import (
	"fmt"
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
	discordConfig := config.Section("discordBot")
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
		fmt.Println("[crontab] Task starting!!!")
		task(jar, config)
	})
	fmt.Println("[crontab] Crontab running!!!")
	job.Start()
	select {}
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
			fmt.Println("[crontab] Work detected")
			fmt.Println("[yunTechSSOCrawler] Try to login yuntech SSO...")
			if loginResult := yunTechSSOCrawler.Login(); !loginResult {
				return
			}
			fmt.Println("[yunTechSSOCrawler] Login successfully")

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
			fmt.Println("[workLogCrawler] Fill out workLog...")
			result := workLogCrawler.FillOutWorkLog()
			if !result {
				return
			}
			fmt.Println("[workLogCrawler] Fill out successfully")

			enableBot, _ := config.Section("discordBot").Key("enableBot").Bool()
			if enableBot {
				channelId := config.Section("discordBot").Key("channelID").String()
				message := workLogCrawler.GetFillSuccessMessage()
				bot.GetDiscordBotInstance().SendMessage(message, channelId)
			}
			enableWebhook, _ := config.Section("discordWebhook").Key("enableWebhook").Bool()
			if enableWebhook {
				webhookURL := config.Section("discordWebhook").Key("webhookURL").String()
				message := workLogCrawler.GetFillSuccessMessage()
				bot.SendWebhookMessage(webhookURL, message)
			}
		}

	}

}
