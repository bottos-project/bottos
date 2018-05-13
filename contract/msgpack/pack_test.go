package msgpack

import (
	"fmt"
	//"bytes"
	"testing"
	"encoding/hex"
)


func BytesToHex(d []byte) string {
	return hex.EncodeToString(d)
}


func HexToBytes(str string) ([]byte, error) {
	h, err := hex.DecodeString(str)

	return h, err
}



func TestMarshalStruct(t *testing.T) {
	type TestStruct struct{
		V1 string
		V2 uint32
	}

	ts := TestStruct {
		V1: "testuser",
		V2: 99,
	}

	b, err := Marshal(ts)
	
	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1)
}


func TestMarshalNestStruct1(t *testing.T) {
	type TestSubStruct struct{
		V1 string
		V2 uint32
	}

	type TestStruct struct{
		V1 string
		V2 uint32
		V3 TestSubStruct
	}
	fmt.Println("TestMarshalNestStruct1...")

	ts := TestStruct {
		V1: "testuser",
		V2: 99,
		V3: TestSubStruct{V1:"123", V2:3},
	}
	b, err := Marshal(ts)
	
	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestMarshalNestStruct2(t *testing.T) {
	type TestSubStruct struct{
		V1 string
		V2 uint32
	}

	type TestStruct struct{
		V1 string
		V2 uint32
		V3 *TestSubStruct
	}
	fmt.Println("TestMarshalNestStruct2...")

	ts := TestStruct {
		V1: "testuser",
		V2: 99,
		V3: &TestSubStruct{V1:"123", V2:3},
	}
	b, err := Marshal(ts)
	
	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestMarshalNestStruct3(t *testing.T) {
	type TestSubStruct struct{
		V1 string
		V2 uint32
	}

	type TestStruct struct{
		V1 string
		V2 uint32
		V3 TestSubStruct
		V4 []byte
	}
	fmt.Println("TestMarshalNestStruct3...")

	ts := TestStruct {
		V1: "testuser",
		V2: 99,
		V3: TestSubStruct{V1:"123", V2:3},
		V4: []byte{99,21,22},
	}
	b, err := Marshal(ts)
	
	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("Unmarshal, ts: ", ts1, err)
}

func TestPackMarshalReguser(t *testing.T) {

	type RegUser struct{
		V1 string
		V2 string
	}

	fmt.Println("TestPackMarshalReguser...")

	ts := RegUser {
		V1: "did:bot:21tDAKCERh95uGgKbJNHYp",
		V2: "this is a test string",
	}
	b, err := Marshal(ts)
	
	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)
}

func TestTransfer(t *testing.T) {

	type Transfer struct{
		From string
		To string
		Value uint64
	}

	fmt.Println("TestTransfer...")

	ts := Transfer {
		From: "delegate1",
		To: "delegate2",
		Value: 9999,
	}
	b, err := Marshal(ts)
	
	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)
}

func TestMarshalNewAccount(t *testing.T) {
	type newaccountparam struct {
		Name		string
		Pubkey		string
	}
	param := newaccountparam {
		Name: "testuser",
		Pubkey: "7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF",
	}

	fmt.Println("TestMarshalNewAccount...")
	b, err := Marshal(param)
	
	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)
	
}
