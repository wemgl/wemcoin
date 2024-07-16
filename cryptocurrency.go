package wemcoin

import (
	"crypto/ecdsa"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) Balance() float64 {
	var bal float64
	for _, output := range BC().UTXOs() {
		if !output.IsMine(w.publicKey) {
			continue
		}
		bal += output.Amount()
	}
	return bal
}

func (w *Wallet) Transfer(receiver *ecdsa.PublicKey, amount float64) *Transaction {
	inputs := make([]*TransactionInput, 0)

	for _, output := range BC().UTXOs() {
		if !output.IsMine(w.publicKey) {
			continue
		}
		inputs = append(inputs, NewTransactionInput(output.ID(), output))
	}

	t := NewTransaction(
		w.publicKey,
		receiver,
		amount,
	)

	for _, input := range inputs {
		t.AddInput(input)
	}

	t.GenerateSignature(w.privateKey)
	return t
}

func NewWallet() *Wallet {
	keyPair := GenKey()
	return &Wallet{
		privateKey: keyPair,
		publicKey:  &keyPair.PublicKey,
	}
}
