package env

import (
	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/chain/extra"
	"github.com/bottos-project/core/db"
)

type ActorEnv struct {
	Db      *db.DBService
	Chain   chain.BlockChainInterface
	TxStore *txstore.TransactionStore
}
