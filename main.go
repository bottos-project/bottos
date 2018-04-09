package main

import (
	"time"
	"fmt"
	"os"
	"path/filepath"
	//"time"
	//	"github.com/bottos-project/bottos/core/account"
	"github.com/bottos-project/bottos/core/api"
	"github.com/bottos-project/bottos/core/common"
	"github.com/bottos-project/bottos/core/common/types"
	"github.com/bottos-project/bottos/core/event"
	tr "github.com/bottos-project/bottos/core/trx"
	pro "github.com/bottos-project/bottos/core/producer"
	"github.com/bottos-project/bottos/core/account"
	"github.com/bottos-project/bottos/core/db"
	//"github.com/bottos-project/bottos/core/p2p"

	"github.com/micro/go-micro"
	log "github.com/sirupsen/logrus"
)

var (
	DataDir = "./datadir/"
)

/*
	// subscribe a NewMinedBlockEvent
	sub := emux.Subscribe(common.NewMinedBlockEvent{})
	go func (minedBlockSub event.Subscription) {
		for obj := range minedBlockSub.Chan() {
			switch ev := obj.(type) {
			case common.NewMinedBlockEvent:
				fmt.Printf("ProcessBlockLoop : recv new mined block, %d\n", ev.Block.Number())
				fmt.Printf("\n")
			}
		}
	}(sub)
	*/

func main() {
	fmt.Println("init db")
	
	blockDb, err := db.NewKVDatabase(filepath.Join(DataDir, "blockchain"))
	if err != nil {
		fmt.Println("init kv database error")
		return
	}

	fmt.Println("init eventmux")
	var emux event.TypeMux

	fmt.Println("init account")
	account.CreateAccountManager()

	fmt.Println("init blockchain")
	bc, _ := common.CreateBlockChain(blockDb, &emux)

	fmt.Println("init txpool")
	txpool, _:= tr.CreateTxPool(&emux, bc)

	fmt.Println("init block producer")
	producer := pro.NewProducer(&emux, bc)

	fmt.Println("init p2p")
	//proto := p2p.NewProtocol(&emux, bc)

	fmt.Println("init done \n\n")

	go txpool.TxPoolLoop()
	go producer.ProducerLoop()

	// test
	go func() {
		for {
			txpool.Add(&types.Transaction{Id: "test", AccountName: "testname"})
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	
	
	for {

	}

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	svc := micro.NewService(
		micro.Name("core"),
		micro.RegisterTTL(30),
		micro.RegisterInterval(1000),
		micro.Version(""),
	)
	svc.Init()
	repo := core.NewCoreSrvice(txpool)
	core.RegisterCoreHandler(svc.Server(), repo)
	fmt.Println("fmt")
	if err := svc.Run(); err != nil {
		panic(err)
	}
}
