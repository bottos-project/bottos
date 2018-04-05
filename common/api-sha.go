package common

import (
	"github.com/bottos-project/bottos/core/common/types"
	"github.com/bottos-project/bottos/core/db"
)

type ApiList interface {
	Len() int
	GetRlp(i int) []byte
}

func ApiSha(list ApiList) types.Hash {
	db, _ := db.NewMemDatabase()
	//lmq
	//	trie := trie.New(nil, db)
	//	for i := 0; i < list.Len(); i++ {
	//		key, _ := rlp.EncodeToBytes(uint(i))
	//		trie.Update(key, list.GetRlp(i))
	//	}

	return types.BytesToHash(trie.Root())
}
