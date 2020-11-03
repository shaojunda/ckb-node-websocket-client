package service

import "github.com/shaojunda/ckb-node-websocket-client/internal/model"

func (svc Service) GetUdtByTypeHash(typeHash string) (*model.Udt, error) {
	return svc.dao.GetUdtByTypeHash(typeHash)
}
