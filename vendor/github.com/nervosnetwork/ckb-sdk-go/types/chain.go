package types

import (
	"math/big"

	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
)

type ScriptHashType string
type DepType string
type TransactionStatus string

const (
	HashTypeData ScriptHashType = "data"
	HashTypeType ScriptHashType = "type"

	DepTypeCode     DepType = "code"
	DepTypeDepGroup DepType = "dep_group"

	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusProposed  TransactionStatus = "proposed"
	TransactionStatusCommitted TransactionStatus = "committed"
)

type Epoch struct {
	CompactTarget uint64 `json:"compact_target"`
	Length        uint64 `json:"length"`
	Number        uint64 `json:"number"`
	StartNumber   uint64 `json:"start_number"`
}

type Header struct {
	CompactTarget    uint     `json:"compact_target"`
	Dao              Hash     `json:"dao"`
	Epoch            uint64   `json:"epoch"`
	Hash             Hash     `json:"hash"`
	Nonce            *big.Int `json:"nonce"`
	Number           uint64   `json:"number"`
	ParentHash       Hash     `json:"parent_hash"`
	ProposalsHash    Hash     `json:"proposals_hash"`
	Timestamp        uint64   `json:"timestamp"`
	TransactionsRoot Hash     `json:"transactions_root"`
	UnclesHash       Hash     `json:"uncles_hash"`
	Version          uint     `json:"version"`
}

type OutPoint struct {
	TxHash Hash `json:"tx_hash"`
	Index  uint `json:"index"`
}

type CellDep struct {
	OutPoint *OutPoint `json:"out_point"`
	DepType  DepType   `json:"dep_type"`
}

type Script struct {
	CodeHash Hash           `json:"code_hash"`
	HashType ScriptHashType `json:"hash_type"`
	Args     []byte         `json:"args"`
}

func (script *Script) OccupiedCapacity() uint64 {
	return uint64(len(script.Args)) + uint64(len(script.CodeHash.Bytes())) + 1
}

func (script *Script) Hash() (Hash, error) {
	data, err := script.Serialize()
	if err != nil {
		return Hash{}, err
	}

	hash, err := blake2b.Blake256(data)
	if err != nil {
		return Hash{}, err
	}

	return BytesToHash(hash), nil
}

func (script *Script) Equals(obj *Script) bool {
	if obj == nil {
		return false
	}

	sh, _ := script.Hash()
	oh, _ := obj.Hash()
	return sh.String() == oh.String()
}

type CellInput struct {
	Since          uint64    `json:"since"`
	PreviousOutput *OutPoint `json:"previous_output"`
}

type CellOutput struct {
	Capacity uint64  `json:"capacity"`
	Lock     *Script `json:"lock"`
	Type     *Script `json:"type"`
}

func (o CellOutput) OccupiedCapacity(outputData []byte) uint64 {
	occupiedCapacity := 8 + uint64(len(outputData)) + o.Lock.OccupiedCapacity()
	if o.Type != nil {
		occupiedCapacity += o.Type.OccupiedCapacity()
	}
	return occupiedCapacity
}

type Transaction struct {
	Version     uint          `json:"version"`
	Hash        Hash          `json:"hash"`
	CellDeps    []*CellDep    `json:"cell_deps"`
	HeaderDeps  []Hash        `json:"header_deps"`
	Inputs      []*CellInput  `json:"inputs"`
	Outputs     []*CellOutput `json:"outputs"`
	OutputsData [][]byte      `json:"outputs_data"`
	Witnesses   [][]byte      `json:"witnesses"`
}

func (t *Transaction) ComputeHash() (Hash, error) {
	data, err := t.Serialize()
	if err != nil {
		return Hash{}, err
	}

	hash, err := blake2b.Blake256(data)
	if err != nil {
		return Hash{}, err
	}

	return BytesToHash(hash), nil
}

func (t *Transaction) SizeInBlock() (uint64, error) {
	// raw tx serialize
	rawTxBytes, err := t.Serialize()
	if err != nil {
		return 0, err
	}

	var witnessBytes [][]byte
	for _, witness := range t.Witnesses {
		witnessBytes = append(witnessBytes, SerializeBytes(witness))
	}
	witnessesBytes := SerializeDynVec(witnessBytes)
	//tx serialize
	txBytes := SerializeTable([][]byte{rawTxBytes, witnessesBytes})
	txSize := uint64(len(txBytes))
	// tx offset cost
	txSize += 4
	return txSize, nil
}

func (t *Transaction) OutputsCapacity() (totalCapacity uint64) {
	for _, output := range t.Outputs {
		totalCapacity += output.Capacity
	}
	return
}

type WitnessArgs struct {
	Lock       []byte `json:"lock"`
	InputType  []byte `json:"input_type"`
	OutputType []byte `json:"output_type"`
}

type UncleBlock struct {
	Header    *Header  `json:"header"`
	Proposals []string `json:"proposals"`
}

type Block struct {
	Header       *Header        `json:"header"`
	Proposals    []string       `json:"proposals"`
	Transactions []*Transaction `json:"transactions"`
	Uncles       []*UncleBlock  `json:"uncles"`
}

type Cell struct {
	BlockHash     Hash      `json:"block_hash"`
	Capacity      uint64    `json:"capacity"`
	Lock          *Script   `json:"lock"`
	OutPoint      *OutPoint `json:"out_point"`
	Type          *Script   `json:"type"`
	Cellbase      bool      `json:"cellbase,omitempty"`
	OutputDataLen uint64    `json:"output_data_len,omitempty"`
}

type CellData struct {
	Content []byte `json:"content"`
	Hash    Hash   `json:"hash"`
}

type CellInfo struct {
	Data   *CellData   `json:"data"`
	Output *CellOutput `json:"output"`
}

type CellWithStatus struct {
	Cell   *CellInfo `json:"cell"`
	Status string    `json:"status"`
}

type TxStatus struct {
	BlockHash *Hash             `json:"block_hash"`
	Status    TransactionStatus `json:"status"`
}

type TransactionWithStatus struct {
	Transaction *Transaction `json:"transaction"`
	TxStatus    *TxStatus    `json:"tx_status"`
}

type BlockReward struct {
	Primary        *big.Int `json:"primary"`
	ProposalReward *big.Int `json:"proposal_reward"`
	Secondary      *big.Int `json:"secondary"`
	Total          *big.Int `json:"total"`
	TxFee          *big.Int `json:"tx_fee"`
}

type BlockEconomicState struct {
	Issuance    BlockIssuance `json:"issuance"`
	MinerReward MinerReward   `json:"miner_reward"`
	TxsFee      *big.Int      `json:"txs_fee"`
	FinalizedAt Hash          `json:"finalized_at"`
}

type BlockIssuance struct {
	Primary   *big.Int `json:"primary"`
	Secondary *big.Int `json:"secondary"`
}

type MinerReward struct {
	Primary   *big.Int `json:"primary"`
	Secondary *big.Int `json:"secondary"`
	Committed *big.Int `json:"committed"`
	Proposal  *big.Int `json:"proposal"`
}
