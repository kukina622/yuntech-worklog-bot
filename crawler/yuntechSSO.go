package crawler

import (
	"io/ioutil"
	"net/http"
)

const (
	LOGIN_URL       = "https://webapp.yuntech.edu.tw/YunTechSSO/Account/Login"
	CHECK_LOGIN_URL = "https://webapp.yuntech.edu.tw/YunTechSSO/Account/IsLogined"
)

type YunTechSSOCrawler struct {
	Client *http.Client
}

func (crawler *YunTechSSOCrawler) Login() bool {
	return crawler.checkLogin()
}

func (crawler *YunTechSSOCrawler) checkLogin() bool {

	req, err := http.NewRequest("GET", CHECK_LOGIN_URL, nil)

	if err != nil {
		panic(err)
	}

	resp, err := crawler.Client.Do(req)

	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body) == "True"
}
