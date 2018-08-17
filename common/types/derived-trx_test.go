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
	"testing"
	"github.com/bottos-project/bottos/bpl"
	"github.com/stretchr/testify/assert"
)

func newTrx(contract string, method string, param []byte) *Transaction {
	trx := &Transaction{
		Version: 1,
		CursorNum: 1,
		CursorLabel: 1,
		Lifetime: 1,
		Sender:   contract,
		Contract: contract,
		Method:   method,
		Param:    param,
		SigAlg: 1,
		Signature: []byte{},
	}

	return trx
}

func TestDerivedTransaction_onelayer(t *testing.T) {
	trx := newTrx("test", "testmethod", []byte{})

	trx1 := &DerivedTransaction{Transaction: newTrx("subtest1", "testmethod", []byte{}), DerivedTrx:make([]*DerivedTransaction,0)}
	trx2 := &DerivedTransaction{Transaction: newTrx("subtest2", "testmethod", []byte{}), DerivedTrx:make([]*DerivedTransaction,0)}
	trx3 := &DerivedTransaction{Transaction: newTrx("subtest3", "testmethod", []byte{}), DerivedTrx:make([]*DerivedTransaction,0)}
	var derivedTrx []*DerivedTransaction
	derivedTrx = append(derivedTrx, trx1)
	derivedTrx = append(derivedTrx, trx2)
	derivedTrx = append(derivedTrx, trx3)
	handledTrx := &HandledTransaction{Transaction: trx, DerivedTrx: derivedTrx}

	data, err := bpl.Marshal(handledTrx)
	assert.Nil(t, err)

	unmTrx := &HandledTransaction{}
	err = bpl.Unmarshal(data, unmTrx)
	assert.Nil(t, err)
	assert.Equal(t, handledTrx, unmTrx)
}

func TestDerivedTransaction_twolayer(t *testing.T) {
	trx := newTrx("test", "testmethod", []byte{})

	trx11 := &DerivedTransaction{Transaction: newTrx("subtest11", "testmethod", []byte{5, 6, 7}), DerivedTrx:make([]*DerivedTransaction,0)}
	trx12 := &DerivedTransaction{Transaction: newTrx("subtest12", "testmethod", []byte{6, 7, 8}), DerivedTrx:make([]*DerivedTransaction,0)}
	trx13 := &DerivedTransaction{Transaction: newTrx("subtest13", "testmethod", []byte{7, 8, 9}), DerivedTrx:make([]*DerivedTransaction,0)}

	var derivedTrx1 []*DerivedTransaction
	derivedTrx1 = append(derivedTrx1, trx11)

	var derivedTrx2 []*DerivedTransaction
	derivedTrx2 = append(derivedTrx2, trx12)
	derivedTrx2 = append(derivedTrx2, trx13)

	trx1 := &DerivedTransaction{Transaction: newTrx("subtest1", "testmethod", []byte{1, 2, 3}), DerivedTrx: derivedTrx1}
	trx2 := &DerivedTransaction{Transaction: newTrx("subtest2", "testmethod", []byte{2, 3, 4}), DerivedTrx: derivedTrx2}
	trx3 := &DerivedTransaction{Transaction: newTrx("subtest3", "testmethod", []byte{3, 4, 5}), DerivedTrx:make([]*DerivedTransaction,0)}

	var derivedTrx []*DerivedTransaction
	derivedTrx = append(derivedTrx, trx1)
	derivedTrx = append(derivedTrx, trx2)
	derivedTrx = append(derivedTrx, trx3)
	handledTrx := &HandledTransaction{Transaction: trx, DerivedTrx: derivedTrx}

	data, err := bpl.Marshal(handledTrx)
	assert.Nil(t, err)

	unmHandledTrx := &HandledTransaction{}
	err = bpl.Unmarshal(data, unmHandledTrx)
	assert.Nil(t, err)
	assert.Equal(t, handledTrx, unmHandledTrx)
}
