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
 * file description:  context definition
 * @Author: Gong Zibin
 * @Date:   2017-01-15
 * @Last Modified by:
 * @Last Modified time:
 */

package contract

import (
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/role"
)

//Context for contracts
type Context struct {
	RoleIntf   role.RoleInterface
	Trx        *types.Transaction
}

//GetTrxParam for contracts
func (ctx *Context) GetTrxParam() []byte {
	return ctx.Trx.Param
}

//GetTrxParamSize for contracts
func (ctx *Context) GetTrxParamSize() uint32 {
	size := len(ctx.Trx.Param)
	return uint32(size)
}
