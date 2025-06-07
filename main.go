package main

import (
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/nocturna-ta/result/cmd"
)

// @title Result Service
// @version 1.0.0
// @description Result Service.
// @BasePath /
func main() {
	cmd.Execute()
}
