package handler

import "time"

type Todo struct {
	Errcode uint32 `json:"errcode"`
	Msg     string `json:"msg"`
	Result    interface{} `json:"result"`


	Name      string      `json:"name"`
	Completed bool        `json:"completed"`
	Due       time.Time   `json:"due"`
}

type Todos []Todo
