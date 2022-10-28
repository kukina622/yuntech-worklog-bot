package crawler

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

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

func (crawler *WorkLogCrawler) loginWorkStudy() {
	crawler.Client.Get(WORKLOG_LOGIN)
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
