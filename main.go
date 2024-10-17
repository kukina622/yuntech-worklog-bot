package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yuntech-worklog-bot/bot"
	"yuntech-worklog-bot/crawler"
	"yuntech-worklog-bot/util"

	"github.com/fatih/color"
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
	job.AddFunc(config.Section("program").Key("execFreq").String(), func() {
		config, _ = ini.ShadowLoad("config.ini")
		fmt.Println("\n[crontab] Task starting!!!")
		task(jar, config)
	})
	fmt.Println("[crontab] Crontab running with", config.Section("program").Key("execFreq").String())
	job.Start()
	select {}
}

func task(jar *cookiejar.Jar, config *ini.File) {
	var workList []string = config.Section("work").Key("work").ValueWithShadows()

	yunTechSSOCrawler := crawler.YunTechSSOCrawler{
		Client: &http.Client{Jar: jar},
	}

	for i := 0; i < len(workList); i++ {
		workItem := strings.Split(workList[i], ",")

		var workDay time.Time
		var workType string

		if strings.Contains(workItem[1], "/") {
			workDay, _ = time.Parse("2006/01/02 -0700", workItem[1]+" +0800")
			workType = "(Once)"
		} else {
			now := time.Now()
			zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			targetWeek, _ := strconv.Atoi(workItem[1])
			workDay = zeroTime.AddDate(0, 0, targetWeek-int(time.Now().Weekday()))
			workType = "(Weekly)"
		}

		if time.Now().After(workDay) && workItem[len(workItem)-1] != workDay.Format("2006/01/02") {
			fmt.Println("\n[crontab] Work detected at", workDay.Format("2006/01/02"), "with", workList[i], workType)
			fmt.Println("[yunTechSSOCrawler] Try to login yuntech SSO...")
			if loginResult := yunTechSSOCrawler.Login(); !loginResult {
				color.Red("[yunTechSSOCrawler] Login Failed! Please Check Your Cookie!!!")
				return
			}

			fmt.Println("[yunTechSSOCrawler] Login successfully")

			startTimeText := workItem[2]
			endTimeText := workItem[3]

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
			file, _ := os.Open("config.ini")
			var lines []string
			scanner := bufio.NewScanner(file)
			if strings.Contains(workItem[1], "/") {
				for scanner.Scan() {
					line := scanner.Text()
					if strings.Contains(line, workList[i]) && strings.HasPrefix(line, "work =") {
						line = "# " + line
					}
					lines = append(lines, line)
				}
			} else {
				for scanner.Scan() {
					line := scanner.Text()
					if strings.Contains(line, workList[i]) && strings.HasPrefix(line, "work =") {
						re := regexp.MustCompile(`(.),?(\d{4}/\d{2}/\d{2})?$`)
						line = re.ReplaceAllString(line, "$1,"+workDay.Format("2006/01/02"))
					}
					lines = append(lines, line)
				}
			}
			file.Close()
			outputFile, _ := os.Create("config.ini")
			defer outputFile.Close()
			for _, line := range lines {
				outputFile.WriteString(line + "\n")
			}

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
