package wemcoin

import (
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestGenesisBlockCreated(t *testing.T) {
	_, _, tearDown := testBlockchain()
	defer tearDown()

	genesisBlock := BC().Blocks()[0]

	if err := uuid.Validate(genesisBlock.ID().String()); err != nil {
		t.Errorf("genesisBlock missing ID")
	}

	if genesisBlock.Timestamp() == 0 {
		t.Errorf("genesisBlock missing timestamp")
	}

	if len(genesisBlock.Transactions()) != 1 {
		t.Errorf("genesisBlock transactions, got: %d", len(genesisBlock.Transactions()))
	}

	if !(genesisBlock.Nonce() > 0) {
		t.Errorf("genesisBlock nonce is not greater than 0, got: %d", genesisBlock.Nonce())
	}

	prefix := strings.Repeat("0", Difficulty)
	if !strings.HasPrefix(genesisBlock.Hash(), prefix) {
		t.Errorf("genesisBlock hash is invalid, got: %s", genesisBlock.Hash())
	}

	if len(genesisBlock.MerkleRoot()) < 0 {
		t.Errorf("genesisBlock merkleRoot is invalid, got: %s", genesisBlock.MerkleRoot())
	}
}

func testBlockchain() (*Wallet, *Wallet, func()) {
	const amount = 1500.0

	sender := NewWallet()
	receiver := NewWallet()
	lender := NewWallet()

	genesisTransaction := NewTransaction(lender.PublicKey(), sender.PublicKey(), amount)
	transactionOutput := NewTransactionOutput("0", genesisTransaction.Receiver(), amount)

	genesisTransaction.AddOutput(transactionOutput)
	genesisTransaction.GenerateSignature(lender.PrivateKey())

	bc := BC()
	bc.AddUTXO(transactionOutput)

	genesisBlock := NewBlock(GenesisPrevHash)
	genesisBlock.AddTransaction(genesisTransaction)

	miner := NewMiner()
	miner.Mine(genesisBlock)

	return sender, receiver, func() {
		bc.Reset()
	}
}

func TestMining(t *testing.T) {
	sender, receiver, tearDown := testBlockchain()
	defer tearDown()

	prevBlock := BC().Blocks()[0]
	miner := NewMiner()

	times := 5
	for i := 0; i < times; i++ {
		transaction := NewTransaction(sender.PublicKey(), receiver.PublicKey(), 1.0)
		transaction.GenerateSignature(sender.PrivateKey())
		block := NewBlock(prevBlock.Hash())
		block.AddTransaction(transaction)

		if err := uuid.Validate(block.ID().String()); err != nil {
			t.Errorf("block.ID() = %s", block.ID())
		}

		if block.Timestamp() == 0 {
			t.Errorf("block.Timestamp() missing timestamp")
		}

		if len(block.Transactions()) != 1 && block.Transactions()[0] == transaction {
			t.Errorf("block.Transactions() did not match")
		}

		if block.Nonce() != 0 {
			t.Errorf("block.Nonce() is not initialized to 0, got: %d", block.Nonce())
		}

		if len(block.MerkleRoot()) < 0 {
			t.Errorf("block.MerkleRoot() is invalid, got: %s", block.MerkleRoot())
		}

		miner.Mine(block)
		if got := block.Nonce(); got < 0 {
			t.Fatalf("block.Nonce() = %d", got)
		}

		prefix := strings.Repeat("0", Difficulty)
		if !strings.HasPrefix(block.Hash(), prefix) {
			t.Errorf("block.Hash() = %s", block.Hash())
		}

		prevBlock = block
	}

	// Account for the genesis block by adding/subtracting 1.
	blockCount := len(BC().Blocks())
	if blockCount != times+1 {
		t.Fatalf("blockCount = %d, want = %d", blockCount, times+1)
	}

	reward := float64(blockCount-1) * Reward
	if miner.Reward() != reward {
		t.Fatalf("miner.Reward() = %f, want = %f", miner.Reward(), reward)
	}
}

func TestTransferCoinsBetweenWallets(t *testing.T) {
	sender, receiver, tearDown := testBlockchain()
	defer tearDown()

	miner := NewMiner()

	assertBalance(t, sender, 1500.0)
	assertBalance(t, receiver, 0.0)

	amount := 100.0
	got := sender.Transfer(receiver.PublicKey(), amount)
	if got == nil {
		t.Fatalf("sender.Transfer(%v, %f) = %v", receiver.PublicKey(), amount, got)
	}

	blocks := BC().Blocks()
	b := NewBlock(blocks[len(blocks)-1].hash)
	b.AddTransaction(got)

	miner.Mine(b)

	assertBalance(t, sender, 1400.0)
	assertBalance(t, receiver, 100.0)

	if len(BC().Blocks()) != 2 {
		t.Fatalf("BC().Blocks() = %d, want = %d", len(BC().Blocks()), 2)
	}

	if len(BC().UTXOs()) != 2 {
		t.Fatalf("BC().UTXOs() = %d, want = %d", len(BC().UTXOs()), 2)
	}
}

func assertBalance(t *testing.T, w *Wallet, want float64) {
	if w.Balance() != want {
		t.Fatalf("w.Balance() = %f, want = %f", w.Balance(), want)
	}
}
