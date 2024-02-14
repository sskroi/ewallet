package service

import (
	"errors"
	"ewallet/pkg/storage"
)

type Wallet interface {
	Create() (*storage.Wallet, error)
	Transfer(fromId, toId string, amount float64) error
	History(id string) ([]storage.Transfer, error)
	Balance(id string) (float64, error)
}

type Service struct {
	Wallet
}

var (
	ErrWalletNoExist       = errors.New("the wallet with the specified id does not exist")
	ErrOutWalletNoExist    = errors.New("the outgoing wallet does not exist")
	ErrInWalletNoExist     = errors.New("the incoming wallet does not exist")
	ErrNoBalanceForTrans   = errors.New("there are not enough funds on the outgoing wallet")
	ErrNegativeTransAmount = errors.New("it is not possible to transfer a negative amount")
)

func NewService(storage storage.Storage) *Service {
	return &Service{
		Wallet: NewWalletService(storage),
	}
}
