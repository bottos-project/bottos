package role

import (
	"math/big"

	"encoding/json"

	"errors"
	"fmt"

	"github.com/bottos-project/core/common"

	"github.com/bottos-project/core/db"
)

const DelegateVotesObjectName string = "delegate"
const DelegateVotesObjectKeyName string = "account_name"
const DelegateVotesObjectIndexVote string = "votes"
const DelegateVotesObjectIndexFinishTime string = "finish_time"

type Serve struct {
	Votes          uint64
	Position       *big.Int //uint128
	TermUpdateTime *big.Int //uint128
	TermFinishTime *big.Int //uint128

}
type DelegateVotes struct {
	OwnerAccount string
	Serve
}

func CreateDelegateVotesRole(ldb *db.DBService) error {
	err := ldb.CreatObjectIndex(DelegateVotesObjectName, DelegateVotesObjectKeyName, DelegateVotesObjectKeyName)
	if err != nil {
		return err
	}
	err = ldb.CreatObjectIndex(DelegateVotesObjectName, DelegateVotesObjectIndexVote, DelegateVotesObjectIndexVote)
	if err != nil {
		return err
	}
	err = ldb.CreatObjectIndex(DelegateVotesObjectName, DelegateVotesObjectIndexFinishTime, DelegateVotesObjectIndexFinishTime)
	if err != nil {
		return err
	}
	return nil
}

func SetDelegateVotesRole(ldb *db.DBService, key string, value *DelegateVotes) error {
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ldb.SetObject(DelegateVotesObjectName, key, string(jsonvalue))
}

func GetDelegateVotesRoleByAccountName(ldb *db.DBService, key string) (*DelegateVotes, error) {
	value, err := ldb.GetObject(DelegateVotesObjectName, key)
	if err != nil {
		return nil, err
	}

	res := &DelegateVotes{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func GetDelegateVotesRoleByVote(ldb *db.DBService, key string) (*DelegateVotes, error) {
	value, err := ldb.GetObjectByIndex(DelegateVotesObjectName, DelegateVotesObjectIndexVote, key)
	if err != nil {
		return nil, err
	}

	res := &DelegateVotes{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func GetDelegateVotesRoleByFinishTime(ldb *db.DBService, key string) (*DelegateVotes, error) {
	value, err := ldb.GetObjectByIndex(DelegateVotesObjectName, DelegateVotesObjectIndexFinishTime, key)
	if err != nil {
		return nil, err
	}

	res := &DelegateVotes{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (d *DelegateVotes) update(currentVotes uint64, currentPosition *big.Int, currentTermTime *big.Int) {
	if currentTermTime.Cmp(big.NewInt(0)) == -1 || currentTermTime.Cmp(big.NewInt(0)) == -1 {
		return
	}
	termTimeToFinish := new(big.Int)
	remaining := termTimeToFinish.Sub(common.MaxUint128(), currentPosition)
	timeFinish := termTimeToFinish.Div(remaining, new(big.Int).SetUint64(currentVotes))
	if currentVotes > 0 {
		termTimeToFinish = timeFinish
	} else {
		termTimeToFinish = common.MaxUint128()
	}

	if currentTermTime.Cmp(new(big.Int).Sub(common.MaxUint128(), termTimeToFinish)) == -1 {
		fmt.Println("currentTermTime  time overflow", currentTermTime)
		return
	}
	termFinishTime := new(big.Int).Add(currentTermTime, termTimeToFinish)
	d.Serve.Votes = currentVotes
	d.Serve.Position = currentPosition
	d.Serve.TermUpdateTime = currentTermTime
	d.Serve.TermFinishTime = termFinishTime
}
func GetAllDelegateVotes(ldb *db.DBService) ([]*DelegateVotes, error) {
	objects, err := ldb.GetAllObjects(DelegateVotesObjectName)
	if err != nil {
		return nil, err
	}
	var dgates = []*DelegateVotes{}
	for _, object := range objects {
		res := &DelegateVotes{}
		err = json.Unmarshal([]byte(object), res)
		if err != nil {
			return nil, errors.New("invalid object to Unmarshal" + object)
		}
		dgates = append(dgates, res)
	}
	return dgates, nil

}

//TODO
func ResetAllDelegateNewTerm(ldb *db.DBService) {

	voteDelegates, err := GetAllDelegateVotes(ldb)
	if err != nil {
		return
	}
	for _, object := range voteDelegates {
		dvotes := object.startNewTerm(big.NewInt(0))
		dvotes.OwnerAccount = object.OwnerAccount
		SetDelegateVotesRole(ldb, object.OwnerAccount, dvotes)
		fmt.Println("ResetAllDelegateNewTerm", object.OwnerAccount, dvotes)
	}
}

func SetDelegateListNewTerm(ldb *db.DBService, termTime *big.Int, lists []string) {
	for _, accountName := range lists {
		delegate, err := GetDelegateVotesRoleByAccountName(ldb, accountName)
		if err != nil {
			return
		}
		dvotes := delegate.startNewTerm(termTime)
		SetDelegateVotesRole(ldb, accountName, dvotes)
		fmt.Println("key", accountName, dvotes)

	}
}

func (d *DelegateVotes) startNewTerm(currentTermTime *big.Int) *DelegateVotes {
	d.update(d.Serve.Votes, big.NewInt(0), currentTermTime)
	return &DelegateVotes{
		OwnerAccount: d.OwnerAccount,
		Serve: Serve{
			Votes:          d.Serve.Votes,
			Position:       d.Serve.Position,
			TermUpdateTime: d.Serve.TermUpdateTime,
			TermFinishTime: d.Serve.TermFinishTime,
		},
	}

}

func (d *DelegateVotes) UpdateVotes(votes uint64, currentTermTime *big.Int) {
	timeSinceLastUpdate := new(big.Int).Sub(currentTermTime, d.Serve.TermUpdateTime)
	myVotes := new(big.Int).Mul(new(big.Int).SetUint64(d.Serve.Votes), timeSinceLastUpdate)
	newPosition := new(big.Int).Add(d.Serve.Position, myVotes)
	newSpeed := d.Serve.Votes + votes

	d.update(newSpeed, newPosition, currentTermTime)
}

func GetAllSortVotesDelegates(ldb *db.DBService) []string {
	objects, err := ldb.GetAllObjectsSortByIndex(DelegateVotesObjectName, DelegateVotesObjectIndexVote)
	if err != nil {
		return nil
	}
	var accounts = []string{}
	for _, object := range objects {
		res := &DelegateVotes{}
		err = json.Unmarshal([]byte(object), res)
		if err != nil {
			return nil
		}
		accounts = append(accounts, res.OwnerAccount)
	}
	return accounts
}

func GetAllSortFinishTimeDelegates(ldb *db.DBService) []string {
	objects, err := ldb.GetAllObjectsSortByIndex(DelegateVotesObjectName, DelegateVotesObjectIndexFinishTime)
	if err != nil {
		return nil
	}
	var accounts = []string{}
	for _, object := range objects {
		res := &DelegateVotes{}
		err = json.Unmarshal([]byte(object), res)
		if err != nil {
			return nil
		}
		accounts = append(accounts, res.OwnerAccount)
	}
	return accounts
}
