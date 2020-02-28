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
	"math/big"
	"testing"
	"fmt"
	log "github.com/cihub/seelog"

	"github.com/bottos-project/bottos/db"
)

func TestTransitVotes_writedb(t *testing.T) {
	log.Info("gtest")
	ins := db.NewDbService("./vote1", "vote2")
	err := CreateTransitVotesRole(ins)
	if err != nil {
		log.Error(err)
	}
	value := &TransitVotes{
		ProducerAccount: "nodepad",
		TransitVotes:           big.NewInt(10000000),
	}
	err = SetTransitVotesRole(ins, value.ProducerAccount, value)
	if err != nil {
		log.Error("SetTransitVotesRole", err)
	}
	value2 := &TransitVotes{
		ProducerAccount: "nodepad111",
		TransitVotes:           big.NewInt(20000000),
	}
	err = SetTransitVotesRole(ins, value2.ProducerAccount, value2)
	if err != nil {
		log.Error("SetTransitVotesRole", err)
	}
	value3 := &TransitVotes{
		ProducerAccount: "nodepad222",
		TransitVotes:           big.NewInt(20000000),
	}
	err = SetTransitVotesRole(ins, value3.ProducerAccount, value3)
	if err != nil {
		log.Error("SetTransitVotesRole", err)
	}

	myValue, err1 := GetTransitVotesRole(ins, value.ProducerAccount)
	if err1 != nil {
		log.Error("GetDelegateVotesRoleByFinishTime", err1)
	}
	log.Info(myValue)

	svots, nerr := GetAllSortTransitVotesDelegates(ins)
	if nerr != nil {
		log.Error("GetAllSortTransitVotesDelegates", nerr)
	}
	log.Info(len(svots))
	fmt.Println(len(svots))

	schedule := ElectTransitPeriodDelegatesRole(ins, false)
	log.Info(len(schedule))
	log.Info(schedule)

}
