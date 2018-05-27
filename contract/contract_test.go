package contract

import (
	"fmt"
	"testing"

	"github.com/bottos-project/bottos/contract/msgpack"
)


func TestTransfer(t *testing.T) {
	type transferparam struct {
		From		string
		To			string
		Value		uint64
	}

	param := transferparam {
		From: "delegate1",
		To: "delegate2",
		Value: 100,
	}
	data, _ :=  msgpack.Marshal(param)
	fmt.Printf("transfer struct: %v, msgpack: %x\n", param, data)
}

func TestNewAccount(t *testing.T) {
	type newaccountparam struct {
		Name		string
		Pubkey		string
	}

	param := newaccountparam {
		Name: "testuser",
		Pubkey: "7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF",
	}
	data, _ :=  msgpack.Marshal(param)
	fmt.Printf("transfer struct: %v, msgpack: %x\n", param, data)
}


func TestGrantCredit(t *testing.T) {
	type GrantCreditParam struct {
		Name		string		`json:"name"`
		Spender		string 		`json:"spender"`
		Limit		uint64		`json:"limit"`
	}

	type CancelCreditParam struct {
		Name		string		`json:"name"`
		Spender		string 		`json:"spender"`
	}
	
	type TransferFromParam struct {
		From		string		`json:"from"`
		To			string		`json:"to"`
		Value		uint64		`json:"value"`
	}

	param := GrantCreditParam {
		Name: "alice",
		Spender: "bob",
		Limit: 100,
	}
	data, _ :=  msgpack.Marshal(param)
	fmt.Printf("grant credit struct: %v, msgpack: %x\n", param, data)

	param1 := CancelCreditParam {
		Name: "alice",
		Spender: "bob",
	}
	data, _ =  msgpack.Marshal(param1)
	fmt.Printf("cancel credit struct: %v, msgpack: %x\n", param1, data)

	param2 := TransferFromParam {
		From: "alice",
		To: "toliman",
		Value: 150,
	}
	data, _ =  msgpack.Marshal(param2)
	fmt.Printf("transfer from credit struct: %v, msgpack: %x\n", param2, data)

}
