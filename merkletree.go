package wemcoin

type MerkleTree struct {
	transactions []*Transaction
	root         string
}

func (m *MerkleTree) Root() string {
	return m.root
}

func NewMerkleTree(transactions []*Transaction) *MerkleTree {
	m := &MerkleTree{transactions: transactions}
	m.construct()
	return m
}

func (m *MerkleTree) construct() {
	if len(m.transactions) == 0 {
		m.root = ""
		return
	}

	var hashTransactionIDs func(transactionIDs []string) []string

	hashTransactionIDs = func(transactionIDs []string) []string {
		if len(transactionIDs) == 1 {
			return transactionIDs
		}

		hashes := make([]string, 0)
		for i := 0; i < len(transactionIDs); i += 2 {
			hash := m.mergeHashes(transactionIDs[i], transactionIDs[i+1])
			hashes = append(hashes, hash)
		}

		if len(transactionIDs)%2 == 1 {
			last := len(transactionIDs) - 1
			hash := m.mergeHashes(transactionIDs[last], transactionIDs[last])
			hashes = append(hashes, hash)
		}

		return hashTransactionIDs(hashes)
	}

	transactionIDs := make([]string, len(m.transactions))
	for _, transaction := range m.transactions {
		transactionIDs = append(transactionIDs, transaction.TransactionId())
	}

	root := hashTransactionIDs(transactionIDs)
	m.root = root[0]
}

func (m *MerkleTree) mergeHashes(h1, h2 string) string {
	return SHA256hash(h1 + h2)
}
