package rpc

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shaojunda/ckb-node-websocket-client/internal/model"
	"gorm.io/datatypes"
)

type NewTransactionSubscriptionResponse struct {
	JsonRPC string                           `json:"jsonrpc"`
	Method  string                           `json:"method"`
	Params  NewTransactionSubscriptionParams `json:"params"`
}

type NewTransactionSubscriptionParams struct {
	Result string `json:"result"`
}

type NewTransactionSubscriptionResult struct {
	Transaction PoolTransactionEntry `json:"transaction"`
}

type PoolTransactionEntry struct {
	CellDeps    []CellDep       `json:"cell_deps"`
	Hash        model.Hash      `json:"hash"`
	HeaderDeps  []model.Hash    `json:"header_deps"`
	Inputs      []CellInput     `json:"inputs"`
	Outputs     []CellOutput    `json:"outputs"`
	OutputsData []hexutil.Bytes `json:"outputs_data"`
	Version     hexutil.Uint    `json:"version"`
	Witnesses   []hexutil.Bytes `json:"witnesses"`
	Fee         hexutil.Uint64  `json:"fee"`
	Cycles      hexutil.Uint64  `json:"cycles"`
	Size        hexutil.Uint64  `json:"size"`
}

type CellDep struct {
	OutPoint OutPoint      `json:"out_point"`
	DepType  model.DepType `json:"dep_type"`
}

type OutPoint struct {
	TxHash model.Hash `json:"tx_hash"`
	Index  uint       `json:"index"`
}

type CellInput struct {
	Since          hexutil.Uint64 `json:"since"`
	PreviousOutput OutPoint       `json:"previous_output"`
}

type CellOutput struct {
	Capacity hexutil.Uint64 `json:"capacity"`
	Lock     *Script        `json:"lock"`
	Type     *Script        `json:"type"`
}

type Script struct {
	CodeHash model.Hash           `json:"code_hash"`
	HashType model.ScriptHashType `json:"hash_type"`
	Args     hexutil.Bytes        `json:"args"`
}

func (t PoolTransactionEntry) ToPoolTransactionEntryModel() (model.PoolTransactionEntry, error) {
	cellDeps, err := toCellDeps(t.CellDeps)
	if err != nil {
		return model.PoolTransactionEntry{}, err
	}
	return model.PoolTransactionEntry{
		CellDeps:       cellDeps,
		TxHash:         nil,
		HeaderDeps:     nil,
		Inputs:         nil,
		Outputs:        nil,
		OutputsData:    nil,
		Version:        0,
		Witnesses:      nil,
		TransactionFee: 0,
		BlockNumber:    0,
		BlockTimestamp: 0,
		Cycles:         0,
		TxSize:         0,
		DisplayInputs:  nil,
		DisplayOutputs: nil,
	}, nil
}

func toCellDeps(deps []CellDep) (datatypes.JSON, error) {
	result := make([]model.CellDep, len(deps))
	for i := 0; i < len(deps); i++ {
		dep := deps[i]
		result[i] = model.CellDep{
			OutPoint: model.OutPoint{
				TxHash: dep.OutPoint.TxHash,
				Index:  uint(dep.OutPoint.Index),
			},
			DepType: dep.DepType,
		}
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return []byte{}, err
	}
	return datatypes.JSON(bytes), nil
}
