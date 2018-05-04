package role

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/bottos-project/core/db"
)

type ScheduleDelegate struct {
	CurrentTermTime *big.Int
	DelegateVotes
}

func GetScheduleDelegateRole(ldb *db.DBService, slotNum uint32) (string, error) {
	chainObject, err := GetChainStateRole(ldb)
	if err != nil {
		fmt.Println("err")
		return "", err
	}
	fmt.Println("currentSlotNum", chainObject.CurrentAbsoluteSlot, slotNum)
	currentSlotNum := chainObject.CurrentAbsoluteSlot + uint64(slotNum)
	currentCoreState, err := GetCoreStateRole(ldb)
	fmt.Println("currentSlotNum", currentSlotNum)
	if err != nil {
		fmt.Println("err")
		return "", err
	}
	size := uint64(len(currentCoreState.CurrentDelegates))
	if size == 0 {
		return "", errors.New("delegate is null, please check configuration")
	}
	fmt.Println("size", size)
	accountName := currentCoreState.CurrentDelegates[currentSlotNum%size]
	return accountName, nil

}

func (s *ScheduleDelegate) ResetDelegateTerm(ldb *db.DBService) {
	s.CurrentTermTime = big.NewInt(0)
	s.DelegateVotes.ResetAllDelegateNewTerm(ldb)
}

func (s *ScheduleDelegate) ElectNextTermDelegates(ldb *db.DBService) {

}
