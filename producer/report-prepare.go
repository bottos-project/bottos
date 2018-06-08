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
 * file description:  producer entry
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package producer

import (
	"fmt"
	"math"

	log "github.com/cihub/seelog"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/config"
)

//ReportState is recording the state of reporters
type ReportState struct {
	ScheduledTime     uint64
	ScheduledReporter string
	PrivateKey        string
	IsReporting       bool
	CheckFlag         uint32
	ReportEnable      bool
}

//IsReady is check if repoter state
func (r *Reporter) IsReady() bool {
	now := GetReportTimeNow()
	r.state.SetCheckFlag(1)

	slot := r.roleIntf.GetSlotAtTime(now)
	if slot == 0 {
		//log.Info("slot is 0,not time yet")
		return false
	}

	object, err := r.roleIntf.GetChainState()
	if err != nil {
		return false
	}
	if (now < object.LastBlockTime+uint64(config.DEFAULT_BLOCK_INTERVAL)) && object.LastBlockNum != 0 {
		//log.Infof("time not ready", now, object.LastBlockTime, uint64(config.DEFAULT_BLOCK_INTERVAL))
		return false
	}
	if r.IsMyTurn(now, slot) == false {
		return false
	}
	return true

}

//GetReportTimeNow is to count reporter's time
func GetReportTimeNow() uint64 {
	systemNow := common.NowToMicroseconds()
	nowMicro := common.Microsecond{}
	nowMicro.Count = (systemNow + config.DEFALT_SLOT_CHECK_INTERVAL)
	now := common.ToSeconds(nowMicro)
	return now
}

//StartReport is to start
func (r *Reporter) StartReport() {
	r.state.IsReporting = true
}

//EndReport is to stop report
func (r *Reporter) EndReport() {
	r.state.IsReporting = false
}

//SetCheckFlag is to set check flags
func (r *ReportState) SetCheckFlag(flag uint32) {
	r.CheckFlag |= flag
}

//IsSynced is to check synced flags
func (r *Reporter) IsSynced(when uint64) bool {
	if r.state.ReportEnable == true {
		return true
	}
	r.roleIntf.GetSlotTime(1)
	if r.roleIntf.GetSlotTime(1) >= when {
		r.state.ReportEnable = true
		return true
	}
	return false
}

//IsMyTurn is to check if is my turn
func (r *Reporter) IsMyTurn(startTime uint64, slot uint64) bool {
	accountName, err := r.roleIntf.GetCandidateBySlot(slot)
	if err != nil {
		log.Infof("cannot get delegate by slot", slot)
		return false
	}
	if r.roleIntf.IsAccountExist(accountName) == false {
		log.Infof("account not exist", accountName)
		return false
	}

	scheduledTime := r.roleIntf.GetSlotTime(slot)

	delegate, err := r.roleIntf.GetDelegateByAccountName(accountName)
	if err != nil {
		log.Infof("find delegate by account failed", accountName)
		return false
	}

	found := false
	for _, v := range config.Param.Delegates {
		if accountName == v {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("current delegate: %v, not found in this node\n", accountName)
		return false
	}

	prate := r.roleIntf.GetDelegateParticipationRate()

	if prate < config.DELEGATE_PATICIPATION {
		//	log.Info("delegate paticipate rate is too low")
		return false
	}

	if math.Abs(float64(scheduledTime)-float64(startTime)) > 500 {
		//	log.Info("delegate  is too slow")
		return false
	}
	r.state.ScheduledTime = scheduledTime
	r.state.ScheduledReporter = accountName
	r.state.PrivateKey = delegate.ReportKey

	return true
}
