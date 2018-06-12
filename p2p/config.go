// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

// This program is free software: you can distribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Bottos.  If not, see <http://www.gnu.org/licenses/>.

// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package exec provides functions for executing WebAssembly bytecode.

/*
 * file description: the interface for WASM execution
 * @Author: Stewart Li
 * @Date:   2018-02-08
 * @Last Modified by:
 * @Last Modified time:
 */

package p2pserver

const (
	//CONF_FILE is definition of config file name
	CONF_FILE = "config.json"
	//TIME_INTERVAL is definition of time interval
	TIME_INTERVAL = 10
	//TST *WRAN* set the variable as "true" before starting test
	TST            = 0
	//MIN_NODE_NUM min peer before sync
	MIN_NODE_NUM = 2
	//INIT_SYNC_WAIT wait time before sync when startup
	INIT_SYNC_WAIT = 30

	//TIME_PNE_START start wait time , Minute
	TIME_PNE_START = 2
	//TIME_PNE_EXCHANGE exchange time , Minute
	TIME_PNE_EXCHANGE = 1
	//MAX_NEIGHBOR_NUM  max neighbor number
	MAX_NEIGHBOR_NUM = 200
	//NEIGHBOR_DISCOVER_COUNT  bunch of neighbor when discover at the same time
	NEIGHBOR_DISCOVER_COUNT = 10
)
