package postgres

import (
	"fmt"
)

func (d *PostgresDB) InitDB() error {
	const fn = "postgres.Init"

	initExtQ := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	if _, err := d.db.Exec(initExtQ); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	initWalletsQ := `CREATE TABLE IF NOT EXISTS wallets(
		id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
		balance numeric(12, 2) DEFAULT 100.00 NOT NULL CHECK (balance >= 0)
	);`

	if _, err := d.db.Exec(initWalletsQ); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	initHistoryQ := `CREATE TABLE IF NOT EXISTS history(
		time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		fromId UUID NOT NULL,
		toId UUID NOT NULL,
		amount numeric(12, 2) NOT NULL
	);`

	if _, err := d.db.Exec(initHistoryQ); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
