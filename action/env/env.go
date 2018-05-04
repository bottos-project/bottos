package env

import (
	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/chain/extra"
	"github.com/bottos-project/core/role"
	"github.com/bottos-project/core/contract"
	"github.com/bottos-project/core/contract/contractdb"
)

type ActorEnv struct {
	RoleIntf	role.RoleInterface
	ContractDB  *contractdb.ContractDB
	Chain   	chain.BlockChainInterface
	TxStore 	*txstore.TransactionStore
	NcIntf		contract.NativeContractInterface
}
