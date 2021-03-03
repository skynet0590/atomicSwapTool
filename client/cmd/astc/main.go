// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package main

import (
	"context"
	"fmt"
	gioApp "gioui.org/app"
	"github.com/skynet0590/atomicSwapTool/client/gui"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	_ "github.com/skynet0590/atomicSwapTool/client/asset/btc" // register btc asset
	_ "github.com/skynet0590/atomicSwapTool/client/asset/dcr" // register dcr asset
	_ "github.com/skynet0590/atomicSwapTool/client/asset/ltc" // register ltc asset
	"github.com/skynet0590/atomicSwapTool/client/cmd/astc/version"
	"github.com/skynet0590/atomicSwapTool/client/core"
	"github.com/skynet0590/atomicSwapTool/dex"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())

	// Parse configuration.
	cfg, err := configure()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configration error: %v\n", err)
		os.Exit(1)
	}

	// Initialize logging.
	utc := !cfg.LocalLogs
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

/*done:
	wg.Wait()
	log.Info("Exiting dexc main.")
	closeFileLogger()*/
	w := gui.NewWindow(clientCore)

	// This creates a new application window and starts the UI.
	go func() {
		if err := w.Loop(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	// Starts Gio main.
	gioApp.Main()
}
