package producer

import (
	"fmt"
	"math"

	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/config"
)

type ReportState struct {
	IsReporting  bool
	CheckFlag    uint32
	ReportEnable bool
}

func GetReportTimeNow() uint64 {
	systemNow := common.NowToMicroseconds()
	now := systemNow + config.DEFALT_SLOT_CHECK_INTERVAL
	return now
}
func (r *Reporter) StartReport() {
	r.state.IsReporting = true
}
func (r *Reporter) EndReport() {
	r.state.IsReporting = false
}
func (r *ReportState) SetCheckFlag(flag uint32) {
	r.CheckFlag |= flag
}
func (r *Reporter) IsSynced(when uint64) bool {
	if r.state.ReportEnable == true {
		return true
	}
	if r.roleIntf.GetSlotTime(1) >= when {
		r.state.ReportEnable = true
		return true
	}
	return false
}
func (r *Reporter) IsOnMySlotTime(when uint64) bool {
	slot := r.roleIntf.GetSlotAtTime(when)
	if slot == 0 {
		return false
	}
	return true
}

func (r *Reporter) IsMyTurn(startTime uint64, slot uint32) bool {
	accountName, err := r.roleIntf.GetScheduleDelegateRole(slot)
	if err != nil {
		fmt.Println("cannot get delegate by slot", slot)
		return false
	}
	if r.roleIntf.IsAccountExist(accountName) == false {
		fmt.Println("account not exist", accountName)
		return false
	}

	scheduledTime := r.roleIntf.GetSlotTime(slot)
	delegate, err := r.roleIntf.GetDelegateByAccountName(accountName)
	if err != nil {
		fmt.Println("find delegate by account failed", accountName)
		return false
	}
	// TODO check   delegate.SigningKey
	fmt.Println(delegate.SigningKey)

	prate := r.roleIntf.GetDelegateParticipationRate()
	if prate < config.DELEGATE_PATICIPATION {
		fmt.Println("delegate paticipate rate is too low")
		return false
	}

	if math.Abs(float64(scheduledTime)-float64(startTime)) > 500 {
		fmt.Println("delegate  is too slow")
		return false
	}

	return true
}
