package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/shaojunda/ckb-node-websocket-client/global"
	"github.com/shaojunda/ckb-node-websocket-client/internal/model"
	"github.com/shaojunda/ckb-node-websocket-client/internal/rpc"
	"gorm.io/datatypes"
	"strconv"
)

const MinSudtAmountByteSize = 16

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
	displayInputs, err := buildDisplayInputs(entry)
	if err != nil {
		return err
	}

	poolTx.DisplayInputs = displayInputs
	return svc.dao.CreateOrUpdatePoolTransactionEntry(&poolTx)
}

func buildDisplayInputs(entry rpc.PoolTransactionEntry) (datatypes.JSON, error) {
	displayInputs := make([]model.DisplayInput, 0)
	for _, input := range entry.Inputs {
		previousOutput := input.PreviousOutput
		tx, err := global.CKBClient.GetTransaction(context.Background(), previousOutput.TxHash)
		if err != nil {
			return datatypes.JSON{}, err
		}
		output := tx.Transaction.Outputs[previousOutput.Index]
		var mode address.Mode
		if global.RPCSetting.Mode == "mainnet" {
			mode = address.Mainnet
		} else {
			mode = address.Testnet
		}

		addressHash, err := address.Generate(mode, output.Lock)
		if err != nil {
			return datatypes.JSON{}, err
		}
		cellIndex := strconv.FormatUint(uint64(previousOutput.Index), 10)
		outputData := tx.Transaction.OutputsData[previousOutput.Index]
		cellType := getCellType(output, outputData)
		displayInput := model.DisplayInput{
			FromCellbase:    false,
			Capacity:        output.Capacity,
			AddressHash:     addressHash,
			GeneratedTxHash: previousOutput.TxHash.String(),
			CellIndex:       cellIndex,
			CellType:        cellType,
			CellInfo:        buildCellInfo(output, outputData),
		}
		if cellType == "udt" {
			displayInput.UdtInfo = buildUdtInfo(output)
		}

		displayInputs = append(displayInputs, displayInput)
	}
	displayInputBytes, err := json.Marshal(displayInputs)
	if err != nil {
		return datatypes.JSON{}, err
	}

	return datatypes.JSON(displayInputBytes), nil
}

func buildUdtInfo(output *ckbTypes.CellOutput) *model.UdtInfo {
	return nil
}

func buildCellInfo(output *ckbTypes.CellOutput, outputData []byte) *model.CellInfo {
	cellInfo := &model.CellInfo{
		Lock: &model.Script{
			CodeHash: output.Lock.CodeHash,
			HashType: output.Lock.HashType,
			Args:     hexutil.Encode(output.Lock.Args),
		},
		Data: hexutil.Encode(outputData),
	}

	if output.Type != nil {
		cellInfo.Type = &model.Script{
			CodeHash: output.Type.CodeHash,
			HashType: output.Type.HashType,
			Args:     hexutil.Encode(output.Type.Args),
		}
	}
	return cellInfo
}

func getCellType(output *ckbTypes.CellOutput, outputData []byte) string {
	if output.Type == nil {
		return "normal"
	}
	switch output.Type.CodeHash.String() {
	case global.SystemCodeHash.Dao:
		if bytes.Compare(outputData, make([]byte, 8)) == 0 {
			return "nervos_dao_deposit"
		} else {
			return "nervos_dao_withdrawing"
		}
	case global.SystemCodeHash.Sudt:
		if len(outputData) >= MinSudtAmountByteSize {
			return "udt"
		} else {
			return "normal"
		}
	default:
		return "normal"
	}
}
