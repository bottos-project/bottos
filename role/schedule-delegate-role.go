package role

import (
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/bottos-project/core/config"
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

func (s *ScheduleDelegate) ElectNextTermDelegates(ldb *db.DBService) []string {
	var tmpList []string
	dgates := GetAllSortVotesDelegates(ldb)
	fDgates := FilterOutgoingDelegate(ldb)
	for _, dgate := range dgates {
		for _, fdgate := range fDgates {
			if dgate == fdgate {
				continue
			}
			tmpList = append(tmpList, dgate)
		}
	}
	if uint32(len(tmpList)) <= config.BLOCKS_PER_ROUND {
		return nil
	}
	candidates := tmpList[0:17]
	sort.Strings(candidates)

	//TODO Check exist ownername
	var eligibleList []string
	ftdelegates := GetAllSortFinishTimeDelegates(ldb)
	for _, ft := range ftdelegates {
		for _, fdgate := range fDgates {
			if ft == fdgate {
				continue
			}
			eligibleList = append(eligibleList, ft)
		}
	}

	//filter votesList
	var eligibles []string

	for _, list := range eligibleList {
		for _, votes := range tmpList {
			if list == votes {
				continue
			}
			eligibles = append(eligibles, list)
		}
	}
	count := config.BLOCKS_PER_ROUND - config.VOTED_DELEGATES_PER_ROUND
	lastTermUp := eligibles[0 : count-1]
	//get final reporter lists
	reporterList := append(candidates, lastTermUp...)

	return reporterList

}
