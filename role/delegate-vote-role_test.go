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
 * file description: code database test
 * @Author: May Luo
 * @Date:   2017-12-04
 * @Last Modified by:
 * @Last Modified time:
 */
package role

import (
	//	"encoding/json"
	//"math/big"
	"testing"

	log "github.com/cihub/seelog"

	//"github.com/bottos-project/bottos/db"
)

func TestDelegateVotes_writedb(t *testing.T) {
	log.Info("gtest")
	//	ins := db.NewDbService("./file2", "./file2/db.db")
	//	err := CreateDelegateVotesRole(ins)
	//	if err != nil {
	//		log.Error(err)
	//	}
	//	value := &DelegateVotes{
	//		OwnerAccount: "nodepad",
	//		Serve: Serve{
	//			Votes:          big.NewInt(2),
	//			Position:       big.NewInt(2),
	//			TermUpdateTime: big.NewInt(2),
	//			TermFinishTime: big.NewInt(2),
	//		},
	//	}
	//	err = SetDelegateVotesRole(ins, value.OwnerAccount, value)
	//	if err != nil {
	//		log.Error("SetDelegateVotesRole", err)
	//	}

	//	value, err = GetDelegateVotesRoleByFinishTime(ins, value.Serve.TermFinishTime)
	//	if err != nil {
	//		log.Error("GetDelegateVotesRoleByFinishTime", err)
	//	}
	//	log.Info(value)

	//	values, nerr := GetAllDelegateVotesRole(ins)
	//	if nerr != nil {
	//		log.Error("GetAllDelegateVotes", nerr)
	//	}
	//	log.Info(len(values))

	//	svotes, nerr := GetAllSortVotesDelegates(ins)
	//	if nerr != nil {
	//		log.Error("GetAllSortVotesDelegates", nerr)
	//	}
	//	log.Info(len(svotes))
	//	tvotes, nerr := GetAllSortFinishTimeDelegates(ins)
	//	if nerr != nil {
	//		log.Error("GetAllSortFinishTimeDelegates", nerr)
	//	}
	//	log.Info(len(tvotes))
}
