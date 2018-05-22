package contract

import (
	"bytes"
	"encoding/json"
)

type ABI struct {
	Types   []interface{} `json:"types"`
	Structs []struct {
		Name   string `json:"name"`
		Base   string `json:"base"`
		Fields map[string]string `json:"fields"`
	} `json:"structs"`
	Actions []struct {
		ActionName string `json:"action_name"`
		Type       string `json:"type"`
	} `json:"actions"`
	Tables []interface{} `json:"tables"`
}


func ParseAbi(abiRaw []byte) (*ABI, error) {
	abi := &ABI{}
	err := json.Unmarshal(abiRaw, abi)
	if err != nil {
		return  &ABI{}, err
	}

	return abi, nil
}

func AbiToJson(abi *ABI) (string, error) {
	data, err := json.Marshal(abi)
	if err != nil {
		return "", err
	}
	return jsonFormat(data), nil
}

func jsonFormat(data []byte) string {
	var out bytes.Buffer  
	json.Indent(&out, data, "", "    ")
	
	return string(out.Bytes())
}
