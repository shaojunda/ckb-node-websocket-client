package dao

import "gorm.io/gorm"

type Dao struct {
	engine *gorm.DB
}

func NewDao(engine *gorm.DB) *Dao {
	return &Dao{engine: engine}
}
