package handler

type ResponseStruct struct {
	Errcode uint32      `json:"errcode"`
	Msg     string      `json:"msg"`
	Result  interface{} `json:"result"`
}

type ResponseStructs []ResponseStruct

type NewAccount struct {
	Account   string ` json:"account"`
	PublicKey string `json:"public_key"`
}
