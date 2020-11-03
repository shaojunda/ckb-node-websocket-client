package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	ckbUtils "github.com/nervosnetwork/ckb-sdk-go/utils"
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
	displayInputs, err := buildDisplayInputs(svc, entry)
	if err != nil {
		return err
	}
	poolTx.DisplayInputs = displayInputs
	displayOutputs, err := buildDisplayOutputs(svc, entry)
	poolTx.DisplayOutputs = displayOutputs
	return svc.dao.CreateOrUpdatePoolTransactionEntry(&poolTx)
}

func buildDisplayOutputs(svc Service, entry rpc.PoolTransactionEntry) (datatypes.JSON, error) {
	displayOutputs := make([]model.DisplayOutput, 0)
	for i, output := range entry.Outputs {
		var mode address.Mode
		if global.RPCSetting.Mode == "mainnet" {
			mode = address.Mainnet
		} else {
			mode = address.Testnet
		}
		addressHash, err := address.Generate(mode, &ckbTypes.Script{
			CodeHash: output.Lock.CodeHash,
			HashType: output.Lock.HashType,
			Args:     output.Lock.Args,
		})
		if err != nil {
			return datatypes.JSON{}, err
		}
		outputData := entry.OutputsData[i]
		cellOutput := ckbTypes.CellOutput{
			Capacity: uint64(output.Capacity),
			Lock: &ckbTypes.Script{
				CodeHash: output.Lock.CodeHash,
				HashType: output.Lock.HashType,
				Args:     output.Lock.Args,
			},
		}
		if output.Type != nil {
			cellOutput.Type = &ckbTypes.Script{
				CodeHash: cellOutput.Type.CodeHash,
				HashType: cellOutput.Type.HashType,
				Args:     cellOutput.Type.Args,
			}
		}
		cellType := getCellType(&cellOutput, outputData)
		displayOutput := model.DisplayOutput{
			Capacity:    cellOutput.Capacity,
			AddressHash: addressHash,
			Status:      "live",
			CellType:    cellType,
		}
		if cellType == "udt" {
			udtInfo, err := buildUdtInfo(svc, &cellOutput, outputData)
			if err != nil {
				return datatypes.JSON{}, err
			}
			displayOutput.UdtInfo = udtInfo
		}

		displayOutputs = append(displayOutputs, displayOutput)
	}

	displayOutputBytes, err := json.Marshal(displayOutputs)
	if err != nil {
		return datatypes.JSON{}, err
	}

	return datatypes.JSON(displayOutputBytes), nil
}

func buildDisplayInputs(svc Service, entry rpc.PoolTransactionEntry) (datatypes.JSON, error) {
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
			udtInfo, err := buildUdtInfo(svc, output, outputData)
			if err != nil {
				return datatypes.JSON{}, err
			}
			displayInput.UdtInfo = udtInfo
		}

		displayInputs = append(displayInputs, displayInput)
	}
	displayInputBytes, err := json.Marshal(displayInputs)
	if err != nil {
		return datatypes.JSON{}, err
	}

	return datatypes.JSON(displayInputBytes), nil
}

func buildUdtInfo(svc Service, output *ckbTypes.CellOutput, outputData []byte) (*model.UdtInfo, error) {
	typeHash, err := output.Type.Hash()
	if err != nil {
		return nil, err
	}
	udt, err := svc.dao.GetUdtByTypeHash(typeHash.String())
	if err != nil {
		return nil, err
	}
	udtAmount, err := ckbUtils.ParseSudtAmount(outputData)
	if err != nil {
		return nil, err
	}
	udtInfo := model.UdtInfo{
		Symbol:    udt.Symbol,
		Amount:    udtAmount.String(),
		Decimal:   strconv.FormatInt(int64(udt.Decimal), 10),
		TypeHash:  udt.TypeHash,
		Published: udt.Published,
	}
	return &udtInfo, nil
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
