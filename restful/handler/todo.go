package handler

type Todo struct {
	Errcode uint32 `json:"errcode"`
	Msg     string `json:"msg"`
	Result    interface{} `json:"result"`
}

type Todos []Todo
