package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/chain/extra"
	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/role"
	"github.com/bottos-project/core/api"

	"github.com/bottos-project/core/common/types"

	"github.com/bottos-project/core/contract"
	"github.com/bottos-project/core/contract/contractdb"

	"github.com/micro/go-micro"
	cactor "github.com/bottos-project/core/action/actor"
	caapi "github.com/bottos-project/core/action/actor/api"
	"github.com/bottos-project/core/action/actor/transaction"
	actionenv "github.com/bottos-project/core/action/env"
	"github.com/bottos-project/core/transaction"
	wasm "github.com/bottos-project/core/vm/wasm/exec"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Load config fail")
		os.Exit(1)
	}

	dbInst := db.NewDbService(config.Param.DataDir, filepath.Join(config.Param.DataDir, "blockchain"))
	if dbInst == nil {
		fmt.Println("Create DB service fail")
		os.Exit(1)
	}

	roleIntf := role.NewRole(dbInst)
	contractDB := contractdb.NewContractDB(dbInst)

	nc, err := contract.NewNativeContract(roleIntf)
	if err != nil {
		fmt.Println("Create Native Contract error: ", err)
		os.Exit(1)
	}

	chain, err := chain.CreateBlockChain(dbInst, roleIntf)
	if err != nil {
		fmt.Println("Create BlockChain error: ", err)
		os.Exit(1)
	}

	txStore := txstore.NewTransactionStore(chain, roleIntf)

	actorenv := &actionenv.ActorEnv{
		RoleIntf:	roleIntf, 
		ContractDB: contractDB,
		Chain:		chain, 
		TxStore:	txStore,
		NcIntf:		nc,
	}
	cactor.InitActors(actorenv)
	//caapi.PushTransaction(2876568)

	//caapi.InitTrxActorAgent()
	var trxPool = transaction.InitTrxPool(actorenv)
	trxactor.SetTrxPool(trxPool)

	if config.Param.ApiServiceEnable {
		repo := caapi.NewApiService(actorenv)

		service := micro.NewService(
			micro.Name("core"),
			micro.Version("2.0.0"),
		)

		service.Init()
		api.RegisterCoreApiHandler(service.Server(), repo)
		if err := service.Run(); err != nil {
			panic(err)
		}
	}

	//bf  := []byte{0xdc, 0x00, 0x02, 0xda, 0x00, 0x08, 0x74, 0x65, 0x73, 0x74, 0x75, 0x73, 0x65, 0x72, 0xce, 0x00, 0x00, 0x00, 0x63}
	bf  := []byte{0xdc, 0x00, 0x02, 0xda, 0x00, 0x08, 0x74, 0x65, 0x73, 0x74, 0x75, 0x73, 0x65, 0x72, 0xda, 0x00, 0x08, 0x74, 0x65, 0x73, 0x74, 0x75, 0x73, 0x65, 0x72}

	trx := &types.Transaction{
		Version:1,
		CursorNum:1,
		CursorLabel:1,
		Lifetime:1,
		Sender:"bottos",
		Contract:"usermng",
		Method:"reguser",
		Param: bf,
		SigAlg:1,
		Signature:[]byte{},
	}
	ctx := &contract.Context{RoleIntf: roleIntf, ContractDB: contractDB, Trx:trx}
	_, err = wasm.GetInstance().Apply2(ctx, 1, false)
	fmt.Println(err)

	WaitSystemDown()
}

func WaitSystemDown() {
	exit := make(chan bool, 0)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	defer signal.Stop(sigc)

	go func() {
		<-sigc
		fmt.Println("System shutdown")
		close(exit)
	}()

	<-exit
}
