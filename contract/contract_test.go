package contract

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestContract(t *testing.T) {
	param := TransferParam {
		From: "delegate1",
		To: "delegate2",
		Value: 10000,
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
