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
 * file description:  define constants for this blockchain
 * @Author: May Luo
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */

package config

const DEFAULT_BLOCK_INTERVAL uint32 = 3
const BLOCKS_PER_ROUND uint32 = 19
const VOTED_DELEGATES_PER_ROUND uint32 = 18
const CONSENSUS_BLOCKS_PERCENT uint32 = 80
const MAX_DELEGATE_VOTES uint32 = 40
const DELEGATE_PATICIPATION uint64 = 33
const MAX_BLOCK_SIZE uint32 = 32000000 //2048000000
const DEFALT_SLOT_CHECK_INTERVAL = 500000

const BOTTOS_CONTRACT_NAME string = "bottos"
const BOTTOS_INIT_SUPPLY uint64 = 1000000000

const INIT_DELEGATE_NUM uint32 = 21
const DEFAULT_BLOCK_TIME_LIMIT uint64 = 200

const DEFAULT_MAX_LIFE_TIME uint64 = 600 //unit: second

const DEFAULT_MAX_PENDING_TRX_IN_POOL uint64 = 1000

const MAX_ACCOUNT_NAME_LENGTH int = 16

const DEFAULT_OPTIONDB_NAME string = "bottos"
const DEFAULT_OPTIONDB_TABLE_BLOCK_NAME string = "Blocks"
const DEFAULT_OPTIONDB_TABLE_TRX_NAME string = "Transactions"
const DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME string = "Accounts"
