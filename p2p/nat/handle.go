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

/*
 * file description:  nat handle
 * @Author: eripi
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package nat

import (
	"net"
	"time"
)

// Handle is definition of interface
type Handle interface {
	Mapping(p string, internalPort int, externalPort int, time int) error
	EIP() (net.IP, error)
	String() string
}
// Interface is definition of net
type Interface interface {
	
	AddMapping(protocol string, extport, intport int, name string, lifetime time.Duration) error
	DeleteMapping(protocol string, extport, intport int) error
	ExternalIP() (net.IP, error)
	String() string
}
// GetHandle is to get handler func by config
func GetHandle(config string) Handle {
	switch config {
	case "pmp":
		return getPmpClient()
	case "upnp":
		return nil
	default:
		return nil
	}
}

// Map is to handle a struct from intport to extport
func Map(h Handle, c chan struct{}, p string, intport int, extport int) {
	update := time.NewTimer(10 * time.Minute)
	err := h.Mapping(p, intport, extport, 1200)
	if err != nil {
	}

	defer func() {
		update.Stop()
		h.Mapping(p, intport, 0, 0)
	}()

	for {
		select {
		case _, ok := <-c:
			if !ok {
				break
			}
		case <-update.C:
			err := h.Mapping(p, intport, extport, 1200)
			if err != nil {
			}

			update.Reset(10 * time.Minute)
		}
	}
}
