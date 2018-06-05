package types

import (
	"crypto/sha256"
	"github.com/bottos-project/bottos/common"
	"github.com/golang/protobuf/proto"
)

func (trx *Transaction) Hash() common.Hash {
	data, _ := proto.Marshal(trx)
	temp := sha256.Sum256(data)
	hash := sha256.Sum256(temp[:])
	return hash
}

func (trx *Transaction) ValidateSign() bool {
	return true
}
