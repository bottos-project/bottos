package contract

import (
	"encoding/json"
	"fmt"
	"testing"
)


type transferparam struct {
	From		string		`json:"from"`
	To			string		`json:"to"`
	Value		uint64		`json:"value"`
}

func TestTransfer(t *testing.T) {
	param := transferparam {
		From: "delegate1",
		To: "delegate2",
		Value: 100,
	}
	data, _ :=  json.Marshal(param)
	fmt.Printf("[")
	for i, v := range data {
		if i == len(data)-1 {
			fmt.Printf("%v", v)
		} else {
			fmt.Printf("%v,", v)
		}
		
	}
	fmt.Printf("]\n")
	fmt.Printf("%s\n", data)
}



type newaccountparam struct {
	Creator		string		`json:"creator"`
	Name		string		`json:"name"`
	Pubkey		string 		`json:"pubkey"`
	Deposit		uint64		`json:"deposit"`
}

func TestNewAccount(t *testing.T) {
	param := newaccountparam {
		Creator: "bottos",
		Name: "testuser",
		Pubkey: "7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF",
		Deposit: 10000,
	}
	data, _ :=  json.Marshal(param)
	fmt.Printf("[")
	for i, v := range data {
		if i == len(data)-1 {
			fmt.Printf("%v", v)
		} else {
			fmt.Printf("%v,", v)
		}
		
	}
	fmt.Printf("]\n")
	fmt.Printf("%s\n", data)
}
