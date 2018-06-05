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
 * file description:  time point
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

// Now current system time seconds
func Now() uint64 {
	return uint64(time.Now().Unix())
}

// NowToSeconds return seconds of now
func NowToSeconds() uint64 {
	return uint64(time.Now().UnixNano() / 1000 / 1000000)
}

// NowToMicroseconds return microseconds of now
func NowToMicroseconds() uint64 {
	return uint64(time.Now().UnixNano() / 1000)
}

// ToNanoseconds return nanoseconds of now
func ToNanoseconds(current time.Time) uint64 {
	cur := current.UnixNano()
	return uint64(cur)
}

// MeasureStart current monotonic clock time use to measure time
func MeasureStart() uint64 {
	return monotime.Now()
}

// Elapsed calculate elapsed microseconds
func Elapsed(t uint64) uint64 {
	nanoElapse := MeasureStart() - t
	microElapse := nanoElapse / 1000
	//	fmt.Println(microElapse)
	return microElapse

}

// NanoToMicroSec convert nanosecond to microsecond
func NanoToMicroSec(src uint64) uint64 {
	return uint64(src / 1000)
}

// MicrosecondsAddToSec add seconds
func MicrosecondsAddToSec(src uint64, des uint64) uint64 {
	addNew := src + des
	return uint64(addNew / 1000000)
}

// NowToSlotSec convert now seconds to slot second
func NowToSlotSec(current time.Time, loopMicroSec uint64) uint64 {
	cur := ToMicroseconds(current)
	value := MicrosecondsAddToSec(cur, loopMicroSec)
	fmt.Println(cur)
	return value
}
