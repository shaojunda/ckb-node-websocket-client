package service

import (
	"context"
	"github.com/shaojunda/ckb-node-websocket-client/global"
	"github.com/shaojunda/ckb-node-websocket-client/internal/dao"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.NewDao(global.DBEngine)
	return svc
}
