package explorer

// storm db keeps last transactions of an address
type Transaction struct {
	Hash        string `json:"hash" storm:"index"`
	From        string `json:"from" storm:"id"`
	To          string `json:"to" storm:"index"`
	Gas         int64  `json:"gas" storm:"index"`
	GasUsed     int64  `json:"gas_used" storm:"index"`
	GasPrice    int64  `json:"gas_price" storm:"index"`
	BlockNumber int64  `json:"block_number" storm:"index"`
	Timestamp   int64  `json:"timestamp" storm:"index"`
}

type SortableTransactions []Transaction

func (s SortableTransactions) Len() int {
	return len(s)
}

func (s SortableTransactions) Less(i, j int) bool {
	return s[i].Timestamp < s[j].Timestamp
}

func (s SortableTransactions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type Explorer interface {
	NewTransactions(last Transaction) []Transaction // sorted in descending order by timestamp
}
