package main

import (
	"fmt"
	//	"os"
	//	"time"
	//	//"time"

	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/db"
	//	"github.com/bottos-project/core/account"
	//	"github.com/bottos-project/core/api"
	//	"github.com/bottos-project/core/common"
	//	"github.com/bottos-project/core/common/types"

	//	pro "github.com/bottos-project/core/producer"
	//	//"github.com/bottos-project/core/p2p"

	//	"github.com/micro/go-micro"
	//	log "github.com/sirupsen/logrus"
	cactor "github.com/bottos-project/core/action/actor"
	caapi "github.com/bottos-project/core/action/actor/api"
)

var (
	DataDir = "./datadir/"
)

func main() {
	//	fmt.Println("init db")

	blockDb := db.NewDbService(config.Param.DataDir, config.Param.DataDir)

	//	fmt.Println("init account")
	//	account.CreateAccountManager()

	fmt.Println("init blockchain")
	chain.CreateBlockChain(blockDb)

	cactor.InitActors()
	caapi.PushTransaction(2876568)

	//console.ReadLine()
	//	fmt.Println("init txpool")
	//	txpool, _ := tr.CreateTxPool(&emux, bc)

	//	fmt.Println("init block producer")
	//	producer := pro.NewProducer(&emux, bc)

	//	fmt.Println("init p2p")
	//	//proto := p2p.NewProtocol(&emux, bc)

	//	fmt.Println("init done \n\n")

	//	go txpool.TxPoolLoop()
	//	go producer.ProducerLoop()

	//	// test
	//	go func() {
	//		for {
	//			txpool.Add(&types.Transaction{Id: "test", AccountName: "testname"})
	//			time.Sleep(1000 * time.Millisecond)
	//		}
	//	}()

	//	for {

	//	}

	//	log.SetOutput(os.Stdout)
	//	log.SetLevel(log.DebugLevel)

	//	svc := micro.NewService(
	//		micro.Name("core"),
	//		micro.RegisterTTL(30),
	//		micro.RegisterInterval(1000),
	//		micro.Version(""),
	//	)
	//	svc.Init()
	//	repo := core.NewCoreSrvice(txpool)
	//	core.RegisterCoreHandler(svc.Server(), repo)
	//	fmt.Println("fmt")
	//	if err := svc.Run(); err != nil {
	//		panic(err)
	//	}
}
