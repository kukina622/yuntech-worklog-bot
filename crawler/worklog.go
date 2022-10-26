package crawler

import "time"

const (
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
