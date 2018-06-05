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
 * file description:  general Hash type
 * @Author: Gong Zibin
 * @Date:   2017-12-05
 * @Last Modified by:
 * @Last Modified time:
 */
package types

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
)

func newTrx(contract string, method string, param []byte) *Transaction {
	trx := &Transaction{
		Sender:   contract,
		Contract: contract,
		Method:   method,
		Param:    param,
	}

	return trx
}

func TestDerivedTransaction_onelayer(t *testing.T) {
	trx := newTrx("test", "testmethod", []byte{})

	trx1 := &DerivedTransaction{Transaction: newTrx("subtest1", "testmethod", []byte{})}
	trx2 := &DerivedTransaction{Transaction: newTrx("subtest2", "testmethod", []byte{})}
	trx3 := &DerivedTransaction{Transaction: newTrx("subtest3", "testmethod", []byte{})}
	var derivedTrx []*DerivedTransaction
	derivedTrx = append(derivedTrx, trx1)
	derivedTrx = append(derivedTrx, trx2)
	derivedTrx = append(derivedTrx, trx3)
	handledTrx := &HandledTransaction{Transaction: trx, DerivedTrx: derivedTrx}

	data, err := proto.Marshal(handledTrx)
	fmt.Println(handledTrx)
	fmt.Printf("data:%x\n", data)

	unmTrx := &HandledTransaction{}
	err = proto.Unmarshal(data, unmTrx)
	fmt.Println(err)
	fmt.Println(unmTrx)
}

func TestDerivedTransaction_twolayer(t *testing.T) {
	trx := newTrx("test", "testmethod", []byte{})

	trx11 := &DerivedTransaction{Transaction: newTrx("subtest11", "testmethod", []byte{5, 6, 7})}
	trx12 := &DerivedTransaction{Transaction: newTrx("subtest12", "testmethod", []byte{6, 7, 8})}
	trx13 := &DerivedTransaction{Transaction: newTrx("subtest13", "testmethod", []byte{7, 8, 9})}

	var derivedTrx1 []*DerivedTransaction
	derivedTrx1 = append(derivedTrx1, trx11)

	var derivedTrx2 []*DerivedTransaction
	derivedTrx2 = append(derivedTrx2, trx12)
	derivedTrx2 = append(derivedTrx2, trx13)

	trx1 := &DerivedTransaction{Transaction: newTrx("subtest1", "testmethod", []byte{1, 2, 3}), DerivedTrx: derivedTrx1}
	trx2 := &DerivedTransaction{Transaction: newTrx("subtest2", "testmethod", []byte{2, 3, 4}), DerivedTrx: derivedTrx2}
	trx3 := &DerivedTransaction{Transaction: newTrx("subtest3", "testmethod", []byte{3, 4, 5})}

	var derivedTrx []*DerivedTransaction
	derivedTrx = append(derivedTrx, trx1)
	derivedTrx = append(derivedTrx, trx2)
	derivedTrx = append(derivedTrx, trx3)
	handledTrx := &HandledTransaction{Transaction: trx, DerivedTrx: derivedTrx}

	data, err := proto.Marshal(handledTrx)
	fmt.Println(handledTrx)
	fmt.Printf("data:%x\n", data)

	unmHandledTrx := &HandledTransaction{}
	err = proto.Unmarshal(data, unmHandledTrx)
	fmt.Println(err)
	fmt.Println(unmHandledTrx.Transaction)
	fmt.Println(unmHandledTrx.DerivedTrx[0].Transaction)
	fmt.Println(unmHandledTrx.DerivedTrx[0].DerivedTrx)
	fmt.Println(unmHandledTrx.DerivedTrx[1].Transaction)
	fmt.Println(unmHandledTrx.DerivedTrx[1].DerivedTrx)
	fmt.Println(unmHandledTrx.DerivedTrx[2].Transaction)
	fmt.Println(unmHandledTrx.DerivedTrx[2].DerivedTrx)
}
