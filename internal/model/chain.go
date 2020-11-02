package model

type DepType string
type ScriptHashType string

type CellDep struct {
	OutPoint OutPoint `json:"out_point"`
	DepType  DepType  `json:"dep_type"`
}

type OutPoint struct {
	TxHash Hash `json:"tx_hash"`
	Index  uint `json:"index"`
}

type Script struct {
	CodeHash Hash           `json:"code_hash"`
	HashType ScriptHashType `json:"hash_type"`
	Args     []byte         `json:"args"`
}

type CellInput struct {
	Since          uint64   `json:"since"`
	PreviousOutput OutPoint `json:"previous_output"`
}

type CellOutput struct {
	Capacity uint64 `json:"capacity"`
	Lock     Script `json:"lock"`
	Type     Script `json:"type"`
}
