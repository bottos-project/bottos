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
 * file description: producer
 * @Author: May Luo
 * @Date:   2017-12-11
 * @Last Modified by:
 * @Last Modified time:
 */
package producer

import (
	"fmt"
	"testing"
	"time"

	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/db"
)

func startup() *Reporter {
	dbInst := db.NewDbService("./temp/db", "./temp/codedb")
	if dbInst == nil {
		fmt.Println("Create DB service fail")
	}
	bc, err := chain.CreateBlockChain(dbInst)
	if err != nil {
		fmt.Println("Create DB service fail")
	}
	reportIns := &Reporter{false, bc, dbInst}
	return reportIns
}
func tearDown(r *Reporter) {
	r.db.Close()
}
func TestReporter_GetSlotAtTime(t *testing.T) {
	ins := startup()
	cbegin := time.Time{}
	slot := ins.GetSlotAtTime(cbegin)
	fmt.Println(slot)
	cUnix := cbegin.Unix()
	fmt.Println(cUnix)
	//	slot = ins.GetSlotAtTime(cUnix)
	//	fmt.Println(slot)
	now := common.NowToSeconds(time.Now().Unix())
	slot = ins.GetSlotAtTime(now)
	fmt.Println(slot)

	nowMicroSec := common.NowToSlotSec(time.Now(), 500000)
	slot = ins.GetSlotAtTime(nowMicroSec)
	fmt.Println(slot)

}
