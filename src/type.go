package src

type YouDaoTranslateResp struct {
	Result GoResult `json:"result"`
	Data   GoData   `json:"data"`
}

type GoResult struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type GoEntries struct {
	Explain string `json:"explain"`
	Entry   string `json:"entry"`
}

type GoData struct {
	Entries  []GoEntries `json:"entries"`
	Query    string      `json:"query"`
	Language string      `json:"language"`
	Type     string      `json:"type"`
}
