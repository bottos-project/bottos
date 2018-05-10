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
	"fmt"
	"math/big"
	"testing"

	"github.com/bottos-project/core/db"
)

func TestDelegateVotes_writedb(t *testing.T) {
	ins := db.NewDbService("./file2", "./file2/db.db", "")
	err := CreateDelegateVotesRole(ins)
	if err != nil {
		fmt.Println(err)
	}
	value := &DelegateVotes{
		OwnerAccount: "nodepad",
		Serve: Serve{
			Votes:          1,
			Position:       big.NewInt(2),
			TermUpdateTime: big.NewInt(2),
			TermFinishTime: big.NewInt(2),
		},
	}
	err = SetDelegateVotesRole(ins, value.OwnerAccount, value)
	if err != nil {
		fmt.Println("SetDelegateVotesRole", err)
	}

	value, err = GetDelegateVotesRoleByAccountName(ins, value.OwnerAccount)
	if err != nil {
		fmt.Println("GetDelegateVotesRoleByAccountName", err)
	}
	fmt.Println(value)

	value, err = GetDelegateVotesRoleByVote(ins, value.Serve.Votes)
	if err != nil {
		fmt.Println("GetDelegateVotesRoleByVote", err)
	}
	fmt.Println(value)

	value, err = GetDelegateVotesRoleByFinishTime(ins, value.Serve.TermFinishTime)
	if err != nil {
		fmt.Println("GetDelegateVotesRoleByFinishTime", err)
	}
	fmt.Println(value)

	values, nerr := GetAllDelegateVotesRole(ins)
	if nerr != nil {
		fmt.Println("GetAllDelegateVotes", nerr)
	}
	fmt.Println(len(values))

	svotes, nerr := GetAllSortVotesDelegates(ins)
	if nerr != nil {
		fmt.Println("GetAllSortVotesDelegates", nerr)
	}
	fmt.Println(len(svotes))
	tvotes, nerr := GetAllSortFinishTimeDelegates(ins)
	if nerr != nil {
		fmt.Println("GetAllSortFinishTimeDelegates", nerr)
	}
	fmt.Println(len(tvotes))
	db.Close()
}
