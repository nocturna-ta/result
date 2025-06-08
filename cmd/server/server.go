package server

import (
	"context"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/nocturna-ta/golib/database/sql"

	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/result/config"
	"github.com/nocturna-ta/result/internal/handler/api"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var (
	serverHTTPCmd = &cobra.Command{
		Use:   "serve-http",
		Short: "Result Service HTTP",
		Long:  "Result Service HTTP",
		RunE:  run,
	}
)

func ServeHttpCmd() *cobra.Command {
	serverHTTPCmd.Flags().StringP("config", "c", "", "Config Path, both relative or absolute. i.e: /usr/local/bin/config/files")
	return serverHTTPCmd
}

func run(cmd *cobra.Command, args []string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configLocation, _ := cmd.Flags().GetString("config")
	cfg := &config.MainConfig{}
	config.ReadConfig(cfg, configLocation)

	database := sql.New(sql.DBConfig{
		SlaveDSN:        cfg.Database.SlaveDSN,
		MasterDSN:       cfg.Database.MasterDSN,
		RetryInterval:   cfg.Database.RetryInterval,
		MaxIdleConn:     cfg.Database.MaxIdleConn,
		MaxConn:         cfg.Database.MaxConn,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}, sql.DriverClickHouse)

	//client, err := ethereum.GetEthereumClient(&cfg.Blockchain)
	//if err != nil {
	//	return err
	//}
	//
	//defer client.Close()

	appContainer := newContainer(&options{
		Cfg: cfg,
		DB:  database,
		Ctx: ctx,
		//Client:    client,
	})

	server := api.New(&api.Options{
		Cfg:          appContainer.Cfg,
		VoteResult:   appContainer.VoteResultUc,
		LiveResult:   appContainer.LiveResultUc,
		WebsocketHub: appContainer.WebSocketHub,
	})

	go server.Run()

	log.WithFields(log.Fields{
		"port":         cfg.Server.Port,
		"websocket":    true,
		"live_results": true,
	}).Info("Result Service started with WebSocket support")

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Info("Exiting gracefully...")
		cancel()
	case err := <-server.ListenError():
		log.Error("Error starting web server, exiting gracefully:", err)
		cancel()
	}

	appContainer.WebSocketHub.Stop()

	return nil
}
