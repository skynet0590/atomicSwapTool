package main


import (
	"context"
	"fmt"
	_ "github.com/skynet0590/atomicSwapTool/client/asset/btc" // register btc asset
	_ "github.com/skynet0590/atomicSwapTool/client/asset/dcr" // register dcr asset
	_ "github.com/skynet0590/atomicSwapTool/client/asset/ltc" // register ltc asset
	"github.com/skynet0590/atomicSwapTool/client/cmd/astc/version"
	"github.com/skynet0590/atomicSwapTool/client/core"
	"github.com/skynet0590/atomicSwapTool/dex"
	"github.com/skynet0590/atomicSwapTool/dex/encode"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	fmt.Sprintf("%v %v", appCtx,cancel)
	// Parse configuration.
	cfg, err := configure()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configration error: %v\n", err)
		os.Exit(1)
	}

	// Initialize logging.
	var utc = !cfg.LocalLogs
	if cfg.Net == dex.Simnet {
		utc = false
	}
	logMaker := initLogging(cfg.DebugLevel, utc)
	log = logMaker.Logger("DEXC")
	log.Infof("%s version %v (Go version %s)", version.AppName, version.Version(), runtime.Version())
	if utc {
		log.Infof("Logging with UTC time stamps. Current local time is %v",
			time.Now().Local().Format("15:04:05 MST"))
	}

	// Prepare the Core.
	clientCore, err := core.New(&core.Config{
		DBPath:       cfg.DBPath, // global set in config.go
		Net:          cfg.Net,
		Logger:       logMaker.Logger("CORE"),
		TorProxy:     cfg.TorProxy,
		TorIsolation: cfg.TorIsolation,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating client core: %v\n", err)
		os.Exit(1)
	}

	// Catch interrupt signal (e.g. ctrl+c), prompting to shutdown if the user
	// is logged in, and there are active orders or matches.
	killChan := make(chan os.Signal)
	signal.Notify(killChan, os.Interrupt)
	go func() {
		for range killChan {
			if clientCore.PromptShutdown() {
				cancel()
				return
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		clientCore.Run(appCtx)
		cancel() // in the event that Run returns prematurely prior to context cancellation
		wg.Done()
	}()

	<-clientCore.Ready()
	appPW := encode.PassBytes("vietanh123")
	appPW.MarshalJSON()
	clientCore.InitializeClient(appPW)

	walletPW := encode.PassBytes("123456")

	err = clientCore.CreateWallet(appPW, walletPW, &core.WalletForm{
		AssetID: 0,
		Config: map[string]string{
			"fallbackfee": "0.001",
			"redeemconftarget": "2",
			"rpcbind": "127.0.0.1",
			"rpcpassword": "VietAnhDepTrai",
			"rpcport": "18332",
			"rpcuser": "alice",
			"txsplit": "0",
			"walletname": "skynet",
		},
	})

	log.Error("Create Wallet Error: ", err)

	err = clientCore.ConnectWallet(0)

	log.Error("Connect Wallet Error: ", err)

	wb, err := clientCore.AssetBalance(0)
	fmt.Println(wb.Available, err)
	clientCore.Trade()
}
