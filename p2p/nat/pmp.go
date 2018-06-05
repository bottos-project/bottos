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
	"fmt"
	"github.com/jackpal/gateway"
	pmp "github.com/jackpal/go-nat-pmp"
	"net"
)

type client struct {
	gw  net.IP
	pmp *pmp.Client
}

/*mapping,  time is second */
func (c *client) Mapping(p string, internalPort int, externalPort int, time int) (err error) {
	if externalPort == 0 && time != 0 ||
		externalPort != 0 && time == 0 {
		return fmt.Errorf("error param")
	}
	_, err = c.pmp.AddPortMapping(p, internalPort, externalPort, time)
	return err
}

func (c *client) EIP() (net.IP, error) {
	r, err := c.pmp.GetExternalAddress()
	if err != nil {
		return nil, err
	}
	return r.ExternalIPAddress[:], nil
}

func (c *client) String() string {
	return fmt.Sprintf("nat pmp(%v)", c.gw)
}

func getPmpClient() Handle {
	gw, err := gateway.DiscoverGateway()
	if err != nil {
		return nil
	}

	c := pmp.NewClient(gw)
	_, err = c.GetExternalAddress()
	if err != nil {
		return nil
	}

	return &client{gw, c}

}
