package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Translate(word string) {
	url := "https://dict.youdao.com/suggest?num=5&ver=3.0&doctype=json&cache=false&le=en&q=" + word
	resp, err := http.Get(url)
	errCheck(err)

	defer resp.Body.Close()
	bys, err := ioutil.ReadAll(resp.Body)
	errCheck(err)

	var respData YouDaoTranslateResp
	err = json.Unmarshal(bys, &respData)
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
		fmt.Fprintf(os.Stderr, "请求出错: %s", err)
		os.Exit(0)
	}
}
