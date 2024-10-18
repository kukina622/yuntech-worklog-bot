package util

import (
	"net/http"
	"os"
	"strings"
)

func GetExternalCookie() (string, error) {
	cookie, err := os.ReadFile("./cookie.txt")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(cookie)), err
}

func GetRequestWithCookie(url string, method string) (*http.Request, error) {

	cookie, err := GetExternalCookie()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", cookie)

	return req, nil
}

func ParseCookieString(cookieStr string) []*http.Cookie {
	cookies := []*http.Cookie{}
	pairs := strings.Split(cookieStr, ";")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			cookie := &http.Cookie{
				Name:  strings.TrimSpace(kv[0]),
				Value: strings.TrimSpace(kv[1]),
			}
			cookies = append(cookies, cookie)
		}
	}
	return cookies
}
