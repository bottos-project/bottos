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
	"fmt"
	"time"

	"github.com/aristanetworks/goarista/monotime"
)

//current system time seconds
func Now() uint64 {
	return uint64(time.Now().Unix())
}
func NowToSeconds() uint64 {
	return uint64(time.Now().Unix())
}
func NowToMicroseconds() uint64 {
	return uint64(time.Now().UnixNano() / 1000)
}
func ToNanoseconds(current time.Time) uint64 {
	cur := current.UnixNano()
	fmt.Println(cur)
	return 0
}

//current monotonic clock time use to measure time
func MeasureStart() uint64 {
	return monotime.Now()
}
func Elapsed(t uint64) uint64 {

	elapse := MeasureStart() - t
	fmt.Println(elapse)
	return elapse

}

func MicrosecondsAddToSec(src uint64, des uint64) uint64 {
	addNew := src + des
	return uint64(addNew / 1000000)
}

func NowToSlotSec(current time.Time, loopMicroSec uint64) uint64 {
	cur := ToMicroseconds(current)
	value := MicrosecondsAddToSec(cur, loopMicroSec)
	fmt.Println(cur)
	return value
}
