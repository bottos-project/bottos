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
	"github.com/bottos-project/bottos/action/actor/transaction/trxprehandleactor"
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
	"github.com/bottos-project/bottos/plugin/mongodb"
	"github.com/bottos-project/bottos/common/types"
	"github.com/AsynkronIT/protoactor-go/actor"
)

var (
	app = cli.NewApp()
)

func init() {
	app.Usage = config.USAGE
	app.Version = version.GetAppVersionString()
	app.Copyright = config.COPYRIGHT
	app.Flags = []cli.Flag{
		cmd.ConfigFileFlag,
		cmd.GenesisFileFlag,
		cmd.DataDirFlag,
		cmd.DisableRESTFlag,
		cmd.RESTPortFlag,
		//cmd.EnableRPCFlag,
		//cmd.RPCPortFlag,
		cmd.P2PPortFlag,
		cmd.P2PServerAddrFlag,
		cmd.RESTServerAddrFlag,
		cmd.PeerListFlag,
		cmd.DelegateSignkeyFlag,
		cmd.DelegateFlag,
		cmd.DelegatePrateFlag,
		cmd.EnableMongoDBFlag,
		cmd.MongoDBFlag,
		cmd.LogConfigFlag,
		cmd.WalletDirFlag,
		cmd.EnableWalletFlag,
		cmd.WalletRESTPortFlag,
		cmd.WalletRESTServerAddrFlag,
		cmd.DebugFlag,
		cmd.LogMinLevelFlag,
		cmd.LogMaxLevelFlag,
		cmd.LogLevelsFlag,
		cmd.LogMaxrollsFlag,
		cmd.RecoverAtBlockNumFlag,
		cmd.RecoverFromDataDirFlag,
		cmd.RestMaxLimit,
		cmd.WalletRestMaxLimit,
	}
	app.Action = startBottos
	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		if ctx.GlobalBool(cmd.DebugFlag.Name) {
			go func() {
				http.ListenAndServe("0.0.0.0:6060", nil)
			}()
		}
		return nil
	}
	app.After = func(ctx *cli.Context) error {
		if ctx.GlobalBool(cmd.DebugFlag.Name) {
			saveHeapProfile()
		}
		return nil
	}
}

func main() {

	signal.Ignore(syscall.SIGHUP)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadConfig(ctx *cli.Context) {
	config.InitConfig()

	if err := config.LoadConfig(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := config.InitLogConfig(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initVersion() {
	if err := version.Init(); err != nil {
		log.Critical(err)
		log.Flush()
		os.Exit(1)
	}
}

func recover(ctx *cli.Context) {
	var err error
	blockNum := uint64(ctx.GlobalInt(cmd.RecoverAtBlockNumFlag.Name))
	datadir := ctx.GlobalString(cmd.RecoverFromDataDirFlag.Name)
	if blockNum == 0 && len(datadir) == 0 {
		err = errors.New("param error.")
	}
	if err == nil {
		err = recoverAtBlockNumber(ctx, blockNum, datadir)
	}
	if err != nil {
		fmt.Printf("recover error: %v\n", err)
		log.Critical(err)
		log.Flush()
		os.Exit(1)
	} else {
		fmt.Println("recover finished.")
		log.Info("recover finished.")
	}
}

func startBottos(ctx *cli.Context) error {
	loadConfig(ctx)

	initVersion()

        //start Wallet REST Api
	if ctx.GlobalBool(cmd.EnableWalletFlag.Name) {
		var walletRestMaxLimit int
		if ctx.GlobalIsSet(cmd.WalletRestMaxLimit.Name) {
			walletRestMaxLimit = ctx.GlobalInt(cmd.WalletRestMaxLimit.Name)
		}
		go startWalletRestApi(walletRestMaxLimit)
	}

	if ctx.GlobalIsSet(cmd.RecoverAtBlockNumFlag.Name) || ctx.GlobalIsSet(cmd.RecoverFromDataDirFlag.Name) {
		recover(ctx)
	}

	blockDBPath := filepath.Join(config.BtoConfig.Node.DataDir, "data/block/")
	stateDBPath := filepath.Join(config.BtoConfig.Node.DataDir, "data/state.db")
	dbInst := db.NewDbService(blockDBPath, stateDBPath)
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

	if ctx.GlobalBool(cmd.EnableMongoDBFlag.Name) {
		var mdbActor *actor.PID = nil
		if ctx.GlobalBool(cmd.EnableMongoDBFlag.Name) {
			mdbActor = startMangoDB(roleIntf)
			if mdbActor == nil {
				log.Critical("Start MongoDB service fail")
				log.Flush()
				os.Exit(1)
			}
		}
		chain.RegisterCommittedBlockCallback(func (block *types.Block) {
			mdbActor.Tell(block)
		})
	}

	if err := chain.Init(); err != nil {
		log.Critical("Initialize BlockChain error: ", err)
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
	trxprehandleactor.SetTrxPool(trxPool)
	trxprehandleactor.SetTrxActor(multiActors.GetTrxActor())

	//start RESTful Api
	if !ctx.GlobalBool(cmd.DisableRESTFlag.Name) {
		go startRestApi(roleIntf)
	}

	//enable Rpc Api
	if ctx.GlobalBool(cmd.EnableRPCFlag.Name) {
		go startRPCService(actorenv)
	}

	WaitSystemDown(chain, multiActors)
	return nil
}

func saveHeapProfile() {
	log.Infof("begin save memory")
	runtime.GC()
	f, err := os.Create(fmt.Sprintf("heap_%s_%d_%s.prof", "bottos", 112233, time.Now().Format("2006_01_02_03_04_05")))
	if err != nil {
		log.Infof("error save memory")
		return
	}
	defer f.Close()
	pprof.Lookup("heap").WriteTo(f, 1)

	log.Infof("end save memory")
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

func startMangoDB(roleIntf role.RoleInterface) *actor.PID {
	optiondb := db.NewOptionDbService(config.BtoConfig.Plugin.MongoDB.URL)
	if optiondb == nil {
		log.Errorf("Start optional db fail")
		return nil
	}

	pid := mongodb.NewMdbActor(roleIntf, optiondb)
	if pid == nil {
		log.Errorf("Start mongodb fail")
	}

	return pid
}