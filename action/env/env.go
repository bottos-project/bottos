package env

import (
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/chain/extra"
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/contract/contractdb"
)

type ActorEnv struct {
	RoleIntf	role.RoleInterface
	ContractDB  *contractdb.ContractDB
	Chain   	chain.BlockChainInterface
	TxStore 	*txstore.TransactionStore
	NcIntf		contract.NativeContractInterface
}
