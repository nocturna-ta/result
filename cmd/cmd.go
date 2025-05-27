package cmd

import (
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/result/cmd/server"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "Result  Service",
		Short: "Result Service",
	}
)

func Execute() {
	log.SetFormatter("json")
	rootCmd.AddCommand(server.ServeHttpCmd())
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error: ", err.Error())
		os.Exit(-1)
	}
}
