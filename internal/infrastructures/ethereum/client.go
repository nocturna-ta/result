package ethereum

import (
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/result/config"
)

func GetEthereumClient(cfg *config.BlockchainConfig) (ethereum.Client, error) {
	return ethereum.New(&ethereum.Options{
		URL: cfg.GanacheURL,
	})
}
