package dao

import (
	"github.com/shaojunda/ckb-node-websocket-client/internal/model"
)

func (d Dao) CreatePoolTransactionEntry(param *model.PoolTransactionEntry) (*model.PoolTransactionEntry, error) {
	poolTx := model.PoolTransactionEntry{
		CellDeps:       param.CellDeps,
		TxHash:         param.TxHash,
		HeaderDeps:     param.HeaderDeps,
		Inputs:         param.Inputs,
		Outputs:        param.Outputs,
		OutputsData:    param.OutputsData,
		Version:        param.Version,
		Witnesses:      param.Witnesses,
		TransactionFee: param.TransactionFee,
		BlockNumber:    param.BlockNumber,
		BlockTimestamp: param.BlockTimestamp,
		Cycles:         param.Cycles,
		TxSize:         param.TxSize,
		DisplayInputs:  param.DisplayInputs,
		DisplayOutputs: param.DisplayOutputs,
	}
	return poolTx.Create(d.engine)
}

func (d Dao) CreateOrUpdatePoolTransactionEntry(param *model.PoolTransactionEntry) error {
	poolTx := model.PoolTransactionEntry{
		CellDeps:       param.CellDeps,
		TxHash:         param.TxHash,
		HeaderDeps:     param.HeaderDeps,
		Inputs:         param.Inputs,
		Outputs:        param.Outputs,
		OutputsData:    param.OutputsData,
		Version:        param.Version,
		Witnesses:      param.Witnesses,
		TransactionFee: param.TransactionFee,
		BlockNumber:    param.BlockNumber,
		BlockTimestamp: param.BlockTimestamp,
		Cycles:         param.Cycles,
		TxSize:         param.TxSize,
		DisplayInputs:  param.DisplayInputs,
		DisplayOutputs: param.DisplayOutputs,
	}
	return poolTx.Upsert(d.engine)
}
