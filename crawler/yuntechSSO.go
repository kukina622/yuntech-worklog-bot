package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	LOGIN_URL       = "https://webapp.yuntech.edu.tw/YunTechSSO/Account/Login"
	CHECK_LOGIN_URL = "https://webapp.yuntech.edu.tw/YunTechSSO/Account/IsLogined"
)

type YunTechSSOCrawler struct {
	Client *http.Client
}

func (crawler *YunTechSSOCrawler) Login() bool {
	cookie, err := crawler.getExternalCookie()

	if err != nil {
		fmt.Println("[yunTechSSOCrawler] Please complete or create cookie.txt")
		return false
	}

	return crawler.checkLogin(cookie)
}


func (crawler *YunTechSSOCrawler) checkLogin(cookie string) bool {

	req, err := http.NewRequest("GET", CHECK_LOGIN_URL, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Cookie", cookie)

	resp, err := crawler.Client.Do(req)

	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body) == "True"
}

func (crawler *YunTechSSOCrawler) getExternalCookie() (string, error) {
	cookie, err := os.ReadFile("./cookie.txt")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(cookie)), err
}
