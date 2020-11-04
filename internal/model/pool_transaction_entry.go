package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PoolTransactionEntry struct {
	*Model
	CellDeps       datatypes.JSON `json:"cell_deps"`
	TxHash         string         `json:"tx_hash"`
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

type UdtInfo struct {
	Symbol    string `json:"symbol"`
	Amount    string `json:"amount"`
	Decimal   string `json:"decimal"`
	TypeHash  string `json:"type_hash"`
	Published bool   `json:"published"`
}

type CellInfo struct {
	Lock *Script `json:"lock,omitempty"`
	Type *Script `json:"type,omitempty"`
	Data string  `json:"data"`
}

type DisplayInput struct {
	FromCellbase    bool      `json:"from_cellbase"`
	Capacity        uint64    `json:"capacity"`
	AddressHash     string    `json:"address_hash"`
	GeneratedTxHash string    `json:"generated_tx_hash"`
	CellIndex       string    `json:"cell_index"`
	CellType        string    `json:"cell_type"`
	CellInfo        *CellInfo `json:"cell_info,omitempty"`
	UdtInfo         *UdtInfo  `json:"udt_info,omitempty"`
}

type DisplayOutput struct {
	Capacity       uint64    `json:"capacity"`
	AddressHash    string    `json:"address_hash"`
	Status         string    `json:"status"`
	ConsumedTxHash string    `json:"consumed_tx_hash"`
	CellType       string    `json:"cell_type"`
	CellInfo       *CellInfo `json:"cell_info,omitempty"`
	UdtInfo        *UdtInfo  `json:"udt_info,omitempty"`
}

func (t PoolTransactionEntry) Create(db *gorm.DB) (*PoolTransactionEntry, error) {
	if err := db.Create(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (t PoolTransactionEntry) Get(db *gorm.DB) (PoolTransactionEntry, error) {
	var poolTx PoolTransactionEntry
	db = db.Where("id = ?", t.ID)
	err := db.First(&poolTx).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return poolTx, err
	}
	return poolTx, nil
}

func (t PoolTransactionEntry) GetByTxHash(db *gorm.DB) (PoolTransactionEntry, error) {
	var poolTx PoolTransactionEntry
	db = db.Where("tx_hash = ?", t.TxHash)
	err := db.First(&poolTx).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return poolTx, err
	}
	return poolTx, nil
}

func (t PoolTransactionEntry) Upsert(db *gorm.DB) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "tx_hash"}},
		DoUpdates: clause.AssignmentColumns([]string{"cell_deps", "header_deps", "inputs", "outputs", "outputs_data", "version",
			"witnesses", "transaction_fee", "block_number", "block_timestamp", "cycles", "tx_size", "display_inputs", "display_outputs", "updated_at"}),
	}).Create(&t).Error
}
