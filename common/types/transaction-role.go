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
 * file description:  transaction
 * @Author: Gong Zibin
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */
package types

import (
	"crypto/sha256"
	"github.com/bottos-project/bottos/common"
	"github.com/golang/protobuf/proto"
)

func (trx *Transaction) Hash() common.Hash {
	data, _ := proto.Marshal(trx)
	temp := sha256.Sum256(data)
	hash := sha256.Sum256(temp[:])
	return hash
}

func (trx *Transaction) ValidateSign() bool {
	return true
}
