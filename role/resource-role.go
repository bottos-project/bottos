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
 * file description:  resource role
 * @Author: leo
 * @Date:   2018-11-1
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"

	"github.com/bottos-project/bottos/db"
	log "github.com/cihub/seelog"
)

// ResourceObjectName is definition of object name of balance
const ResourceLimitObjectName string = "resource_limit"

// StakedBalanceObjectName is definition of object name of stake balance
const ResourceUsageObjectName string = "resource_usage"

// ResourceLimit 单账户在窗口期内的可使用总量最大值，每次触发执行合约时更新，更新时必须保证用户使用量不超过此上限
type ResourceLimit struct {
	AccountName                string `json:"account_name"`
	PledgedSpaceLimitInWin     uint64 `json:"pledged_space_limit_in_win"`
	PledgedTimeTokenLimitInWin uint64 `json:"pledged_time_token_limit_in_win"`
	FreeTimeTokenLimitInWin    uint64 `json:"free_time_token_limit_in_win"`
	FreeSpaceTokenLimitInWin   uint64 `json:"free_space_token_limit_in_win"`
}
// GetResourceLimitRole is to get resource limit
func GetResourceLimitRole(ldb *db.DBService, accountName string) (*ResourceLimit, error) {
	key := accountName
	value, err := ldb.GetObject(ResourceLimitObjectName, key)
	if err != nil {
		return nil, err
	}

	res := &ResourceLimit{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		log.Error("ROLE Unmarshal failed accountName ", accountName)
		return nil, err
	}

	return res, nil
}
// GetResourceUsageRoleByName is to get resource usage
func GetResourceUsageRoleByName(ldb *db.DBService, name string) (*ResourceUsage, error) {
	key := name
	value, err := ldb.GetObject(ResourceUsageObjectName, key)
	if err != nil {
		log.Errorf("DB: get resource role:%v", err)
		return nil, err
	}

	res := &ResourceUsage{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		log.Errorf("ROLE Unmarshal failed %v", err)
		return nil, err
	}

	return res, nil
}

// SetResourceLimitRole is to set resource limit
/*func SetResourceReceiptRole(ldb *db.DBService, accountName string, value *ResourceReceipt) error {
	key := accountName
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ldb.SetObject(ResourceLimitObjectName, key, string(jsonvalue))
}*/

// GetResourceLimitRole is to get resource limit
/*func GetResourceReceiptRole(ldb *db.DBService, accountName string) (*ResourceReceipt, error) {
	key := accountName
	value, err := ldb.GetObject(ResourceLimitObjectName, key)
	if err != nil {
		return nil, err
	}

	res := &ResourceReceipt{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}*/
