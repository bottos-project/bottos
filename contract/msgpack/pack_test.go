package msgpack

import (
	"fmt"
	//"bytes"
	"encoding/hex"
	"testing"
)

func BytesToHex(d []byte) string {
	return hex.EncodeToString(d)
}

func HexToBytes(str string) ([]byte, error) {
	h, err := hex.DecodeString(str)

	return h, err
}

func TestMarshalStruct(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 uint32
	}

	ts := TestStruct{
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
	type TestSubStruct struct {
		V1 string
		V2 uint32
	}

	type TestStruct struct {
		V1 string
		V2 uint32
		V3 TestSubStruct
	}
	fmt.Println("TestMarshalNestStruct1...")

	ts := TestStruct{
		V1: "testuser",
		V2: 99,
		V3: TestSubStruct{V1: "123", V2: 3},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestMarshalNestStruct2(t *testing.T) {
	type TestSubStruct struct {
		V1 string
		V2 uint32
	}

	type TestStruct struct {
		V1 string
		V2 uint32
		V3 *TestSubStruct
	}
	fmt.Println("TestMarshalNestStruct2...")

	ts := TestStruct{
		V1: "testuser",
		V2: 99,
		V3: &TestSubStruct{V1: "123", V2: 3},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestMarshalNestStruct3(t *testing.T) {
	type TestSubStruct struct {
		V1 string
		V2 uint32
	}

	type TestStruct struct {
		V1 string
		V2 uint32
		V3 TestSubStruct
		V4 []byte
	}
	fmt.Println("TestMarshalNestStruct3...")

	ts := TestStruct{
		V1: "testuser",
		V2: 99,
		V3: TestSubStruct{V1: "123", V2: 3},
		V4: []byte{99, 21, 22},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("Unmarshal, ts: ", ts1, err)
}

func TestPackMarshalReguser(t *testing.T) {

	type RegUser struct {
		V1 string
		V2 string
	}

	fmt.Println("TestPackMarshalReguser...")

	ts := RegUser{
		V1: "did:bot:21tDAKCERh95uGgKbJNHYp",
		V2: "this is a test string",
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)
}

func TestTransfer(t *testing.T) {

	type Transfer struct {
		From  string
		To    string
		Value uint64
	}

	fmt.Println("TestTransfer...")

	// marshal
	ts := Transfer{
		From:  "bottos",
		To:    "bot",
		Value: 100000000000,
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := &Transfer{}
	err = Unmarshal(b, ts1)
	fmt.Println("ts1: ", ts1)
	cc, _ := HexToBytes("dc0004da00057474747474da000461666166cf000000003b9aca00da000c417072696c27732072656e74")
	err = Unmarshal(cc, ts1)
	fmt.Println("ts1: ", ts1)
}

func TestNewAccount(t *testing.T) {
	type newaccountparam struct {
		Name   string
		Pubkey string
	}
	param := newaccountparam{
		Name:   "testuser",
		Pubkey: "7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF",
	}

	fmt.Println("TestNewAccount...")
	b, err := Marshal(param)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	param1 := &newaccountparam{}
	err = Unmarshal(b, param1)
	fmt.Println("param1: ", param1)
}

func TestDatafileReg(t *testing.T) {
	type TestSubStruct struct {
		V1 string
		V2 string
		V3 uint64
		V4 string
		V5 string
		V6 string
		V7 uint64
		V8 string
	}

	type TestStruct struct {
		V1 string
		V2 *TestSubStruct
	}
	fmt.Println("TestDatafileReg...")

	ts := TestStruct{
		V1: "12345678901234567890",
		V2: &TestSubStruct{V1: "salertest", V2: "sissidTest", V3: 111, V4: "filenameTest", V5: "filepolicytest", V6: "authpathtest", V7: 222, V8: "sign"},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestAssetfileReg(t *testing.T) {
	type TestSubStruct struct {
		UserName    string
		AssetName   string
		AssetType   uint64
		FeatureTag  string
		SampleHash  string
		StorageHash string
		ExpireTime  uint32
		OpType      uint32
		Price       uint64
		Description string
	}

	type TestStruct struct {
		AssetId string
		V2      TestSubStruct
	}
	fmt.Println("TestAssetfileReg...")

	ts := TestStruct{
		AssetId: "98e0b84063b311e8a5e3d1b3c579b67f",
		V2: TestSubStruct{
			UserName:    "btd121",
			AssetName:   "assetnametest",
			AssetType:   12,
			FeatureTag:  "FeatureTag",
			SampleHash:  "43162f44ff5565b317ef3904c93ca525f0d739a86823bb5c1d08dbfcabbcde8c",
			StorageHash: "43162f44ff5565b317ef3904c93ca525f0d739a86823bb5c1d08dbfcabbcde8c",
			ExpireTime:  1527478061,
			OpType:      1,
			Price:       999999900000000,
			Description: "destest",
		},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
	cc, _ := HexToBytes("dc0002da00206230356563613430363362363131653861313164623763303833663930643061dc000ada0003626f74da00046e616d65cf000000000000000eda0005312d312d31da0000da004066636336386466646632316639343432616134306361363062313262396639653332383239663332346566343532653730656533623434313465363164396434ce5b195680ce00000001cf00038d7e9ed09f00da0003313233")
	err = Unmarshal(cc, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestUserReg(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 string
	}
	fmt.Println("TestUserReg...")

	ts := TestStruct{
		V1: "didinfotest",
		V2: "userinfotest",
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestAssetReg(t *testing.T) {
	type TestSubStruct struct {
		V1  string
		V2  string
		V3  string
		V4  string
		V5  string
		V6  string
		V7  string
		V8  string
		V9  uint32
		V10 uint64
		V11 string
		V12 uint32
		V13 string
	}

	type TestStruct struct {
		V1 string
		V2 *TestSubStruct
	}
	fmt.Println("TestAssetReg...")

	ts := TestStruct{
		V1: "23456789012345678901",
		V2: &TestSubStruct{V1: "usernametest", V2: "assetname", V3: "assettypetest", V4: "tagtest", V5: "pathtest", V6: "hasttest", V7: "storepathtest", V8: "12345678901234567890", V9: 11, V10: 22, V11: "desctriptest", V12: 333, V13: "signtest"},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestDataReqReg(t *testing.T) {
	type TestSubStruct struct {
		V1  string
		V2  string
		V3  uint64
		V4  uint64
		V5  string
		V6  uint64
		V7  uint32
		V8  uint64
		V9  uint32
		V10 string
	}

	type TestStruct struct {
		V1 string
		V2 *TestSubStruct
	}
	fmt.Println("TestDataReqReg...")

	ts := TestStruct{
		V1: "12345678901234567890",
		V2: &TestSubStruct{V1: "usernametest", V2: "reqnametest", V3: 111, V4: 222, V5: "hasttest", V6: 222, V7: 2, V8: 333, V9: 444, V10: "desctriptest"},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestGoodsProReq(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 uint32
		V3 string
		V4 string
	}

	fmt.Println("TestGoodsProReq...")

	ts := TestStruct{V1: "usernametest", V2: 2, V3: "asset", V4: "goodsIdTest"}

	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestDataDeal_PreSaleReq(t *testing.T) {
	type TestSubStruct struct {
		V1 string
		V2 string
		V3 string
		V4 string
		V5 uint32
		V6 uint64
	}

	type TestStruct struct {
		V1 string
		V2 *TestSubStruct
	}
	fmt.Println("TestDataDeal_PreSaleReq...")

	ts := TestStruct{
		V1: "12345678901234567899",
		V2: &TestSubStruct{V1: "usernametest", V2: "assetidTest", V3: "requireId", V4: "consumerId", V5: 2, V6: 222},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestDataDeal_BuyAssetReq(t *testing.T) {
	type TestSubStruct struct {
		V1 string
		V2 string
		V3 uint64
	}

	type TestStruct struct {
		V1 string
		V2 *TestSubStruct
	}
	fmt.Println("TestDataDeal_BuyAssetReq...")

	ts := TestStruct{
		V1: "12345678901234567899",
		V2: &TestSubStruct{V1: "buyertest", V2: "23456789012345678901", V3: 1234},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
}
