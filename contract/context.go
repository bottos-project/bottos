package contract

import (
	"github.com/bottos-project/core/role"
	"github.com/bottos-project/core/common/types"
)

type Context struct {
	roleIntf role.RoleInterface
	Trx *types.Transaction
}
