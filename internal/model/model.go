package model

import (
	"fmt"
	"github.com/shaojunda/ckb-node-websocket-client/pkg/setting"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
