package version

import (
	"reflect"
	"strconv"
	"strings"
	"fmt"
	"github.com/bottos-project/bottos/common/types"
	berr "github.com/bottos-project/bottos/common/errors"
	log "github.com/cihub/seelog"
)

type VersionInfo struct {
	BlockNum uint64//new version start from BlockNum(inclusive)
	VersionString string//new version(string)
	VersionNumber uint32//new version(uint32)
}

var versions []VersionInfo = []VersionInfo{ //in reverse order, from high to low
	{BlockNum: 11456000, VersionString: "1.2.2"},
	{BlockNum: 4533000, VersionString: "1.2.1"},
	{BlockNum: 3110000, VersionString: "1.2.0"},
	{BlockNum: 880000, VersionString: "1.1.0"},
	{BlockNum:0, VersionString:"1.0.0"},
}

func Init() error {
	if len(versions) == 0 {
		return fmt.Errorf("VERSION can't found any version info")
	}

	versionTableString := ""
	maxVersionNumber := ^uint32(0)
	for i := 0; i < len(versions); i++ {
		if n := GetUintVersion(versions[i].VersionString); n == 0 {
			return fmt.Errorf("VERSION found invalid version, version %v", versions[i].VersionString)
		} else {
			versions[i].VersionNumber = n
			if versions[i].VersionNumber >= maxVersionNumber {
				return fmt.Errorf("VERSION the versions's order is wrong")
			}
			maxVersionNumber = versions[i].VersionNumber
			versionTableString += fmt.Sprintf("BlockNum:%v, VersionString:%v, VersionNumber:%v; ", versions[i].BlockNum, versions[i].VersionString, versions[i].VersionNumber)
		}
	}
	log.Errorf("VERSION version table: %v", versionTableString)
	return nil
}

func GetAppVersionString() string {
	return versions[0].VersionString
}

func GetAppVersionNum() uint32 {
	return versions[0].VersionNumber
}

func GetVersionByBlockNum(blockNum uint64) *VersionInfo {
	var info *VersionInfo = nil
	for i := 0; i < len(versions); i++ {
		if blockNum >= versions[i].BlockNum {
			info = &versions[i]
			break
		}
	}
	return info
}

func GetVersionNumByBlockNum(blockNum uint64) uint32 {
	if info := GetVersionByBlockNum(blockNum); info == nil {
		return ^uint32(0)
	} else {
		return info.VersionNumber
	}
}

func CheckBlock(block *types.Block, funcname string) berr.ErrCode {
	rightVersion := GetVersionByBlockNum(block.GetNumber());
	if rightVersion == nil {
		log.Errorf("VERSION handle %v can't found version for number in my version table, block.number %v, block.version %v", funcname, block.GetNumber(), GetStringVersion(block.GetVersion()))
		return berr.ErrBlockVersionError
	}

	if  block.GetVersion() != rightVersion.VersionNumber {
		log.Errorf("VERSION handle %v version match failed, block.number %v, block.version %v, right version %v", funcname, block.GetNumber(), GetStringVersion(block.GetVersion()), rightVersion.VersionString)
		return berr.ErrBlockVersionError
	}
	return berr.ErrNoError
}

func GetStringVersion(version uint32) string {
	s := strconv.Itoa((int(version)&0xff0000)>>16) + "." + strconv.Itoa((int(version)&0x00ff00)>>8) + "." + strconv.Itoa(int(version)&0xff)
	return s
}

func GetUintVersion(version string) uint32 {
	array := strings.Split(version, ".")
	if len(array) != 3 {
		log.Errorf("VERSION wrong version, version: %v", version)
		return 0
	}
	major, err := strconv.Atoi(array[0])
	if err != nil {
		log.Errorf("VERSION wrong version, version: %v, error: %v", version, err)
		return 0
	}
	minor, err := strconv.Atoi(array[1])
	if err != nil {
		log.Errorf("VERSION wrong version, version: %v, error: %v", version, err)
		return 0
	}
	patch, err := strconv.Atoi(array[2])
	if err != nil {
		log.Errorf("VERSION wrong version, version: %v, error: %v", version, err)
		return 0
	}
	nVersion := uint32(major<<16) + uint32(minor<<8) + uint32(patch)
	return nVersion
}


func bplIgnoreFieldRule(f reflect.StructField, version uint32) bool {
	tag := f.Tag.Get("version")
	if len(tag) == 0 {
		return false
	}
	if fieldVersion := GetUintVersion(tag); fieldVersion == 0 {
		log.Errorf("VERSION parse version tag field, field %v, tag %v", f.Name, tag)
		return false
	} else if fieldVersion <= version {
		return false
	}
	return true
}

//Usage guidance:called in main.initVersionCompatibilityRule:bpl.SetIgnoreRule("Block", version.BlockVersionRule)
func BlockVersionRule(f reflect.StructField, index int, curStructValue interface{}, topStructValue interface{}) bool {
	if index == 0 { //Header field(include version field), dont ignore
		return false
	}
	var version uint32
	switch t := curStructValue.(type) {
	case *types.Block:
		v := curStructValue.(*types.Block)
		version = v.GetVersion()
	case types.Block:
		v := curStructValue.(types.Block)
		version = v.GetVersion()
	default:
		log.Errorf("VERSION BlockVersionCompatibilityRule failed, curStructValue.type %v", t)
		return false
	}

	return bplIgnoreFieldRule(f, version)
}

func BlockHeaderVersionRule(f reflect.StructField, index int, curStructValue interface{}, topStructValue interface{}) bool {
	if index == 0 { //Header.Version, dont ignore
		return false
	}
	var version uint32
	switch t := curStructValue.(type) {
	case *types.Header:
		v := curStructValue.(*types.Header)
		version = v.Version
	case types.Header:
		v := curStructValue.(types.Header)
		version = v.Version
	default:
		log.Errorf("VERSION BlockHeaderVersionCompatibilityRule failed, curStructValue.type %v", t)
		return false
	}

	return bplIgnoreFieldRule(f, version)
}

func TrxVersionRule(f reflect.StructField, index int, curStructValue interface{}, topStructValue interface{}) bool {
	if index == 0 { //Version field, dont ignore
		return false
	}
	var version uint32
	switch t := curStructValue.(type) {
	case *types.Transaction:
		v := curStructValue.(*types.Transaction)
		version = v.Version
	case types.Transaction:
		v := curStructValue.(types.Transaction)
		version = v.Version
	default:
		log.Errorf("VERSION TrxVersionCompatibilityRule failed, curStructValue.type %v", t)
		return false
	}

	return bplIgnoreFieldRule(f, version)
}
