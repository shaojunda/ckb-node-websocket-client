package rpc

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
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
	Fee         hexutil.Uint64       `json:"fee"`
	Cycles      hexutil.Uint64       `json:"cycles"`
	Size        hexutil.Uint64       `json:"size"`
}

type PoolTransactionEntry struct {
	CellDeps    []CellDep       `json:"cell_deps"`
	Hash        ckbTypes.Hash   `json:"hash"`
	HeaderDeps  []ckbTypes.Hash `json:"header_deps"`
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
	TxHash ckbTypes.Hash `json:"tx_hash"`
	Index  hexutil.Uint  `json:"index"`
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
	CodeHash ckbTypes.Hash        `json:"code_hash"`
	HashType model.ScriptHashType `json:"hash_type"`
	Args     hexutil.Bytes        `json:"args"`
}

func (t PoolTransactionEntry) ToPoolTransactionEntryModel() (model.PoolTransactionEntry, error) {
	cellDeps, err := t.toCellDeps()
	if err != nil {
		return model.PoolTransactionEntry{}, err
	}
	headerDeps, err := t.toHeaderDeps()
	if err != nil {
		return model.PoolTransactionEntry{}, err
	}
	inputs, err := t.toInputs()
	if err != nil {
		return model.PoolTransactionEntry{}, err
	}
	outputs, err := t.toOutputs()
	if err != nil {
		return model.PoolTransactionEntry{}, err
	}
	outputsData, err := t.toOutputsData()
	if err != nil {
		return model.PoolTransactionEntry{}, nil
	}
	witnesses, err := t.toWitnesses()
	if err != nil {
		return model.PoolTransactionEntry{}, nil
	}
	return model.PoolTransactionEntry{
		CellDeps:       cellDeps,
		TxHash:         t.Hash.String(),
		HeaderDeps:     headerDeps,
		Inputs:         inputs,
		Outputs:        outputs,
		OutputsData:    outputsData,
		Version:        uint(t.Version),
		Witnesses:      witnesses,
		TransactionFee: uint64(t.Fee),
		BlockNumber:    0,
		BlockTimestamp: 0,
		Cycles:         uint64(t.Cycles),
		TxSize:         uint64(t.Size),
		DisplayInputs:  nil,
		DisplayOutputs: nil,
	}, nil
}

func (t PoolTransactionEntry) toHeaderDeps() (datatypes.JSON, error) {
	bytes, err := json.Marshal(t.HeaderDeps)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return datatypes.JSON(bytes), nil
}

func (t PoolTransactionEntry) toCellDeps() (datatypes.JSON, error) {
	result := make([]model.CellDep, len(t.CellDeps))
	for i := 0; i < len(t.CellDeps); i++ {
		dep := t.CellDeps[i]
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
		return datatypes.JSON{}, err
	}

	return datatypes.JSON(bytes), nil
}

func (t PoolTransactionEntry) toInputs() (datatypes.JSON, error) {
	result := make([]model.CellInput, len(t.Inputs))
	for i := 0; i < len(t.Inputs); i++ {
		input := t.Inputs[i]
		result[i] = model.CellInput{
			Since: uint64(input.Since),
			PreviousOutput: model.OutPoint{
				TxHash: input.PreviousOutput.TxHash,
				Index:  uint(input.PreviousOutput.Index),
			},
		}
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return datatypes.JSON(bytes), nil
}
func (t PoolTransactionEntry) toOutputs() (datatypes.JSON, error) {
	result := make([]model.CellOutput, len(t.Outputs))
	for i := 0; i < len(t.Outputs); i++ {
		output := t.Outputs[i]
		result[i] = model.CellOutput{
			Capacity: uint64(output.Capacity),
			Lock: &model.Script{
				CodeHash: output.Lock.CodeHash,
				HashType: output.Lock.HashType,
				Args:     output.Lock.Args.String(),
			},
		}
		if output.Type != nil {
			result[i].Type = &model.Script{
				CodeHash: output.Type.CodeHash,
				HashType: output.Type.HashType,
				Args:     output.Type.Args.String(),
			}
		}
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return datatypes.JSON(bytes), nil
}

func (t PoolTransactionEntry) toOutputsData() (datatypes.JSON, error) {
	bytes, err := json.Marshal(t.OutputsData)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return datatypes.JSON(bytes), nil
}

func (t PoolTransactionEntry) toWitnesses() (datatypes.JSON, error) {
	bytes, err := json.Marshal(t.Witnesses)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return datatypes.JSON(bytes), nil
}
