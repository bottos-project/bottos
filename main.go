package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"net/http"
	"strconv"

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
	"github.com/bottos-project/bottos/restful/handler"
	"github.com/bottos-project/bottos/transaction"
	log "github.com/cihub/seelog"
	"github.com/micro/go-micro"
	"github.com/bottos-project/bottos/cmd"

	cli "gopkg.in/urfave/cli.v1"
	"runtime"
	"fmt"
)

var (
	app = cli.NewApp()
)

func init() {
	app.Usage = "the bottos command line interface"
	app.Version = "3.2.0"
	app.Copyright = "Copyright 2017~2022 The Bottos Authors"
	app.Flags = []cli.Flag {
		cmd.ConfigFileFlag,
		cmd.GenesisFileFlag,
		cmd.DataDirFlag,
		cmd.DisableAPIFlag,
		cmd.APIPortFlag,
		cmd.DisableRPCFlag,
		cmd.RPCPortFlag,
		cmd.P2PPortFlag,
		cmd.ServerAddrFlag,
		cmd.PeerListFlag,
		cmd.DelegateSignkeyFlag,
		cmd.DelegateFlag,
		cmd.EnableStaleReportFlag,
		cmd.EnableMongoDBFlag,
		cmd.MongoDBFlag,
		cmd.LogConfigFlag,
	}
	app.Action = startBottos
	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadConfig(ctx *cli.Context) {
	config.InitConfig()

	if err := config.InitLogConfig(ctx); err != nil {
		os.Exit(1)
	}

	if err := config.LoadConfig(ctx); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
}

func startBottos(ctx *cli.Context) error {
	loadConfig(ctx)

	blockDBPath := filepath.Join(config.Param.DataDir, "block/")
	stateDBPath := filepath.Join(config.Param.DataDir, "state.db")
	dbInst := db.NewDbService(blockDBPath, stateDBPath, config.Param.OptionDb)
	if dbInst == nil {
		log.Critical("Create DB service fail")
		os.Exit(1)
	}

	roleIntf := role.NewRole(dbInst)
	contractDB := contractdb.NewContractDB(dbInst)

	nc, err := contract.NewNativeContract(roleIntf)
	if err != nil {
		log.Critical("Create Native Contract error: ", err)
		os.Exit(1)
	}

	chain, err := chain.CreateBlockChain(dbInst, roleIntf, nc)
	if err != nil {
		log.Critical("Create BlockChain error: ", err)
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

	//Enabled RestFul Api
	if config.Param.RestFulApiServiceEnable {
		go startRestApi(roleIntf, contractDB)
	}

	//Enabled Rpc Api
	if config.Param.RpcServiceEnable {
		go startRPCService(actorenv)
	}

	WaitSystemDown(chain, multiActors)
	return nil
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

func startRestApi(roleIntf role.RoleInterface, contractDB *contractdb.ContractDB) {
	router := handler.NewRouter()
	//transfer to restful handler
	handler.SetRoleIntf(roleIntf)
	handler.SetContractDbIns(contractDB)
	log.Critical(http.ListenAndServe(config.Param.ServInterAddr+":"+strconv.Itoa(config.Param.APIPort), router))
}

func startRPCService(actorenv *actionenv.ActorEnv) {
	repo := caapi.NewApiService(actorenv)

	service := micro.NewService(
		micro.Name(config.Param.RpcServiceName),
		micro.Version(config.Param.RpcServiceVersion),
	)

	api.RegisterChainHandler(service.Server(), repo)
	if err := service.Run(); err != nil {
		log.Critical("RPC Service fail: ", err)
	}
}
