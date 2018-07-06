package role

import (
	"math/big"

	"encoding/json"
	"errors"
	"strconv"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/db"
	log "github.com/cihub/seelog"
)

// DelegateVotesObjectName is definition of delegate vote object name
const DelegateVotesObjectName string = "delegatevotes"

// DelegateVotesObjectKeyName is definition of delegate vote object key name
const DelegateVotesObjectKeyName string = "owner_account"

// DelegateVotesObjectIndexVote is definition of delegate vote object index name
const DelegateVotesObjectIndexVote string = "votes"
const DelegateVotesObjectIndexVoteJSON string = "serve.votes"

// DelegateVotesObjectIndexFinishTime is definition of delegate vote object index finish time
const DelegateVotesObjectIndexFinishTime string = "term_finish_time"
const DelegateVotesObjectIndexFinishTimeJSON string = "serve.term_finish_time"

// Serve is definition of serve
type Serve struct {
	Votes          uint64   `json:"votes"`
	Position       *big.Int `json:"position"`
	TermUpdateTime *big.Int `json:"term_update_time"`
	TermFinishTime *big.Int `json:"term_finish_time"`
}

// DelegateVotes is definition of delegate votes
type DelegateVotes struct {
	OwnerAccount string `json:"owner_account"`
	Serve        `json:"serve"`
}

// CreateDelegateVotesRole is to save initial delegate votes
func CreateDelegateVotesRole(ldb *db.DBService) error {
	err := ldb.CreatObjectIndex(DelegateVotesObjectName, DelegateVotesObjectKeyName, DelegateVotesObjectKeyName)
	if err != nil {
		return err
	}
	err = ldb.CreatObjectMultiIndex(DelegateVotesObjectName, DelegateVotesObjectIndexVote, DelegateVotesObjectIndexVoteJSON, DelegateVotesObjectKeyName)
	if err != nil {
		return err
	}
	err = ldb.CreatObjectMultiIndex(DelegateVotesObjectName, DelegateVotesObjectIndexFinishTime, DelegateVotesObjectIndexFinishTimeJSON, DelegateVotesObjectKeyName)
	if err != nil {
		return err
	}
	return nil
}

// SetDelegateVotesRole is to save delegate votes
func SetDelegateVotesRole(ldb *db.DBService, key string, value *DelegateVotes) error {
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ldb.SetObject(DelegateVotesObjectName, key, string(jsonvalue))
}

// GetDelegateVotesRoleByAccountName is to get delegate votes by account name
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

