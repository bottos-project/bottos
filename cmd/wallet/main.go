package main

import (
	"fmt"
	"github.com/bottos-project/bottos/config"
	comtool "github.com/bottos-project/bottos/restful/common"
	"github.com/bottos-project/bottos/restful/wallet"
	log "github.com/cihub/seelog"
	"gopkg.in/urfave/cli.v1"
	"net"
	"net/http"
	"golang.org/x/net/netutil"
	"os"
	"sort"
	"strconv"
	"os/signal"
	"syscall"
)

func init() {
	logger, err := log.LoggerFromConfigAsFile("./walletlog.xml")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer logger.Flush()
	log.ReplaceLogger(logger)
	
        /*logger2, err := log.LoggerFromConfigAsFile("./walletlog_corefile.xml")
        if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
        }
	defer logger2.Flush()
        
        common.CoreLogger = logger2
	*/
}

func main() {
	signal.Ignore(syscall.SIGHUP)
	//var language string
	app := cli.NewApp()
	app.Name = "wallet"
	app.Usage = "the bottos wallet command line interface"
	app.Version = "3.2.0"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 6869,
			Usage: "listening port",
		},
		cli.StringFlag{
			Name:  "wallet-servaddr,w",
			Value: "127.0.0.1",
			Usage: "wallet service listen address",
			//Destination: &language,
		},
	}

	/*app.Commands = []cli.Command{
		{
			Name:     "port",
			Aliases:  []string{"a"},
			Usage:    "calc 1+1",
			Category: "arithmetic",
			Action: func(c *cli.Context) error {
				fmt.Println("1 + 1 = ", 1+1)
				return nil
			},
		},
		{
			Name:     "sub",
			Aliases:  []string{"s"},
			Usage:    "calc 5-3",
			Category: "arithmetic",
			Action: func(c *cli.Context) error {
				fmt.Println("5 - 3 = ", 5-3)
				return nil
			},
		},
		{
			Name:     "db",
			Usage:    "database operations",
			Category: "database",
			Subcommands: []cli.Command{
				{
					Name:  "insert",
					Usage: "insert data",
					Action: func(c *cli.Context) error {
						fmt.Println("insert subcommand")
						return nil
					},
				},
				{
					Name:  "delete",
					Usage: "delete data",
					Action: func(c *cli.Context) error {
						fmt.Println("delete subcommand")
						return nil
					},
				},
			},
		},
	}*/
	app.Action = startWallet
	/*app.Before = func(c *cli.Context) error {
		fmt.Println("app Before")
		return nil
	}*/
	app.After = func(c *cli.Context) error {
		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "show help!",
	}

	/*	cli.VersionFlag = cli.BoolFlag{
			Name:  "print-version, v",
			Usage: "print version",
		}*/

	err := app.Run(os.Args)
	if err != nil {
		log.Critical(err)
	}
}
func startWallet(c *cli.Context) {
	fmt.Println("Start wallet REST service. Listen IP:", c.String("wallet-servaddr") ," Port:",c.Int("port"))
	log.Info("Start wallet REST service. Listen IP:", c.String("wallet-servaddr") ," Port:",c.Int("port"))

	comtool.VerifyInit()
	listener, err := net.Listen("tcp", c.String("wallet-servaddr")+":"+strconv.Itoa(c.Int("port")))
	if err != nil {
		log.Errorf("Listen: %v, wallet start failed", err)
		fmt.Printf("Listen: %v, wallet start failed\n", err)
		os.Exit(1)
	}
	defer listener.Close()
	listener = netutil.LimitListener(listener, config.WALLET_REST_LIMINT_VALUE)

	router := wallet.NewRouter()
	//config.InitConfig()
	//err := http.ListenAndServe(config.BtoConfig.Plugin.Wallet.WalletRESTHost+":"+strconv.Itoa(config.BtoConfig.Plugin.Wallet.WalletRESTPort), router)
	err = http.Serve(listener, router)

	if err != nil {
		fmt.Println("Start wallet REST service Failed: ", err)
		log.Critical("Start wallet REST service Failed: ", err)
		os.Exit(1)
	}
}
