package main

import (
	cactor "github.com/bottos-project/bottos/action/actor"
	caapi "github.com/bottos-project/bottos/action/actor/api"
	"github.com/bottos-project/bottos/action/actor/transaction"
	actionenv "github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/chain/extra"
	"github.com/bottos-project/bottos/cmd"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/restful/handler"
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/transaction"
	log "github.com/cihub/seelog"
	"github.com/micro/go-micro"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"fmt"
	"gopkg.in/urfave/cli.v1"
	"runtime"
	"strconv"
)

var (
	app = cli.NewApp()
)

func init() {
	app.Usage = "the bottos command line interface"
	app.Version = "3.2.0"
	app.Copyright = "Copyright 2017~2022 The Bottos Authors"
	app.Flags = []cli.Flag{
		cmd.ConfigFileFlag,
		cmd.GenesisFileFlag,
		cmd.DataDirFlag,
		cmd.DisableRESTFlag,
		cmd.RESTPortFlag,
		cmd.RESTServerAddrFlag,
		cmd.DisableRPCFlag,
		cmd.RPCPortFlag,
		cmd.P2PPortFlag,
		cmd.P2PServerAddrFlag,
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

	log.Infof("Bottos ChainID: %x", config.GetChainID())
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
		Chain:      chain,
		TxStore:    txStore,
		NcIntf:     nc,
	}
	multiActors := cactor.InitActors(actorenv)

	var trxPool = transaction.InitTrxPool(dbInst, roleIntf, nc, multiActors.GetNetActor())
	trxactor.SetTrxPool(trxPool)

	//start RESTful Api
	if !ctx.GlobalBool(cmd.DisableRESTFlag.Name) {
		go startRestApi(roleIntf)
	}

	//start Rpc Api
	if !ctx.GlobalBool(cmd.DisableRPCFlag.Name) {
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

func startRestApi(roleIntf role.RoleInterface) {
	router := handler.NewRouter()
	//transfer to restful handler
	handler.SetRoleIntf(roleIntf)
	err := http.ListenAndServe(config.Param.RESTServAddr+":"+strconv.Itoa(config.Param.RESTPort), router)
	if err != nil {
		log.Critical("RESTful server fail: ", err)
		os.Exit(1)
	}
}

func startRPCService(actorenv *actionenv.ActorEnv) {
	repo := caapi.NewApiService(actorenv)

	service := micro.NewService(
		micro.Name(config.Param.RpcServiceName),
		micro.Version(config.Param.RpcServiceVersion),
	)

	api.RegisterChainHandler(service.Server(), repo)
	if err := service.Run(); err != nil {
		log.Critical("RPC server fail: ", err)
	}
}
