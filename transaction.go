package wemcoin

import (
	"crypto/ecdsa"
	"fmt"
)

type Transaction struct {
	sender        *ecdsa.PublicKey
	receiver      *ecdsa.PublicKey
	amount        float64
	inputs        []*TransactionInput
	outputs       []*TransactionOutput
	transactionId string
	signature     []byte
}

func (t *Transaction) TransactionId() string {
	return t.transactionId
}

func (t *Transaction) Signature() []byte {
	return t.signature
}

func (t *Transaction) Sender() *ecdsa.PublicKey {
	return t.sender
}

func (t *Transaction) Receiver() *ecdsa.PublicKey {
	return t.receiver
}

func (t *Transaction) Amount() float64 {
	return t.amount
}

func (t *Transaction) Inputs() []*TransactionInput {
	return t.inputs
}

func (t *Transaction) AddInput(input *TransactionInput) {
	t.inputs = append(t.inputs, input)
}

func (t *Transaction) Outputs() []*TransactionOutput {
	return t.outputs
}

func (t *Transaction) AddOutput(output *TransactionOutput) {
	t.outputs = append(t.outputs, output)
}

func (t *Transaction) GenerateSignature(p *ecdsa.PrivateKey) {
	h := fmt.Sprintf("%v%v%f", t.sender, t.receiver, t.amount)
	t.signature = Sign(p, h)
}

func (t *Transaction) VerifyTransaction() bool {
	if !t.verifySignature(t.sender) {
		return false
	}

	for _, input := range t.inputs {
		input.UTXO = BC().UTXOs()[input.TransactionOutputId]
	}

	t.outputs = append(t.outputs, NewTransactionOutput(t.transactionId, t.receiver, t.amount))
	t.outputs = append(t.outputs, NewTransactionOutput(t.transactionId, t.sender, t.inputsSum()-t.amount))

	for _, output := range t.outputs {
		BC().AddUTXO(output)
	}

	for _, input := range t.inputs {
		BC().RemoveUTXO(input.UTXO.ID())
	}

	return true
}

func (t *Transaction) verifySignature(pub *ecdsa.PublicKey) bool {
	h := fmt.Sprintf("%v%v%f", t.sender, t.receiver, t.amount)
	return Verify(pub, h, t.signature)
}

func (t *Transaction) inputsSum() float64 {
	var sum float64
	for _, inputs := range t.inputs {
		if inputs.UTXO == nil {
			continue
		}
		sum += inputs.UTXO.Amount()
	}
	return sum
}

func NewTransaction(
	sender *ecdsa.PublicKey,
	receiver *ecdsa.PublicKey,
	amount float64,
) *Transaction {
	transactionId := fmt.Sprintf("%v%v%f", sender, receiver, amount)
	return &Transaction{
		sender:        sender,
		receiver:      receiver,
		amount:        amount,
		inputs:        make([]*TransactionInput, 0),
		outputs:       make([]*TransactionOutput, 0),
		transactionId: SHA256hash(transactionId),
	}
}

type TransactionInput struct {
	TransactionOutputId string
	UTXO                *TransactionOutput
}

func NewTransactionInput(transactionOutputId string, UTXO *TransactionOutput) *TransactionInput {
	return &TransactionInput{TransactionOutputId: transactionOutputId, UTXO: UTXO}
}

type TransactionOutput struct {
	parentTransactionId string
	receiver            *ecdsa.PublicKey
	amount              float64
	id                  string
}

func (t *TransactionOutput) Amount() float64 {
	return t.amount
}

func (t *TransactionOutput) Receiver() *ecdsa.PublicKey {
	return t.receiver
}

func (t *TransactionOutput) ParentTransactionId() string {
	return t.parentTransactionId
}

func (t *TransactionOutput) ID() string {
	return t.id
}

func (t *TransactionOutput) IsMine(receiver *ecdsa.PublicKey) bool {
	return t.receiver.Equal(receiver)
}

func NewTransactionOutput(
	parentTransactionId string,
	receiver *ecdsa.PublicKey,
	amount float64,
) *TransactionOutput {
	id := fmt.Sprintf("%s%f%s", receiver, amount, parentTransactionId)
	return &TransactionOutput{
		parentTransactionId: parentTransactionId,
		receiver:            receiver,
		amount:              amount,
		id:                  SHA256hash(id),
	}
}
