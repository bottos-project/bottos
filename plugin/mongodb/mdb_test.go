package mongodb

import (
	log "github.com/cihub/seelog"
	"testing"

	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/db"
)

func TestPersistanceRole_writedb(t *testing.T) {
	ins := db.NewDbService("./file2", "./file2/db.db", "10.104.14.169:27017")
	block := &types.Block{}
	err := ApplyPersistanceRole(ins, block)
	if err != nil {
		log.Error(err)
	}
}
