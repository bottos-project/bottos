package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

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

func main() {
	var err error

	blockDb := db.NewDbService(config.Param.DataDir, filepath.Join(config.Param.DataDir, "blockchain"))

	//	fmt.Println("init account")
	//	account.CreateAccountManager()

	_, err = chain.CreateBlockChain(blockDb)
	if err != nil {
		fmt.Println("Create BlockChain error: ", err)
		os.Exit(1)
	}

	cactor.InitActors()
	caapi.PushTransaction(2876568)

	WaitSystemDown()

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
