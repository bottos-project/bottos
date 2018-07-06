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
	"errors"
	log "github.com/cihub/seelog"
	"math/rand"
	"reflect"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
)

//StringSliceReflectEqual is reflecting equal to two slices
func StringSliceReflectEqual(a, b []string) bool {
	return reflect.DeepEqual(a, b)
}

//ShuffleEelectCandidateList is to shuffle the candidates in one round
func (r *Reporter) ShuffleEelectCandidateList(block types.Block) ([]string, error) {
	newSchedule := r.roleIntf.ElectNextTermDelegates(&block, true)
	currentState, err := r.roleIntf.GetCoreState()
	if err != nil {
		return nil, err
	}
	changes := common.Filter(currentState.CurrentDelegates, newSchedule)
	equal := reflect.DeepEqual(block.Header.DelegateChanges, changes)
	if equal == false {
		log.Info("invalid block changes")
		return nil, errors.New("Unexpected round changes in new block header")
	}

	h := block.Hash()
	label := h.Label()
	rand.New(rand.NewSource(int64(label)))
	rand.Shuffle(len(newSchedule), func(i, j int) {
		newSchedule[i], newSchedule[j] = newSchedule[j], newSchedule[i]
	})

	return newSchedule, nil
}
