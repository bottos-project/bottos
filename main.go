package main

//	"fmt"
//	"os"
//	"path/filepath"
//	"time"
//	//"time"
//	//	"github.com/bottos-project/core/account"
//	"github.com/bottos-project/core/account"
//	"github.com/bottos-project/core/api"
//	"github.com/bottos-project/core/common"
//	"github.com/bottos-project/core/common/types"
//	"github.com/bottos-project/core/db"
//	pro "github.com/bottos-project/core/producer"
//	//"github.com/bottos-project/core/p2p"

//	"github.com/micro/go-micro"
//	log "github.com/sirupsen/logrus"

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
	//	fmt.Println("init db")

	//	blockDb, err := db.NewKVDatabase(filepath.Join(DataDir, "blockchain"))
	//	if err != nil {
	//		fmt.Println("init kv database error")
	//		return
	//	}

	//	fmt.Println("init eventmux")
	//	var emux event.TypeMux

	//	fmt.Println("init account")
	//	account.CreateAccountManager()

	//	fmt.Println("init blockchain")
	//	bc, _ := common.CreateBlockChain(blockDb, &emux)

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
