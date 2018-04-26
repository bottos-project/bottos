package role

import (
	"fmt"

	"github.com/bottos-project/core/db"
)

func GetScheduleDelegateRole(ldb *db.DBService, slotNum uint32) (string, error) {
	chainObject, err := GetChainStateObjectRole(ldb)
	if err != nil {
		fmt.Println("err")
		return "", err
	}
	currentSlotNum := chainObject.CurrentAbsoluteSlot + uint64(slotNum)
	currentCoreState, err := GetGlobalPropertyRole(ldb)
	if err != nil {
		fmt.Println("err")
		return "", err
	}
	size := uint64(len(currentCoreState.CurrentDelegates))
	accountName := currentCoreState.CurrentDelegates[currentSlotNum%size]
	return accountName, nil

}
