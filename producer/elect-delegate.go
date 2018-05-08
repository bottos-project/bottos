package producer

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"

	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
)

func StringSliceReflectEqual(a, b []string) bool {
	return reflect.DeepEqual(a, b)
}
func (r *Reporter) ShuffleEelectCandidateList(block types.Block) ([]string, error) {
	newSchedule := r.roleIntf.ElectNextTermDelegates()
	currentState, err := r.roleIntf.GetCoreState()
	if err != nil {
		return nil, err
	}
	changes := common.Filter(currentState.CurrentDelegates, newSchedule)
	equal := reflect.DeepEqual(block.Header.DelegateChanges, changes)
	if equal == false {
		fmt.Println("invalid block changes")
		panic(1)
		return nil, errors.New("Unexpected round changes in new block header")
	}
	for i := range newSchedule {
		j := rand.Intn(i + 1)
		newSchedule[i], newSchedule[j] = newSchedule[j], newSchedule[i]
	}
	return newSchedule, nil
}
