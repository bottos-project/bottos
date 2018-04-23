package env

import (
	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/chain"
)

type ActorEnv struct {
	Db		*db.DBService
	Chain	chain.BlockChainInterface
}

