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
