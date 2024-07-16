package wemcoin

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
)

const GenesisPrevHash = "0000000000000000000000000000000000000000000000000000000000000000"
const Reward = 6.25
const Difficulty = 5

var blockchain *Blockchain

func init() {
	var once sync.Once
	once.Do(func() {
		blockchain = &Blockchain{
			utxos: make(map[string]*TransactionOutput),
		}
	})
}

type Blockchain struct {
	utxos  map[string]*TransactionOutput
	blocks []*Block
}

func (b *Blockchain) Blocks() []*Block {
	return b.blocks
}

func (b *Blockchain) AddBlock(block *Block) {
	b.blocks = append(b.blocks, block)
}

func (b *Blockchain) Reset() {
	b.utxos = make(map[string]*TransactionOutput)
	b.blocks = make([]*Block, 0)
}

func (b *Blockchain) UTXOs() map[string]*TransactionOutput {
	return b.utxos
}

func (b *Blockchain) AddUTXO(utxo *TransactionOutput) {
	b.utxos[utxo.ID()] = utxo
}

func (b *Blockchain) RemoveUTXO(id string) {
	delete(b.utxos, id)
}

func BC() *Blockchain {
	return blockchain
}

type Block struct {
	id           uuid.UUID
	timestamp    int64
	prevHash     string
	transactions []*Transaction
	nonce        int
	hash         string
}

func (b *Block) String() string {
	var hash string
	if len(b.hash) > 0 {
		hash = b.hash
	} else {
		hash = "<none>"
	}
	return fmt.Sprintf(
		"Block(id = %s, prevHash = %s, hash = %s, timestamp = %d, nonce = %d, transactions = %v, merkleRoot = %s)\n",
		b.id,
		b.prevHash,
		hash,
		b.timestamp,
		b.nonce,
		b.transactions,
		b.MerkleRoot(),
	)
}

func (b *Block) Hash() string {
	return b.hash
}

func (b *Block) GenerateHash() {
	s := fmt.Sprintf(
		"%d%s%d%d%v%s",
		b.id,
		b.prevHash,
		b.timestamp,
		b.nonce,
		b.transactions,
		b.MerkleRoot(),
	)
	b.hash = SHA256hash(s)
}

func (b *Block) Nonce() int {
	return b.nonce
}

func (b *Block) IncrementNonce() {
	b.nonce += 1
}

func (b *Block) MerkleRoot() string {
	return NewMerkleTree(b.transactions).Root()
}

func (b *Block) Transactions() []*Transaction {
	return b.transactions
}

func (b *Block) AddTransaction(t *Transaction) bool {
	if b.prevHash != GenesisPrevHash {
		if !t.VerifyTransaction() {
			return false
		}
	}

	b.transactions = append(b.transactions, t)
	fmt.Printf("Transaction is valid and has been added to the block: %v", b)
	return true
}

func (b *Block) PrevHash() string {
	return b.prevHash
}

func (b *Block) Timestamp() int64 {
	return b.timestamp
}

func (b *Block) ID() uuid.UUID {
	return b.id
}

func NewBlock(prevHash string) *Block {
	return &Block{
		id:           uuid.New(),
		timestamp:    time.Now().UnixMilli(),
		prevHash:     prevHash,
		transactions: make([]*Transaction, 0),
	}
}

type Miner struct {
	reward float64
}

func (m *Miner) Reward() float64 {
	return m.reward
}

func (m *Miner) Mine(b *Block) {
	for !m.isGoldenHash(b) {
		b.IncrementNonce()
		b.GenerateHash()
	}

	fmt.Printf("%v has just been mined\n", b)
	blockchain.AddBlock(b)
	m.reward += Reward
}

func (m *Miner) isGoldenHash(b *Block) bool {
	prefix := strings.Repeat("0", Difficulty)
	return len(b.hash) > 0 && strings.HasPrefix(b.hash, prefix)
}

func NewMiner() *Miner {
	return &Miner{}
}
