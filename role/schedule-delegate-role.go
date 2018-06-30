package role

import (
	"encoding/json"
	"errors"
	"math/big"
	"math/rand"
	"reflect"
	"sort"

	log "github.com/cihub/seelog"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/db"
)

//ScheduleDelegateObjectName is scheduledelegate
const ScheduleDelegateObjectName string = "scheduledelegate"

//ScheduleDelegate is singleton role
type ScheduleDelegate struct {
	CurrentTermTime *big.Int
}

//SetScheduleDelegateRole is seting scheduled delegate role
func SetScheduleDelegateRole(ldb *db.DBService, value *ScheduleDelegate) error {
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		log.Error("Set object", ScheduleDelegateObjectName, "failed")
		return err
	}

	return ldb.SetObject(ScheduleDelegateObjectName, "my", string(jsonvalue))
}

//GetScheduleDelegateRole is to get schedulated delegate role
func GetScheduleDelegateRole(ldb *db.DBService) (*ScheduleDelegate, error) {
	value, err := ldb.GetObject(ScheduleDelegateObjectName, "my")
	if err != nil {
		log.Error("GetObject object", ScheduleDelegateObjectName, "failed")
		return nil, err
	}

	res := &ScheduleDelegate{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil

}

//GetCandidateBySlot is to get candidate by slot
func GetCandidateBySlot(ldb *db.DBService, slotNum uint64) (string, error) {
	chainObject, err := GetChainStateRole(ldb)
	if err != nil {
		log.Error("err")
		return "", err
	}
	currentSlotNum := chainObject.CurrentAbsoluteSlot + slotNum
	currentCoreState, err := GetCoreStateRole(ldb)
	//log.Info("currentSlotNum", currentSlotNum)
	if err != nil {
		log.Error("err")
		return "", err
	}
	size := uint64(len(currentCoreState.CurrentDelegates))
	if size == 0 {
		return "", errors.New("delegate is null, please check configuration")
	}
	//log.Info("dddd", currentCoreState.CurrentDelegates)
	//log.Info("size", size)
	accountName := currentCoreState.CurrentDelegates[currentSlotNum%size]
	return accountName, nil

}

//ResetCandidatesTerm is reseting candidates term
func ResetCandidatesTerm(ldb *db.DBService) {
	sch := &ScheduleDelegate{big.NewInt(0)}
	SetScheduleDelegateRole(ldb, sch)
	ResetAllDelegateNewTerm(ldb)
}

//SetCandidatesTerm is setting candidates term
func SetCandidatesTerm(ldb *db.DBService, termTime *big.Int, list []string) {
	sch := &ScheduleDelegate{termTime}
	SetScheduleDelegateRole(ldb, sch)
	SetDelegateListNewTerm(ldb, termTime, list)
}

//ElectNextTermDelegatesRole is to elect next term delegates
func ElectNextTermDelegatesRole(ldb *db.DBService) []string {
	var tmpList []string
	var eligibleList []string
	var eligibles []string

	sortedDelegates, err := GetAllSortVotesDelegates(ldb)
	if err != nil {
		return nil
	}

	filterDgates := FilterOutgoingDelegate(ldb)

	if len(filterDgates) == 0 {
		tmpList = sortedDelegates
	} else {
		tmpList = common.Filter(sortedDelegates, filterDgates)
	}
	if uint32(len(tmpList)) < config.BLOCKS_PER_ROUND {
		//panic("Not enough active producers registered to schedule a round")
		return nil
	}

	candidates := tmpList[0:config.VOTED_DELEGATES_PER_ROUND]

	//sort candidates by account name
	sort.Strings(candidates)

	//Check exist ownername
	finishdelegates, err := GetAllSortFinishTimeDelegates(ldb)
	if err != nil {
		return nil
	}
	log.Info("finish delegates", finishdelegates)

	if len(filterDgates) == 0 {
		eligibleList = finishdelegates
	} else {
		eligibleList = common.Filter(finishdelegates, filterDgates)
	}

	//filter from candidates with number config.VOTED_DELEGATES_PER_ROUND

	eligibles = common.Filter(eligibleList, candidates)

	count := config.BLOCKS_PER_ROUND - config.VOTED_DELEGATES_PER_ROUND
	if count != 1 {
		//panic("invalid configuration BLOCKS_PER_ROUND and VOTED_DELEGATES_PER_ROUND")
		return nil
	}
	if len(eligibles) == 0 {
		//panic("not enough eligible delegates")
		return nil
	}
	lastTermUp := eligibles[0] //count -1 = 0

	//get final reporter lists
	reporterList := append(candidates, lastTermUp)
	newCandidates, err := GetDelegateVotesRoleByAccountName(ldb, lastTermUp)
	if err != nil {
		return nil
	}
	if (config.BLOCKS_PER_ROUND >= uint32(len(finishdelegates))) && (newCandidates.TermFinishTime.Cmp(common.MaxUint128()) == -1) {
		ResetCandidatesTerm(ldb)
	} else {
		SetCandidatesTerm(ldb, newCandidates.TermFinishTime, reporterList)
	}

	log.Info("elect next term", reporterList)

	return reporterList

}

//ShuffleEelectCandidateList is to shuffle the candidates in one round
func ShuffleEelectCandidateListRole(ldb *db.DBService, block *types.Block) ([]string, error) {
	newSchedule := ElectNextTermDelegatesRole(ldb)
	currentState, err := GetCoreStateRole(ldb)
	if err != nil {
		return nil, err
	}
	changes := common.Filter(currentState.CurrentDelegates, newSchedule)
	equal := reflect.DeepEqual(block.Header.DelegateChanges, changes)
	if equal == false {
		log.Info("invalid block changes")
		return nil, errors.New("Unexpected round changes in new block header")
	}

	h := block.Hash()
	label := h.Label()
	rand.New(rand.NewSource(int64(label)))
	rand.Shuffle(len(newSchedule), func(i, j int) {
		newSchedule[i], newSchedule[j] = newSchedule[j], newSchedule[i]
	})

	return newSchedule, nil
}
