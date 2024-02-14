package storage

type Storage interface {
	CreateWallet() (*Wallet, error)
	Balance(id string) (*Wallet, error)
	Transfer(fromId, toId string, amount float64) error
	History(id string) ([]Transfer, error)
	IsWalletExists(id string) (bool, error)
}

type Wallet struct {
	Id      string  `db:"id" json:"id"`
	Balance float64 `db:"balance" json:"balance"`
}

type Transfer struct {
	Time   string  `db:"time" json:"time"`
	FromId string  `db:"fromId" json:"from"`
	ToId   string  `db:"toId" json:"to"`
	Amount float64 `db:"amount" json:"amount"`
}
