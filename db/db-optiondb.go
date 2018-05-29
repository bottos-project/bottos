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
 * file description: code database interface
 * @Author: May Luo
 * @Date:   2017-12-05
 * @Last Modified by:
 * @Last Modified time:
 */

package db

import (
	"errors"
	"fmt"
)

func (d *DBService) Insert(collection string, value interface{}) error {
	if d.optDbRepo == nil {
		//fmt.Println("error optiondb is not init")
		return nil
	}
	return d.optDbRepo.InsertOptionDb(collection, value)
}
func (d *DBService) Find(collection string, key string, value interface{}) (interface{}, error) {
	if d.optDbRepo == nil {
		//fmt.Println("error optiondb is not init")
		return nil, errors.New("error optiondb is not init")
	}
	return d.optDbRepo.OptionDbFind(collection, key, value)
}

func (d *DBService) Update(collection string, key string, value interface{}, updatekey string, updatevalue interface{}) error {
    if d.optDbRepo == nil {
        //fmt.Println("error optiondb is not init")
        return errors.New("error optiondb is not init")
    }
    
    return d.optDbRepo.OptionDbUpdate(collection, key, value, updatekey, updatevalue)
}
