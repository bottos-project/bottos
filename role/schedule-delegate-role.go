package role

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/bottos-project/core/db"
)

type ScheduleDelegate struct {
	CurrentRaceTime big.Int
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

func ResetProducerRace(ldb *db.DBService) {
	//	auto ResetRace = [&db](const producer_votes_object& pvo) {
	//	      db.modify(pvo, [](producer_votes_object& pvo) {
	//	         pvo.start_new_race_lap(0);
	//	      });
	//	   };
	//	   const auto& AllProducers = db.get_index<producer_votes_multi_index, by_votes>();

	//	   boost::for_each(AllProducers, ResetRace);
	//	   db.modify(*this, [](producer_schedule_object& pso) {
	//	      pso.currentRaceTime = 0;
	//	   });
	//	}

	//	} } // namespace eosio::chain

}
