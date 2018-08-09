package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/chain/extra"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/role"

	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/contract/contractdb"

	cactor "github.com/bottos-project/bottos/action/actor"
	caapi "github.com/bottos-project/bottos/action/actor/api"
	"github.com/bottos-project/bottos/action/actor/transaction"
	actionenv "github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/transaction"
	log "github.com/cihub/seelog"
	"github.com/micro/go-micro"
	"github.com/bottos-project/bottos/cmdcli"
)

func main() {
	//TODO: The config's cli parameters will be used soon.
	GlobalConf, GenesisConf, err := cmdcli.Init()
	if err != nil {
		log.Error("Parse cmdcli fail")
		os.Exit(1)
	}
	err = config.LoadConfig(&GlobalConf, &GenesisConf)
	if err != nil {
		log.Error("Load config fail")
		os.Exit(1)
	}

	dbInst := db.NewDbService(config.Param.DataDir, filepath.Join(config.Param.DataDir, "blockchain"), config.Param.OptionDb)
	if dbInst == nil {
		log.Error("Create DB service fail")
		os.Exit(1)
	}

	roleIntf := role.NewRole(dbInst)
	contractDB := contractdb.NewContractDB(dbInst)

	nc, err := contract.NewNativeContract(roleIntf)
	if err != nil {
		log.Info("Create Native Contract error: ", err)
		os.Exit(1)
	}

	chain, err := chain.CreateBlockChain(dbInst, roleIntf, nc)
	if err != nil {
		log.Error("Create BlockChain error: ", err)
		os.Exit(1)
	}

	txStore := txstore.NewTransactionStore(chain, roleIntf)

	actorenv := &actionenv.ActorEnv{
		RoleIntf:   roleIntf,
		ContractDB: contractDB,
		Chain:      chain,
		TxStore:    txStore,
		NcIntf:     nc,
	}
	multiActors := cactor.InitActors(actorenv)

	var trxPool = transaction.InitTrxPool(actorenv, multiActors.GetNetActor())
	trxactor.SetTrxPool(trxPool)

	if config.Param.RpcServiceEnable {
		repo := caapi.NewApiService(actorenv)

		service := micro.NewService(
			micro.Name(config.Param.RpcServiceName),
			micro.Version(config.Param.RpcServiceVersion),
		)
		
		//Prompt this due to it parse cli parmeters which conflict to urfave/cli.
		//service.Init()
		api.RegisterChainHandler(service.Server(), repo)
		if err := service.Run(); err != nil {
			panic(err)
		}
	}

	WaitSystemDown(chain, multiActors)
}

//WaitSystemDown is to handle ctrl+C
func WaitSystemDown(chain chain.BlockChainInterface, actors *cactor.MultiActor) {
	exit := make(chan bool, 0)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL)

	go func() {
		for sig := range sigc {
			actors.ActorsStop()
			chain.Close()
			log.Infof("System shutdown, signal: %v", sig.String())
			log.Flush()
			close(exit)
		}
	}()

	<-exit
}
