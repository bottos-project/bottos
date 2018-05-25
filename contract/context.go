package contract

import (
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract/contractdb"
)

type Context struct {
	RoleIntf role.RoleInterface
	ContractDB *contractdb.ContractDB
	Trx *types.Transaction
}

func (ctx *Context) GetTrxParam() []byte {
	return ctx.Trx.Param
}

func (ctx *Context) GetTrxParamSize() uint32 {
	size := len(ctx.Trx.Param)
	return uint32(size)
}
