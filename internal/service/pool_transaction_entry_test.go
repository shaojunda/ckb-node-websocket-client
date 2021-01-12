package service

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/shaojunda/ckb-node-websocket-client/global"
	"github.com/shaojunda/ckb-node-websocket-client/pkg/setting"
	"testing"
)

func TestGetCellTypeReturnNormal(t *testing.T) {
	args, err := hexutil.Decode("0x3954acece65096bfa81258983ddb83915fc56bd8")
	if err != nil {
		t.Error("failed")
	}
	lock := &ckbTypes.Script{
		CodeHash: ckbTypes.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
		HashType: "type",
		Args:     args,
	}
	output := &ckbTypes.CellOutput{
		Capacity: 0,
		Lock:     lock,
		Type:     nil,
	}
	cellType := getCellType(output, []byte{})
	if cellType != "normal" {
		t.Errorf("got %s want %s", cellType, "normal")
	}
}

func TestGetCellTypeReturnNervosDaoDeposit(t *testing.T) {
	global.SystemCodeHash = &setting.SystemCodeHashS{
		Dao:   "0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e",
		Sudts: []string{""},
	}
	args, err := hexutil.Decode("0x3954acece65096bfa81258983ddb83915fc56bd8")
	if err != nil {
		t.Error("failed")
	}
	lock := &ckbTypes.Script{
		CodeHash: ckbTypes.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
		HashType: "type",
		Args:     args,
	}
	typeScript := &ckbTypes.Script{
		CodeHash: ckbTypes.HexToHash("0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e"),
		HashType: "type",
		Args:     args,
	}
	output := &ckbTypes.CellOutput{
		Capacity: 0,
		Lock:     lock,
		Type:     typeScript,
	}
	cellType := getCellType(output, make([]byte, 8))
	if cellType != "nervos_dao_deposit" {
		t.Errorf("got %s want %s", cellType, "nervos_dao_deposit")
	}
}

func TestGetCellTypeReturnNervosDaoWithdrawing(t *testing.T) {
	global.SystemCodeHash = &setting.SystemCodeHashS{
		Dao:   "0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e",
		Sudts: []string{""},
	}
	args, err := hexutil.Decode("0x3954acece65096bfa81258983ddb83915fc56bd8")
	if err != nil {
		t.Error("failed")
	}
	lock := &ckbTypes.Script{
		CodeHash: ckbTypes.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
		HashType: "type",
		Args:     args,
	}
	typeScript := &ckbTypes.Script{
		CodeHash: ckbTypes.HexToHash("0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e"),
		HashType: "type",
		Args:     args,
	}
	output := &ckbTypes.CellOutput{
		Capacity: 0,
		Lock:     lock,
		Type:     typeScript,
	}
	cellType := getCellType(output, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9})
	if cellType != "nervos_dao_withdrawing" {
		t.Errorf("got %s want %s", cellType, "nervos_dao_withdrawing")
	}
}

func TestGetCellTypeReturnUdt(t *testing.T) {
	global.SystemCodeHash = &setting.SystemCodeHashS{
		Dao:   "0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e",
		Sudts: []string{"0x5e7a36a77e68eecc013dfa2fe6a23f3b6c344b04005808694ae6dd45eea4cfd5"},
	}
	args, err := hexutil.Decode("0x3954acece65096bfa81258983ddb83915fc56bd8")
	if err != nil {
		t.Error("failed")
	}
	lock := &ckbTypes.Script{
		CodeHash: ckbTypes.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
		HashType: "type",
		Args:     args,
	}
	typeScript := &ckbTypes.Script{
		CodeHash: ckbTypes.HexToHash("0x5e7a36a77e68eecc013dfa2fe6a23f3b6c344b04005808694ae6dd45eea4cfd5"),
		HashType: "data",
		Args:     args,
	}
	output := &ckbTypes.CellOutput{
		Capacity: 0,
		Lock:     lock,
		Type:     typeScript,
	}
	cellType := getCellType(output, make([]byte, 61))
	if cellType != "udt" {
		t.Errorf("got %s want %s", cellType, "udt")
	}
}

func TestGetCellTypeReturnNormalWhenDataIsInvalid(t *testing.T) {
	global.SystemCodeHash = &setting.SystemCodeHashS{
		Dao:   "0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e",
		Sudts: []string{"0x5e7a36a77e68eecc013dfa2fe6a23f3b6c344b04005808694ae6dd45eea4cfd5"},
	}
	args, err := hexutil.Decode("0x3954acece65096bfa81258983ddb83915fc56bd8")
	if err != nil {
		t.Error("failed")
	}
	lock := &ckbTypes.Script{
		CodeHash: ckbTypes.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
		HashType: "type",
		Args:     args,
	}
	typeScript := &ckbTypes.Script{
		CodeHash: ckbTypes.HexToHash("0x5e7a36a77e68eecc013dfa2fe6a23f3b6c344b04005808694ae6dd45eea4cfd5"),
		HashType: "data",
		Args:     args,
	}
	output := &ckbTypes.CellOutput{
		Capacity: 0,
		Lock:     lock,
		Type:     typeScript,
	}
	cellType := getCellType(output, make([]byte, 10))
	if cellType != "normal" {
		t.Errorf("got %s want %s", cellType, "normal")
	}
}
