package main

import (
	"context"
	"github.com/skynet0590/atomicSwapTool/server/admin"
	_ "decred.org/dcrdex/server/asset/btc" // register btc asset
	_ "decred.org/dcrdex/server/asset/dcr" // register dcr asset
	_ "decred.org/dcrdex/server/asset/ltc" // register ltc asset
	"fmt"
	"github.com/decred/dcrd/dcrec/secp256k1/v3"
	"github.com/skynet0590/atomicSwapTool/dex"
	"github.com/skynet0590/atomicSwapTool/dex/encode"
	dexsrv "github.com/skynet0590/atomicSwapTool/server/dex"
	"os"
	"runtime"
	"strings"
	"sync"
)

func main() {
	// Create a context that is canceled when a shutdown request is received
	// via requestShutdown.
	ctx := withShutdownCancel(context.Background())
	// Listen for both interrupt signals (e.g. CTRL+C) and shutdown requests
	// (requestShutdown calls).
	go shutdownListener()

	err := mainCore(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func mainCore(ctx context.Context) error {
	// Parse the configuration file, and setup logger.
	cfg, _, err := loadConfig()
	if err != nil {
		fmt.Printf("Failed to load server config: %s\n", err.Error())
		return err
	}
	defer func() {
		if logRotator != nil {
			logRotator.Close()
		}
	}()

	// Display app version.
	log.Infof("%s version %v (Go version %s)", AppName, Version(), runtime.Version())
	log.Infof("dcrdex starting for network: %s", cfg.Network)
	log.Infof("swap locktimes config: maker %s, taker %s",
		dex.LockTimeMaker(cfg.Network), dex.LockTimeTaker(cfg.Network))

	// Load the market and asset configurations for the given network.
	markets, assets, err := loadMarketConfFile(cfg.Network, cfg.MarketsConfPath)
	if err != nil {
		return fmt.Errorf("failed to load market and asset config %q: %v",
			cfg.MarketsConfPath, err)
	}

	if err != nil {
		return fmt.Errorf("failed to load market and asset config %q: %v",
			cfg.MarketsConfPath, err)
	}
	log.Infof("Found %d assets, loaded %d markets, for network %s",
		len(assets), len(markets), strings.ToUpper(cfg.Network.String()))
	// NOTE: If MaxUserCancelsPerEpoch is ultimately a setting we want to keep,
	// bake it into the markets.json file and load it per-market in settings.go.
	// For now, patch it into each dex.MarketInfo.
	for _, mkt := range markets {
		mkt.MaxUserCancelsPerEpoch = cfg.MaxUserCancels
	}


	// Load, or create and save, the DEX signing key.
	var privKey *secp256k1.PrivateKey
	if len(cfg.SigningKeyPW) == 0 {
		// cfg.SigningKeyPW, err = admin.PasswordPrompt(ctx, "Signing key password: ")
		//if err != nil {
		//	return fmt.Errorf("cannot use password: %v", err)
		//}
	}
	privKey, err = dexKey(cfg.DEXPrivKeyPath, cfg.SigningKeyPW)
	encode.ClearBytes(cfg.SigningKeyPW)
	if err != nil {
		return err
	}

	// Create the DEX manager.
	dexConf := &dexsrv.DexConf{
		DataDir:    cfg.DataDir,
		LogBackend: cfg.LogMaker,
		Markets:    markets,
		Assets:     assets,
		Network:    cfg.Network,
		DBConf: &dexsrv.DBConf{
			DBName:       cfg.DBName,
			Host:         cfg.DBHost,
			User:         cfg.DBUser,
			Port:         cfg.DBPort,
			Pass:         cfg.DBPass,
			ShowPGConfig: cfg.ShowPGConfig,
		},
		RegFeeXPub:        cfg.RegFeeXPub,
		RegFeeAmount:      cfg.RegFeeAmount,
		RegFeeConfirms:    cfg.RegFeeConfirms,
		BroadcastTimeout:  cfg.BroadcastTimeout,
		CancelThreshold:   cfg.CancelThreshold,
		Anarchy:           cfg.Anarchy,
		FreeCancels:       cfg.FreeCancels,
		BanScore:          cfg.BanScore,
		InitTakerLotLimit: cfg.InitTakerLotLimit,
		AbsTakerLotLimit:  cfg.AbsTakerLotLimit,
		DEXPrivKey:        privKey,
		CommsCfg: &dexsrv.RPCConfig{
			RPCCert:        cfg.RPCCert,
			RPCKey:         cfg.RPCKey,
			ListenAddrs:    cfg.RPCListen,
			AltDNSNames:    cfg.AltDNSNames,
			DisableDataAPI: cfg.DisableDataAPI,
		},
		NoResumeSwaps: cfg.NoResumeSwaps,
	}
	dexMan, err := dexsrv.NewDEX(dexConf)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	if cfg.AdminSrvOn {
		srvCFG := &admin.SrvConfig{
			Core:    dexMan,
			Addr:    cfg.AdminSrvAddr,
			// AuthSHA: adminSrvAuthSHA,
			Cert:    cfg.RPCCert,
			Key:     cfg.RPCKey,
		}
		adminServer, err := admin.NewServer(srvCFG)
		if err != nil {
			return fmt.Errorf("cannot set up admin server: %v", err)
		}
		wg.Add(1)
		go func() {
			adminServer.Run(ctx)
			wg.Done()
		}()
	}

	log.Info("The Server is running. Hit CTRL+C to quit...")
	<-ctx.Done()
	// Wait for the admin server to finish.
	wg.Wait()

	log.Info("Stopping Server...")
	dexMan.Stop()
	log.Info("Bye!")

	return nil
}
