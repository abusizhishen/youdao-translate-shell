package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-resty/resty/v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func Translate(word string) {
	queryUrl := "https://dict.youdao.com/suggest?num=5&ver=3.0&doctype=json&le=en&q=" + url.QueryEscape(word)

	resp, err := resty.New().R().EnableTrace().SetCookies(loadCookie()).Get(queryUrl)
	errCheck(err)

	var respData YouDaoTranslateResp
	fmt.Println(resp.Cookies())
	err = json.Unmarshal(resp.Body(), &respData)
	errCheck(err)

	if respData.Result.Code != 200 {
		fmt.Fprintf(os.Stdout, "请求出错: %s", respData.Result.Msg)
		os.Exit(0)
	}

	data := respData.Data

	//fmt.Fprintf(os.Stdout, "- %s\n", data.Query)

	for _, entry := range data.Entries {
		fmt.Fprintf(os.Stdout, "- %s\n- %s\n\n", entry.Entry, entry.Explain)
	}
}

func errCheck(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(os.Stderr, "请求出错: %s", err)
		os.Exit(0)
	}
}

// {"msg":"NO_LOGIN","code":2035}
//{"msg":"OK","code":0,"data":{"compatible_yduserid":"urs-phoneyd.9e39c49e863a45519@163.com","login_method":"urstoken","yduserid":"urs-phoneyd.9e39c49e863a45519@163.com","userid":"m13571981487@163.com","third_part_info":[{"platform":"ursphone","username":"13571981487"},{"platform":"urstoken","username":"13571981487"},{"platform":"weixin","username":"无形"}],"ssn":"yd.9e39c49e863a45519@163.com"}}
func CheckLoginStatus() (AccountInfo, error) {
	accountInfoUrl := "https://dict.youdao.com/login/acc/query/accountinfo"

	req, _ := http.NewRequest("GET", accountInfoUrl, nil)
	fmt.Println("load")
	cookies := loadCookie()
	fmt.Println("set")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	//fmt.Println("print", cookies)
	//time.Sleep(time.Second * 3)
	//resp, err := resty.New().R().EnableTrace().SetCookies(cookies).Get(accountInfoUrl)
	//errCheck(err)

	for _, c := range cookies {
		fmt.Printf("%+v\n", c.Name)
		req.AddCookie(c)
	}
	result := AccountInfo{}
	resp, err := (&http.Client{}).Do(req)
	bys, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bys, &result)
	errCheck(err)
	if result.Code != 0 {
		return result, errors.New(result.Msg)
	}

	return result, nil
}

type AccountInfo struct {
	Msg  string          `json:"msg"`
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

func loadCookie() []*http.Cookie {
	// 如果存在则读取cookies的数据
	cookiesData, err := ioutil.ReadFile(cookieFileName)
	if err != nil {
		errCheck(err)
		return nil
	}

	// 反序列化
	cookies := []Cookie{}
	if err = json.Unmarshal(cookiesData, &cookies); err != nil {
		errCheck(err)
		return nil
	}

	var result []*http.Cookie

	for _, c := range cookies {
		result = append(result, &http.Cookie{
			Name:       c.Name,
			Value:      strings.Trim(c.Value, "\""),
			Path:       c.Path,
			Domain:     c.Domain,
			Expires:    time.Unix(int64(c.Expires), 0),
			RawExpires: ``,
			MaxAge:     0,
			Secure:     false,
			HttpOnly:   c.HTTPOnly,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		})
	}

	return result
}

type Cookie struct {
	Name         string  `json:"name"`
	Value        string  `json:"value"`
	Domain       string  `json:"domain"`
	Path         string  `json:"path"`
	Expires      float64 `json:"expires"`
	Size         int     `json:"size"`
	HTTPOnly     bool    `json:"httpOnly"`
	Secure       bool    `json:"secure"`
	Session      bool    `json:"session"`
	Priority     string  `json:"priority"`
	SameParty    bool    `json:"sameParty"`
	SourceScheme string  `json:"sourceScheme"`
	SourcePort   int     `json:"sourcePort"`
}

var headers = map[string]string{
	`Accept`: `application/json, text/plain, */*`,
	//`Accept-Encoding`: `gzip, deflate, br`,
	`Accept-Language`:    `zh-CN,zh-TW;q=0.9,zh;q=0.8`,
	`Connection`:         `keep-alive`,
	`DNT`:                `1`,
	`Host`:               `dict.youdao.com`,
	`Origin`:             `https://www.youdao.com`,
	`Referer`:            `https://www.youdao.com/`,
	`sec-ch-ua`:          `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`,
	`sec-ch-ua-mobile`:   `?0`,
	`sec-ch-ua-platform`: `"macOS"`,
	`Sec-Fetch-Dest`:     `empty`,
	`Sec-Fetch-Mode`:     `cors`,
	`Sec-Fetch-Site`:     `same-site`,
	`User-Agent`:         `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36`,
}
