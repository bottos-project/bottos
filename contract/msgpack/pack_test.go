// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  msgpack go
 * @Author: Gong Zibin
 * @Date:   2017-12-20
 * @Last Modified by:
 * @Last Modified time:
 */
package msgpack

import (
	"fmt"
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
	fmt.Println(err)
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
	fmt.Println(err)
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
	fmt.Println(err)
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
		V1: "buyertest",
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
		Signature   string
	}

	type TestStruct struct {
		V1 string
		V2 TestSubStruct
	}
	fmt.Println("TestAssetReg...")

	ts := TestStruct{
		V1: "23456789012345678901",
		V2: TestSubStruct{UserName: "usernametest", AssetName: "assetname",
			AssetType: 11, FeatureTag: "tagtest",
			SampleHash: "hasttest", StorageHash: "12345678901234567890",
			ExpireTime: 1521990, OpType: 22,
			Price: 1, Description: "desctriptest",
			Signature: "signtest"},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1, err)
	bb, _ := HexToBytes("dc0002da00203531653864633130373735393131653862343236643930383735383530383438dc000ada00056565656565da000474657374cf000000000000000bda0005312d312d31da0000da004033333236366339326365643038613861393634326134363465373663643230613537386233613536303435663436376539306639373439633533636532336265ce00000000ce00000001cf000000003b9aca00da0003313233")
	//bb, _ := HexToBytes("dc0002da00203632366431613530373734623131653838306536316633313030663634656561dc000ada0007626f74746f7331da00106865616c746820636172652064617461cf000000000000000bda000869642d6e616d652dda004034336363343433323962626331323238626331383631373431663064656535336264663064326338336666323736666139326132656137366261366332383665da004037336561626366363337643137633731626334356133303266326565633839636264383632333737303833303561313861376433656638376537383137326131ce00000000ce00000001cf000000000bebc200da001074657374206865616c74682064617461")
	err = Unmarshal(bb, &ts1)
	fmt.Println("ts1 ", ts1, err)
}

func TestDataReqReg(t *testing.T) {
	type TestSubStruct struct {
		V1 string
		V2 string
		V3 uint64
		V4 uint64
		V5 string
		V6 uint64
		V7 uint32
		V8 uint64
		V9 string

		V10 uint32
		V11 string
	}

	type TestStruct struct {
		V1 string
		V2 *TestSubStruct
	}
	fmt.Println("TestDataReqReg...")

	ts := TestStruct{
		V1: "12345678901234567890",
		V2: &TestSubStruct{V1: "usernametest", V2: "reqnametest", V3: 111, V4: 222, V5: "hasttest",
			V6: 222, V7: 2, V8: 333, V9: "bto", V10: 444, V11: "desctriptest"},
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




func TestNodeManStruct(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 string
	}

	ts := TestStruct{
		V1: "1234",
		V2: "12345678",
	}
	fmt.Println("TestNodeManStruct...")
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1)
	fmt.Println(err)
}





func TestCreateTokenStruct(t *testing.T) {
	type TestStruct struct {
		V1 string		
		V2 uint64
	}

	ts := TestStruct{
		V1: "DTO",		
		V2: 100000000000000000,
	}
	fmt.Println("TestCreateTokenStruct...")
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1)
	fmt.Println(err)
}



func TestIssueTokenStruct(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 string
		V3 uint64
	}

	ts := TestStruct{
		V1: "DTO",
		V2: "testfrom",
		V3: 100000000000000000,
	}
	fmt.Println("TestIssueTokenStruct...")
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1)
	fmt.Println(err)
}



func TestTransferTokenStruct(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 string
		V3 string
		V4 uint64
	}

	ts := TestStruct{
		V1: "testfrom",
		V2: "testto",
		V3: "DTO",
		V4:  1,
	}
	fmt.Println("TestTransferTokenStruct...")
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1)
	fmt.Println(err)
}



func TestTransferCreditStruct(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 string
		V3 string
		V4 uint64
	}

	ts := TestStruct{
		V1: "testfrom",
		V2: "datadealmng",
		V3: "DTO",
		V4: 100,
	}
	fmt.Println("TestTransferCreditStruct...")
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1)
	fmt.Println(err)
}

func TestDeleteCreditStruct(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 string
		V3 string
	}

	ts := TestStruct{
		V1: "testfrom",
		V2: "datadealmng",
		V3: "DTO",
	}
	fmt.Println("TestDeleteCreditStruct...")
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Println("ts1 ", ts1)
	fmt.Println(err)
}
