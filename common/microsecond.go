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
 * file description:  provide a interface such as time to seconds and epoch time etc.
 * @Author: May Luo
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */
package common

import (
	"time"
)

type Microsecond struct {
	count uint64
}

func maximum() Microsecond {
	return Microsecond{count: 0x7fffffffffffffff}
}

func ToMicroseconds(current time.Time) uint64 {
	cur := current.Unix()
	return uint64(cur * 1000)
}
func ToSeconds(m Microsecond) uint64 {
	return m.count / 1000000
}
func ToMilliseconds(m Microsecond) uint64 {
	return m.count / 1000
}
func SecondsToMicro(s uint64) Microsecond {
	return Microsecond{s * 1000000}
}
func MilliSecToMicro(m uint64) Microsecond {
	return Microsecond{m * 1000}
}
