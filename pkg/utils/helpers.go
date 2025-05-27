package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/result/config"
)

func StringToTx(signedTx string) (*types.Transaction, error) {
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(common.FromHex(signedTx)); err != nil {
		return nil, err
	}

	return tx, nil
}

func ConvertToRouterCorsConfig(configCors *config.CorsConfig) *router.CorsConfig {
	return &router.CorsConfig{
		AllowOrigins:     configCors.AllowOrigins,
		AllowMethods:     configCors.AllowMethods,
		AllowHeaders:     configCors.AllowHeaders,
		AllowCredentials: configCors.AllowCredentials,
		ExposeHeaders:    configCors.ExposeHeaders,
		MaxAge:           configCors.MaxAge,
	}
}
