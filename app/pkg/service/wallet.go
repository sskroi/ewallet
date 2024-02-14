package service

import (
	"ewallet/pkg/storage"
)

type WalletService struct {
	Storage storage.Storage
}

func NewWalletService(storage storage.Storage) *WalletService {
	return &WalletService{Storage: storage}
}

func (w *WalletService) Create() (*storage.Wallet, error) {
	return w.Storage.CreateWallet()
}

func (w *WalletService) Balance(id string) (float64, error) {
	exists, err := w.Storage.IsWalletExists(id)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, ErrWalletNoExist
	}

	wallet, err := w.Storage.Balance(id)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (w *WalletService) History(id string) ([]storage.Transfer, error) {
	exists, err := w.Storage.IsWalletExists(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrWalletNoExist
	}

	transfers, err := w.Storage.History(id)
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

func (w *WalletService) Transfer(fromId, toId string, amount float64) error {
	fromExist, err := w.Storage.IsWalletExists(fromId)
	if err != nil {
		return err
	}
	if !fromExist {
		return ErrOutWalletNoExist
	}

	inExist, err := w.Storage.IsWalletExists(toId)

	if err != nil {
		return err
	}
	if !inExist {
		return ErrInWalletNoExist
	}

	srcBalance, err := w.Balance(fromId)
	if err != nil {
		return err
	}

	if amount < 0 {
		return ErrNegativeTransAmount
	}

	if srcBalance-amount < 0 {
		return ErrNoBalanceForTrans
	}

	if err := w.Storage.Transfer(fromId, toId, amount); err != nil {
		return err
	}

	return nil
}
