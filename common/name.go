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
 * file description:  encoded Name
 * @Author: Gong Zibin
 * @Date:   2018-08-08
 * @Last Modified by:
 * @Last Modified time:
 */

package common

import (
	"regexp"
	"strings"
	//log "github.com/cihub/seelog"
	)

const (
	// MaxNameLength define max account name length
	MaxNameLength int = 21
)

type NameType uint

const (
	NameTypeAccount	    NameType = 1
	NameTypeExContract  NameType = 2
	NameTypeUnknown     NameType = 3
)

// ACCOUNT_NAME_REGEXP define account name format
const CONTRACT_NAME_REGEXP string = "^[a-z][a-z0-9]{2,9}$"
const ACCOUNT_NAME_REGEXP string = "^[a-z][a-z0-9.-]{2,20}$"
const EX_CONTRACT_NAME_REGEXP string = "^[a-z][a-z0-9]{2,9}@[a-z][a-z0-9.-]{2,20}$"
      

var AccountReg *regexp.Regexp = regexp.MustCompile(ACCOUNT_NAME_REGEXP)
var ContractReg *regexp.Regexp = regexp.MustCompile(CONTRACT_NAME_REGEXP)
var ExContractReg *regexp.Regexp = regexp.MustCompile(EX_CONTRACT_NAME_REGEXP)
//var CoreLogger log.LoggerInterface

func CheckAccountNameContent(name string) bool {
	return AccountReg.MatchString(name)
}

func CheckContractNameContent(name string) bool {
	return ContractReg.MatchString(name) 
}

func CheckExContractNameContent(name string) bool {
	return ExContractReg.MatchString(name)
}

func AnalyzeName(name string) (NameType, string) {
	
	if CheckAccountNameContent(name) {
		return NameTypeAccount, name
	} else if CheckExContractNameContent(name) {
		separateSymbol := strings.Index(name, "@")
		return NameTypeExContract,  name[separateSymbol+1:]
	} else {
		return NameTypeUnknown, ""
	}
	
	// else if  nc.checkContractNameContent(name) {
	// 	return NameTypeContract, ""
	// }
	
}

/*func Printf(format string, msg ...interface{}) {
        CoreLogger.Infof(format, msg...)
	//CoreLogger.Flush()
}

func Println(format string, msg ...interface{}) {
        CoreLogger.Infof(format, msg...)
	//CoreLogger.Flush()
}

func Errorf(format string, msg ...interface{}) error {
        err := CoreLogger.Errorf(format, msg...)
	//CoreLogger.Flush()
	return err
}*/
