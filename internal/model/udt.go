package model

import "gorm.io/gorm"

type Udt struct {
	*Model
	CodeHash             string `json:"code_hash"`
	HashType             string `json:"hash_type"`
	Args                 string `json:"args"`
	TypeHash             string `json:"type_hash"`
	FullName             string `json:"full_name"`
	Symbol               string `json:"symbol"`
	Decimal              int32  `json:"decimal"`
	Description          string `json:"description"`
	IconFile             string `json:"icon_file"`
	OperatorWebsite      string `json:"operator_website"`
	AddressCount         uint64 `json:"address_count"`
	TotalAmount          uint64 `json:"total_amount"`
	UdtType              int32  `json:"udt_type"`
	Published            bool   `json:"published"`
	BlockTimestamp       uint64 `json:"block_timestamp"`
	IssuerAddress        string `json:"issuer_address"`
	CkbTransactionsCount uint64 `json:"ckb_transactions_count"`
}

func (u Udt) GetByTypeHash(db *gorm.DB) (*Udt, error) {
	var udt Udt
	db = db.Where("type_hash = ?", u.TypeHash)
	err := db.First(&udt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &udt, err
	}
	return &udt, nil
}
