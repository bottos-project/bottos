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
 * file description:  file introduction for commom tip
 * @Author: May Luo
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */

package common

import (
	"time"
)

type Microseconds struct {
	int64 micro
}

func NowToSeconds() int64 {
	return time.Now().Unix()
}
func microseconds() int64 {

}
func millseconds(s int64) int64 {
	return s * 1000
}
func seconds(s int64) int64 {
	return s * 1000000
}

func getEpochTime() int64 {
	now := time.Now() // get current time
	epoch := time.Since(now)
	return epoch.Nanoseconds * 1000 // microseconds
}
