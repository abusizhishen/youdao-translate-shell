package src

import (
  "net/http"
  "net/url"
)

var (
  addWordUrl = "https://dict.youdao.com/wordbook/webapi/v2/ajax/add?lan=en&word="
)

func add(word string) {
  http.Get(addWordUrl + url.QueryEscape(word))
}
