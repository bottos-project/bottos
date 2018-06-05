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
 * file description:  actor entry
 * @Author:
 * @Date:   2017-12-20
 * @Last Modified by:
 * @Last Modified time:
 */

package env

import (
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/chain/extra"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/contract/contractdb"
	"github.com/bottos-project/bottos/role"
)

//ActorEnv actor external interface
type ActorEnv struct {
	RoleIntf   role.RoleInterface
	ContractDB *contractdb.ContractDB
	Chain      chain.BlockChainInterface
	TxStore    *txstore.TransactionStore
	NcIntf     contract.NativeContractInterface
}
