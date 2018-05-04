package contract

import (
	"github.com/bottos-project/core/role"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/contract/contractdb"
)

type Context struct {
	RoleIntf role.RoleInterface
	ContractDB *contractdb.ContractDB
	Trx *types.Transaction
}
