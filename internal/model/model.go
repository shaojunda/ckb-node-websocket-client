package model

import (
	"fmt"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/shaojunda/ckb-node-websocket-client/pkg/setting"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	var dsn string
	if databaseSetting.DBType == "postgresql" {
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=%s",
			databaseSetting.UserName, databaseSetting.Password,
			databaseSetting.DBName, databaseSetting.Port, databaseSetting.SSLMode,
		)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}

type CellDep struct {
	OutPoint OutPoint         `json:"out_point"`
	DepType  ckbTypes.DepType `json:"dep_type"`
}

type OutPoint struct {
	TxHash ckbTypes.Hash `json:"tx_hash"`
	Index  uint          `json:"index"`
}

type Script struct {
	CodeHash ckbTypes.Hash           `json:"code_hash"`
	HashType ckbTypes.ScriptHashType `json:"hash_type"`
	Args     string                  `json:"args"`
}

type CellInput struct {
	Since          uint64   `json:"since"`
	PreviousOutput OutPoint `json:"previous_output"`
}

type CellOutput struct {
	Capacity uint64  `json:"capacity"`
	Lock     *Script `json:"lock"`
	Type     *Script `json:"type"`
}
