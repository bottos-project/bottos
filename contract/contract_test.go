package contract

import (
	"fmt"
	"testing"

	"github.com/bottos-project/core/contract/msgpack"
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
