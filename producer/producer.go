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
	"time"

	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/consensus/dpos"
)

type Producer struct {
	core *chain.BlockChain
}

func New(b *chain.BlockChain) *Producer {
	return &Producer{b}
}
func (p *Producer) isEligible() bool {
	return true
}
func (p *Producer) isReady() bool {
	return true
	slotTime := dpos.GetSlotTime(1)
	fmt.Println(slotTime)
	if slotTime >= common.NowToSeconds() {
		return true
	}
	return false
}
func (p *Producer) isMyTurn() bool {
	return true

}
func (p *Producer) Woker() *types.Block {

	if p.isEligible() && p.isReady() && p.isMyTurn() {
		now := time.Now()
		slot := dpos.GetSlotAtTime(now)
		scheduledTime := dpos.GetSlotTime(slot)
		fmt.Println(scheduledTime)
		block, err := p.reportBlock()
		if err != nil {
			return nil // errors.New("report Block failed")
		}
		return block
		fmt.Println("brocasting block", block)
	}
	return nil
}

//func reportBlock(reportTime time.Time, reportor role.Delegate) *types.Block {
func (p *Producer) reportBlock() (*types.Block, error) {
	chain := chain.GetChain()
	head := types.NewHeader()
	head.PrevBlockHash = chain.HeadBlockHash().Bytes()
	head.Number = chain.HeadBlockNum() + 1
	head.Timestamp = chain.HeadBlockTime() + uint64(config.DEFAULT_BLOCK_INTERVAL)
	head.Producer = []byte("my")
	block := types.NewBlock(head, nil)
	block.Header.ProducerSign = block.Sign("123").Bytes()
	return block, nil

}
