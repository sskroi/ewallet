package postgres

import (
	"ewallet/pkg/storage"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type PostgresDB struct {
	db *sqlx.DB
}

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func New(cfg Config) (*PostgresDB, error) {
	const fn = "postgres.New"

	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &PostgresDB{db: db}, nil
}

func (d *PostgresDB) CreateWallet() (*storage.Wallet, error) {
	const fn = "postgres.CreateWallet"

	res := d.db.QueryRow(`INSERT INTO wallets DEFAULT VALUES RETURNING id, balance;`)

	if res.Err() != nil {
		return nil, fmt.Errorf("%s: %w", fn, res.Err())
	}

	wallet := &storage.Wallet{}

	if err := res.Scan(&wallet.Id, &wallet.Balance); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return wallet, nil
}

func (d *PostgresDB) Balance(id string) (*storage.Wallet, error) {
	const fn = "postgres.Balance"

	res := d.db.QueryRow(`SELECT id, balance FROM wallets WHERE id = $1;`, id)

	if res.Err() != nil {
		return nil, fmt.Errorf("%s: %w", fn, res.Err())
	}

	wallet := &storage.Wallet{}

	if err := res.Scan(&wallet.Id, &wallet.Balance); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return wallet, nil
}

func (d *PostgresDB) Transfer(fromId, toId string, amount float64) error {
	const fn = "postgres.Transfer"

	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.Exec(`UPDATE wallets SET balance = balance - $1 WHERE id = $2;`, amount, fromId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.Exec(`UPDATE wallets SET balance = balance + $1 WHERE id = $2;`, amount, toId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.Exec(`INSERT INTO history (fromId, toId, amount) VALUES ($1, $2, $3);`, fromId, toId, amount)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", fn, err)
	}

	return tx.Commit()
}

func (d *PostgresDB) History(id string) ([]storage.Transfer, error) {
	const fn = "postgres.History"

	rows, err := d.db.Query(`SELECT time, fromId, toId, amount
				FROM history WHERE fromId = $1 OR toId = $2;`, id, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	transfers := make([]storage.Transfer, 0)

	for rows.Next() {
		newTransfer := storage.Transfer{}
		err := rows.Scan(&newTransfer.Time, &newTransfer.FromId, &newTransfer.ToId, &newTransfer.Amount)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		transfers = append(transfers, newTransfer)
	}

	return transfers, nil
}

func (d *PostgresDB) IsWalletExists(id string) (bool, error) {
	const fn = "postgres.isWalletExists"

	if err := uuid.Validate(id); err != nil {
		return false, nil
	}

	exists := false

	res := d.db.QueryRow(`SELECT EXISTS(SELECT id FROM wallets WHERE id = $1);`, id)
	if res.Err() != nil {
		return false, fmt.Errorf("%s: %w", fn, res.Err())
	}

	if err := res.Scan(&exists); err != nil {
		return false, fmt.Errorf("%s: %w", fn, res.Err())
	}

	return exists, nil
}
