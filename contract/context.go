package contract

import (
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract/contractdb"
	"github.com/bottos-project/bottos/role"
)

//Context for contracts
type Context struct {
	RoleIntf   role.RoleInterface
	ContractDB *contractdb.ContractDB
	Trx        *types.Transaction
}

//GetTrxParam for contracts
func (ctx *Context) GetTrxParam() []byte {
	return ctx.Trx.Param
}

//GetTrxParamSize for contracts
func (ctx *Context) GetTrxParamSize() uint32 {
	size := len(ctx.Trx.Param)
	return uint32(size)
}
