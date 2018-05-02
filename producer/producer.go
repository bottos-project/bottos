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

	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/role"
)

type Reporter struct {
	core     chain.BlockChainInterface
	roleIntf role.RoleInterface
	state    ReportState
}
type ReporterRepo interface {
	Woker(Trxs []*types.Transaction) *types.Block
	VerifyTrxs(Trxs []*types.Transaction) error
	IsReady() bool
}

func New(b chain.BlockChainInterface, roleIntf role.RoleInterface) ReporterRepo {
	stat := ReportState{false, 0, false}
	return &Reporter{core: b, roleIntf: roleIntf, state: stat}
}

func (p *Reporter) Woker(trxs []*types.Transaction) *types.Block {

	now := common.NowToSeconds()
	slot := p.roleIntf.GetSlotAtTime(now)
	scheduledTime := p.roleIntf.GetSlotTime(slot)
	fmt.Println("Woker", scheduledTime, slot)
	accountName, err1 := p.roleIntf.GetScheduleDelegateRole(slot)
	if err1 != nil {
		return nil // errors.New("report Block failed")
	}

	block, err := p.reportBlock(accountName, trxs)
	if err != nil {
		return nil // errors.New("report Block failed")
	}

	fmt.Println("brocasting block", block)
	return block
}

func (p *Reporter) StartTag() error {
	//p.core.

	return nil

}
func (p *Reporter) VerifyTrxs(trxs []*types.Transaction) error {

	return nil
}

//func reportBlock(reportTime time.Time, reportor role.Delegate) *types.Block {
func (p *Reporter) reportBlock(accountName string, trxs []*types.Transaction) (*types.Block, error) {
	head := types.NewHeader()
	head.PrevBlockHash = p.core.HeadBlockHash().Bytes()
	head.Number = p.core.HeadBlockNum() + 1
	head.Timestamp = p.core.HeadBlockTime() + uint64(config.DEFAULT_BLOCK_INTERVAL)
	head.Delegate = []byte(accountName)
	block := types.NewBlock(head, trxs)
	block.Header.DelegateSign = block.Sign("123").Bytes()

	return block, nil
}
