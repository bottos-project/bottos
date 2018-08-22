package handler

type ResponseStruct struct {
	Errcode uint32 `json:"errcode"`
	Msg     string `json:"msg"`
	Result    interface{} `json:"result"`
}

type ResponseStructs []ResponseStruct
