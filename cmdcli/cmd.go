/*
@Time : 2018/7/27 15:36 
@Author : 星空之钥丶
@File : cmd
@Software: GoLang
*/
package cmdcli

import(
	cli "gopkg.in/urfave/cli.v1"
	"os"
	"github.com/bottos-project/bottos/config"
	"io/ioutil"
	"bytes"
	"fmt"
	"encoding/json"
	"strings"
)
var Conf  config.Parameter
var GenConf config.GenesisConfig
var KeyPair config.KeyPair

func Init() (config.Parameter, config.GenesisConfig, error) {

	config.InitParam(&Conf, &GenConf)

	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "config",
			Value:"./chainconfig.json",
			Usage: "config file path the greeting,If without this path, the bottos process will boot up with default config in hardcode",
		},
		cli.StringFlag{
			Name: "genesis",
			Value: Conf.GenesisJson,
			Usage: "genesis config file path the greeting",
		},
		cli.StringFlag{
			Name: "datadir",
			Value: Conf.DataDir,
			Usage: "datadir's path",
		},
		cli.BoolFlag{
			Name: "disable-api",
			Usage: "disable restful api's requests",
		},
		cli.IntFlag{
			Name: "apiport",
			Value: Conf.APIPort,
			Usage: "api service port for the greeting",
		},
		cli.BoolFlag{
			Name: "disable-rpc",
			Usage: "disable rpc requests",
		},
		cli.IntFlag{
			Name: "rpcport",
			Value: 8690,
			Usage: "json-rpc port for the greeting",
		},
		cli.IntFlag{
			Name: "p2pport",
			Value: Conf.P2PPort,
			Usage: "local listen on this p2p port to receive remote p2p messages",
		},
		cli.StringFlag{
			Name: "servaddr",
			Value: Conf.ServAddr,
			Usage: "for p2p sync / reply local server ip& port info",
		},
		cli.StringFlag{
			Name: "peerlist",
			Value: "",
			Usage: "for p2p add pne / add neighbour. Example: 192.168.1.2:9868, 192.168.1.3:9868, 192.168.1.4:9868",
		},
		cli.StringFlag{
			Name: "delegate-signkey",
			Usage: "--delegate-signkey=<pubkey>,<private key>.Param struct needs be modified ,public and private key for native contract, external contracts' accounts",
		},
		cli.StringFlag{
			Name: "delegate",
			Usage: "Assign one producer. Later this section will no more be used.\n Only one delegate is allowed in one node(other than bottos account).",
		},
		cli.BoolFlag{
			Name: "enable-stale-report",
			Usage: "",
		},
		cli.BoolFlag{
			Name: "enable-mongodb",
			Usage: "",
		},
		cli.StringFlag{
			Name: "mongodb",
			Value: Conf.OptionDb,
			Usage: "db inst for load mongodb",
		},
		cli.StringFlag{
			Name: "logconfig",
			Value: Conf.LogConfig,
			Usage: "for seelog config",
		},
	}

	app.Action = func(c *cli.Context) error {
		var ChaincfgExists bool
		var GenesiscfgExists bool

		_, err := os.Stat(c.String("config"))
		if err != nil && os.IsNotExist(err) {
			fmt.Println("'" + c.String("config") + "' file does not exist.")
			ChaincfgExists = false
		} else if err != nil {
			fmt.Println("Read config file status error: ", err)
			return err
		} else {
			ChaincfgExists = true
		}
	
		_, err = os.Stat(c.String("genesis"))
		if err != nil && os.IsNotExist(err) {
			fmt.Println("'" + c.String("genesis") + "' file does not exist.")
			GenesiscfgExists = false
		} else if err != nil {
			fmt.Println("Read config file status error: ", err)
			return err
		} else {
			GenesiscfgExists = true
		}

		if ChaincfgExists == true {
			file, e := loadConfigJson(c.String("config"))
			if e != nil {
				fmt.Println("Read config file error: ", e)
				return e
			}

			e = json.Unmarshal(file, &Conf)
			if e != nil {
				fmt.Println("Unmarshal config file error: ", e)
				return e
			}
		}
		
		if GenesiscfgExists == true {
			file, e := loadConfigJson(c.String("genesis"))
			if e != nil {
				fmt.Println("Read genesis config file error: ", e)
				return e
			}

			e = json.Unmarshal(file, &GenConf)
			if e != nil {
				fmt.Println("Unmarshal config file error: ", e)
				return e
			}
		}
		
		if len(c.String("datadir")) > 0 {
			Conf.DataDir = c.String("datadir")
		}
		
		if len(c.String("logconfig")) > 0 {
			Conf.LogConfig = c.String("logconfig")
		}
		
		if len(c.String("servaddr")) > 0 {
			Conf.ServAddr = c.String("servaddr")
		}

		if c.Int("apiport") > 0 {
			//For new restful api port
			Conf.APIPort = c.Int("apiport")
		}
		
		if c.Int("rpcport") > 0 {
			//TO DO: for micro rpc port. The port is not be used by now.
		}

		if c.Int("p2pport") > 0 {
			Conf.P2PPort = c.Int("p2pport")
		}

		if len(c.String("peerlist")) > 0 {
			var strval string
			peer_list := c.String("peerlist")
			
			strval = strings.Replace(peer_list, " ", "", -1)
			Conf.PeerList = strings.Split(strval, ",")
			if len(Conf.PeerList[0]) <= 0 {
				Conf.PeerList = Conf.PeerList[1:]
			}
			if len(Conf.PeerList[len(Conf.PeerList) - 1]) <= 0 {
				Conf.PeerList = Conf.PeerList[:len(Conf.PeerList) - 1]
			}
		}

		if len(c.String("delegate-signkey")) > 0 {
			strval := strings.Replace(c.String("delegate-signkey"), " ", "", -1)
			key := strings.Split(strval, ",")
			KeyPair.PrivateKey = key[0]
			KeyPair.PublicKey = key[1]
			Conf.DelegateSignKey = KeyPair
		}

		if len(c.String("delegate")) > 0 {
			Conf.Delegates = []string{c.String("delegate")}
		}
		
		if len(c.String("mongodb")) > 0 {
			Conf.OptionDb = c.String("mongodb")
		}
		
		if c.GlobalIsSet("enable-stale-report") {
			fmt.Println(c.String("enable-stale-report"))
			Conf.EnableStaleReport = true
		}

		if ! c.GlobalIsSet("enable-mongodb") {
			Conf.OptionDb = ""
		}
		
		if c.GlobalIsSet("disable-api") {
			//TODO for new restful api
		}
		
		if c.GlobalIsSet("disable-rpc") {
			Conf.RpcServiceEnable = false
		}

		return nil
	}
	err := app.Run(os.Args)
	return Conf, GenConf, err
}

func loadConfigJson(fn string) ([]byte, error) {
	file, e := ioutil.ReadFile(fn)
	if e != nil {
		return nil, e
	}

	// Remove the UTF-8 Byte Order Mark
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	return file, nil
}
