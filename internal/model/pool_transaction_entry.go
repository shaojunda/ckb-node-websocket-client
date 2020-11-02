package model

import "gorm.io/datatypes"

type PoolTransactionEntry struct {
	*Model
	CellDeps       datatypes.JSON `json:"cell_deps"`
	TxHash         []byte         `json:"tx_hash"`
	HeaderDeps     datatypes.JSON `json:"header_deps"`
	Inputs         datatypes.JSON `json:"inputs"`
	Outputs        datatypes.JSON `json:"outputs"`
	OutputsData    datatypes.JSON `json:"outputs_data"`
	Version        uint           `json:"version"`
	Witnesses      datatypes.JSON `json:"witnesses"`
	TransactionFee uint64         `json:"transaction_fee"`
	BlockNumber    uint64         `json:"block_number"`
	BlockTimestamp uint64         `json:"block_timestamp"`
	Cycles         uint64         `json:"cycles"`
	TxSize         uint64         `json:"tx_size"`
	DisplayInputs  datatypes.JSON `json:"display_inputs"`
	DisplayOutputs datatypes.JSON `json:"display_outputs"`
}
