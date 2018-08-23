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
 * file description:  blockchain general interface and logic
 * @Author: Gong Zibin
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package chain

//ConfirmedNum confirm array
type ConfirmedNum []uint64

func (c ConfirmedNum) Len() int {
	return len(c)
}

func (c ConfirmedNum) Swap(i, j int) {
	if i >= 0 && i < c.Len() && j >= 0 && j < c.Len() {
		c[i], c[j] = c[j], c[i]
	}
}

func (c ConfirmedNum) Less(i, j int) bool {
	if i >= 0 && i < c.Len() && j >= 0 && j < c.Len() {
		return c[i] < c[j]
	}
	return false
}
