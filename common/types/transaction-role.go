package types 

import (
	"crypto/sha256"
	"github.com/bottos-project/core/common"
	"github.com/golang/protobuf/proto"
)

func (trx *Transaction) Hash() common.Hash {
	data, _ := proto.Marshal(trx)
	hash := sha256.Sum256(data)
	return hash
}



func (trx *Transaction) ValidateSign() bool {
	return true
}