// GetDelegateVotesRoleByVote is to get delegate votes by vote
func GetDelegateVotesRoleByVote(ldb *db.DBService, vote uint64) (*DelegateVotes, error) {
	value, err := ldb.GetObjectByIndex(DelegateVotesObjectName, DelegateVotesObjectIndexVote, strconv.FormatUint(vote, 10))
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

// GetDelegateVotesRoleByFinishTime is to get delegate votes by finish time
func GetDelegateVotesRoleByFinishTime(ldb *db.DBService, key *big.Int) (*DelegateVotes, error) {
	value, err := ldb.GetObjectByIndex(DelegateVotesObjectName, DelegateVotesObjectIndexFinishTime, key.String())
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

// update is to update delegate
func (d *DelegateVotes) update(currentVotes uint64, currentPosition *big.Int, currentTermTime *big.Int) {
	if currentTermTime.Cmp(big.NewInt(0)) == -1 || currentTermTime.Cmp(big.NewInt(0)) == -1 {
		return
	}
	termTimeToFinish := new(big.Int)
	remaining := termTimeToFinish.Sub(common.MaxUint128(), currentPosition)
	if currentVotes > 0 {
		termTimeToFinish = termTimeToFinish.Div(remaining, new(big.Int).SetUint64(currentVotes))
	} else {
		termTimeToFinish = common.MaxUint128()
	}

	if currentTermTime.Cmp(new(big.Int).Sub(common.MaxUint128(), termTimeToFinish)) == -1 {
		log.Error("ERROR currentTermTime time overflow", currentTermTime)
		return
	}
	termFinishTime := new(big.Int).Add(currentTermTime, termTimeToFinish)
	d.Serve.Votes = currentVotes
	d.Serve.Position = currentPosition
	d.Serve.TermUpdateTime = currentTermTime
	d.Serve.TermFinishTime = termFinishTime
}

// GetAllDelegateVotesRole is to get all delegate votes
func GetAllDelegateVotesRole(ldb *db.DBService) ([]*DelegateVotes, error) {
	objects, err := ldb.GetAllObjects(DelegateVotesObjectKeyName)
	if err != nil {
		log.Error("ERROR failed ", err)
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

// ResetAllDelegateNewTerm is to reset all delegate
func ResetAllDelegateNewTerm(ldb *db.DBService) {

	voteDelegates, err := GetAllDelegateVotesRole(ldb)
	if err != nil {
		return
	}
	for _, object := range voteDelegates {
		dvotes := object.StartNewTerm(big.NewInt(0))
		dvotes.OwnerAccount = object.OwnerAccount
		SetDelegateVotesRole(ldb, object.OwnerAccount, dvotes)
		//log.Info("ResetAllDelegateNewTerm", object.OwnerAccount, dvotes)
	}
}

// SetDelegateListNewTerm is to set delegate list new term
func SetDelegateListNewTerm(ldb *db.DBService, termTime *big.Int, lists []string) {
	var mylists = make([]string, len(lists))
	copy(mylists, lists)
	for _, accountName := range mylists {
		delegate, err := GetDelegateVotesRoleByAccountName(ldb, accountName)
		if err != nil {
			return
		}
		dvotes := delegate.StartNewTerm(termTime)
		SetDelegateVotesRole(ldb, accountName, dvotes)
		//log.Info("set delegate new term", accountName, dvotes)

	}
}

// StartNewTerm is to start new term
func (d *DelegateVotes) StartNewTerm(currentTermTime *big.Int) *DelegateVotes {
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

// UpdateVotes is to update votes
func (d *DelegateVotes) UpdateVotes(votes uint64, currentTermTime *big.Int) {
	timeSinceLastUpdate := new(big.Int).Sub(currentTermTime, d.Serve.TermUpdateTime)
	myVotes := new(big.Int).Mul(new(big.Int).SetUint64(d.Serve.Votes), timeSinceLastUpdate)
	newPosition := new(big.Int).Add(d.Serve.Position, myVotes)
	newSpeed := d.Serve.Votes + votes

	d.update(newSpeed, newPosition, currentTermTime)
}

// GetAllSortVotesDelegates is to get all sort votes delegates
func GetAllSortVotesDelegates(ldb *db.DBService) ([]string, error) {
	var objects []string
	var err error
	objects, err = ldb.GetAllObjectsSortByIndex(DelegateVotesObjectIndexVote)
	if err != nil {
		return nil, err
	}
	var accounts []string

	for _, object := range objects {
		res := new(DelegateVotes)
		err = json.Unmarshal([]byte(object), res)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, res.OwnerAccount)
	}
	var accountRtn = make([]string, len(accounts))
	copy(accountRtn, accounts)
	return accountRtn, nil
}

// GetAllSortFinishTimeDelegates is to get all sort finish time delegates
func GetAllSortFinishTimeDelegates(ldb *db.DBService) ([]string, error) {
	var objects []string
	var err error
	objects, err = ldb.GetAllObjectsSortByIndex(DelegateVotesObjectIndexFinishTime)
	if err != nil {
		return nil, err
	}
	var accounts []string
	for _, object := range objects {
		res := &DelegateVotes{}
		err = json.Unmarshal([]byte(object), res)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, res.OwnerAccount)
	}
	var accountRtn = make([]string, len(accounts))
	copy(accountRtn, accounts)
	return accountRtn, nil
}
