package producer

import (
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

func (r *Reporter) IsMyTurn(slot uint32) bool {
	accountName, err := r.roleIntf.GetScheduleDelegateRole(slot)
	if err != nil {
		return false
	}
	if r.roleIntf.IsAccountExist(accountName) == true {
		return true
	}
	return false
}
