package server

import (
	"github.com/nocturna-ta/golib/database/nosql/clickhouse"
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
	configLocation, _ := cmd.Flags().GetString("config")
	cfg := &config.MainConfig{}
	config.ReadConfig(cfg, configLocation)

	database, err := clickhouse.New(&clickhouse.Config{
		Addrs: cfg.ClickHouse.Addrs,
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouse.Auth.Database,
			Username: cfg.ClickHouse.Auth.Username,
			Password: cfg.ClickHouse.Auth.Password,
		},
		Database:        cfg.ClickHouse.Database,
		DialTimeout:     cfg.ClickHouse.DialTimeout,
		MaxOpenConns:    cfg.ClickHouse.MaxOpenConns,
		MaxIdleConns:    cfg.ClickHouse.MaxIdleConns,
		ConnMaxLifetime: cfg.ClickHouse.ConnMaxLifetime,
		TLS: &clickhouse.TLSConfig{
			Enable:             cfg.ClickHouse.TLS.Enable,
			InsecureSkipVerify: cfg.ClickHouse.TLS.InsecureSkipVerify,
			CertFile:           cfg.ClickHouse.TLS.CertFile,
			KeyFile:            cfg.ClickHouse.TLS.KeyFile,
			CAFile:             cfg.ClickHouse.TLS.CAFile,
		},
		BlockBufferSize:    cfg.ClickHouse.BlockBufferSize,
		MaxCompressionSize: cfg.ClickHouse.MaxCompressionSize,
		AsyncInsert:        cfg.ClickHouse.AsyncInsert,
		AsyncInsertOptions: clickhouse.AsyncInsertOptions{
			MaxBatchSize: cfg.ClickHouse.AsyncInsertOptions.MaxBatchSize,
			MaxDelay:     cfg.ClickHouse.AsyncInsertOptions.MaxDelay,
		},
		Debug: cfg.ClickHouse.Debug,
	})

	if err != nil {
		return err
	}

	log.Info("ClickHouse connections established successfully")

	//client, err := ethereum.GetEthereumClient(&cfg.Blockchain)
	//if err != nil {
	//	return err
	//}
	//
	//defer client.Close()
	//
	//publisher, err := kafka.NewPublisher(context.Background(), cfg.Kafka.Producer)
	//if err != nil {
	//	log.Fatalf("Failed to instantiate kafka publisher: %w", err)
	//	return err
	//}

	appContainer := newContainer(&options{
		Cfg: cfg,
		DB:  database,
		//Client:    client,
		//DB:        database,
		//Publisher: publisher,
	})

	server := api.New(&api.Options{
		Cfg:        appContainer.Cfg,
		VoteResult: appContainer.VoteResultUc,
	})

	go server.Run()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Info("Exiting gracefully...")
	case err := <-server.ListenError():
		log.Error("Error starting web server, exiting gracefully:", err)
	}

	return nil
}
