package service

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shaojunda/ckb-node-websocket-client/internal/rpc"
)

func (svc *Service) CreatePoolTransactionEntry(entry rpc.PoolTransactionEntry, fee hexutil.Uint64, cycles hexutil.Uint64, size hexutil.Uint64) error {
	entry.Fee = fee
	entry.Cycles = cycles
	entry.Size = size
	poolTx, err := entry.ToPoolTransactionEntryModel()
	if err != nil {
		return err
	}
	_, err = svc.dao.CreatePoolTransactionEntry(&poolTx)
	if err != nil {
		return err
	}
	return nil
}

func (svc Service) CreateOrUpdatePoolTransactionEntry(entry rpc.PoolTransactionEntry, fee hexutil.Uint64, cycles hexutil.Uint64, size hexutil.Uint64) error {
	entry.Fee = fee
	entry.Cycles = cycles
	entry.Size = size
	poolTx, err := entry.ToPoolTransactionEntryModel()
	if err != nil {
		return err
	}
	return svc.dao.CreateOrUpdatePoolTransactionEntry(&poolTx)
}
