package crawler

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yuntech-worklog-bot/util"

	"github.com/anaskhan96/soup"
)

const (
	WORKLOG_LOGIN      = "https://webapp.yuntech.edu.tw/workstudy/Home/Login"
	WORKLOG_LIST_URL   = "https://webapp.yuntech.edu.tw/workstudy/StudWorkRecord/ContractList"
	WORKLOG_CREATE_URL = "https://webapp.yuntech.edu.tw/workstudy/StudWorkRecord/ApplyAction"
)

type WorkLogCrawler struct {
	YunTechSSOCrawler
	WorkName    string
	StartTime   time.Time
	EndTime     time.Time
	WorkContent string
}

func (crawler *WorkLogCrawler) FillOutWorkLog() bool {
	crawler.loginWorkStudy()
	workId := crawler.getWorkId()
	payload := crawler.getFormPayload(workId)
	_, err := crawler.Client.PostForm(WORKLOG_CREATE_URL, payload)
	return err == nil
}

func (crawler *WorkLogCrawler) loginWorkStudy() {
	resp, err := crawler.Client.Get(WORKLOG_LOGIN)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	doc := soup.HTMLParse(string(body))
	loginEndPointUrl := doc.Find("a").Attrs()["href"]
	crawler.Client.Get(loginEndPointUrl)
}

func (crawler *WorkLogCrawler) getWorkId() (workId string) {
	url := fmt.Sprintf("%s?date=%d/%d", WORKLOG_LIST_URL, crawler.StartTime.Year(), int(crawler.StartTime.Month()))
	resp, err := crawler.Client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	doc := soup.HTMLParse(string(body))
	workRowList := doc.Find("tbody").FindAll("tr")
	for i := 0; i < len(workRowList); i++ {
		workRow := workRowList[i].FindAll("td")
		if strings.Trim(workRow[2].Text(), " ") == crawler.WorkName {
			workhref := workRow[8].Find("a").Attrs()["href"]
			var rgx = regexp.MustCompile(`ContractId=\d+`)
			rs := rgx.FindStringSubmatch(workhref)
			workId = strings.Split(rs[0], "=")[1]
		}
	}
	return
}

func (crawler *WorkLogCrawler) getFormPayload(workId string) url.Values {
	workHours := fmt.Sprintf("%.1f", util.GetHourDiffer(crawler.StartTime, crawler.EndTime))
	day := strconv.Itoa(crawler.StartTime.Day())
	if crawler.StartTime.Day() < 10 {
		day = "0" + day
	}
	dateContract := fmt.Sprintf(
		"%d/%d/%s,%s",
		crawler.StartTime.Year(),
		int(crawler.StartTime.Month()),
		day,
		workId,
	)
	startTime := crawler.StartTime.Add(-time.Minute * getRandomTimeDuration(5))
	endTime := crawler.EndTime.Add(time.Minute * getRandomTimeDuration(5))

	payload := url.Values{}
	payload.Add("DateContract", dateContract)
	payload.Add("StartHour", strconv.Itoa(startTime.Hour()))
	payload.Add("StartMin", strconv.Itoa(startTime.Minute()))
	payload.Add("EndHour", strconv.Itoa(endTime.Hour()))
	payload.Add("EndMin", strconv.Itoa(endTime.Minute()))
	payload.Add("IsAnnualLeave", "false")
	payload.Add("WorkContent", crawler.WorkContent)
	payload.Add("Hours", workHours)
	return payload
}

func getRandomTimeDuration(n int) time.Duration {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return time.Duration(r1.Intn(n) + 1)
}

func (crawler *WorkLogCrawler) GetFillSuccessMessage() string {
	const timeLayout string = "15:04"
	const DayLayout string = "2006-01-02"
	workHours := fmt.Sprintf("%.1f", util.GetHourDiffer(crawler.StartTime, crawler.EndTime))
	return fmt.Sprintf(
		"[%s]%s %s~%s 共%s小時 填寫完成",
		crawler.WorkContent,
		crawler.StartTime.Format(DayLayout),
		crawler.StartTime.Format(timeLayout),
		crawler.EndTime.Format(timeLayout),
		workHours,
	)
}
