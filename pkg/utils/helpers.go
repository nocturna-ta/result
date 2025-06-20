package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func StringToTx(signedTx string) (*types.Transaction, error) {
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(common.FromHex(signedTx)); err != nil {
		return nil, err
	}

	return tx, nil
}
