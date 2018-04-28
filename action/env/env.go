package env

import (
	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/chain/extra"
	"github.com/bottos-project/core/role"
)

type ActorEnv struct {
	RoleIntf	role.RoleInterface
	Chain   	chain.BlockChainInterface
	TxStore 	*txstore.TransactionStore
}
