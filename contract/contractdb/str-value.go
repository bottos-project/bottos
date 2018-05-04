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
 * file description:  contract db
 * @Author: Gong Zibin
 * @Date:   2017-12-20
 * @Last Modified by:
 * @Last Modified time:
 */

package contractdb

import (
)

func (cdb *ContractDB) getStrValue(objectName string, key string) (string, error) {
	value, err := cdb.Db.GetObject(objectName, key)
	if err != nil {
		return "", err
	}

	return string(value), nil
}

func (cdb *ContractDB) setStrValue(objectName string, key string, value string) error {
	return cdb.Db.SetObject(objectName, key, value)
}

func (cdb *ContractDB) removeStrValue(objectName string, key string) error {
	// TODO
	return nil
}
  