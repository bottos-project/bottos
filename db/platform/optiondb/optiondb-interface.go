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
 * file description: option database interface
 * @Author: May Luo
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */
package optiondb

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
)

type OptionDbRepo interface {
	InsertOptionDb(collection string, value interface{}) error
}

func (r *OptionDbRepository) InsertOptionDb(collection string, value interface{}) error {
	fmt.Println("ddddddddddddddddddddddddddddddddddddddd")
	fmt.Println(value)
	session, err := GetSession(r.mgoEndpoint)
	if err != nil {
		fmt.Println(err)
		return errors.New("Get session faild" + r.mgoEndpoint)
	}

	insert := func(c *mgo.Collection) error {

		return c.Insert(value)
	}
	session.SetCollection(collection, insert)

	return nil
}
