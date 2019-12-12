package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //nolint

	"github.com/ant1k9/deposit-watcher/internal/datastruct"
)

const (
	// DBName is the name of the database with bank and deposits
	DBName = "deposits.db"
)

var (
	db *sqlx.DB
)

func init() {
	conn, err := sqlx.Connect("sqlite3", "deposits.db")
	if err != nil {
		log.Fatal(err)
	}

	db = conn
}

func convertToDepositRow(deposit datastruct.Deposit, bankID int) datastruct.DepositRow {
	return datastruct.DepositRow{
		Alias:            deposit.Alias,
		Name:             deposit.Name,
		BankID:           bankID,
		MinimalAmount:    deposit.Amount.From,
		Rate:             deposit.Rate,
		HasReplenishment: deposit.Replenishment.Available,
		Detail:           deposit.Replenishment.Description,
	}
}

// CreateOrUpdateDeposit make a new database in database.
// If deposit with given alias is existed then it should be updated.
// In other cases new deposit is created.
func CreateOrUpdateDeposit(deposit datastruct.Deposit, bank datastruct.BankRow) error {
	depositRow := datastruct.DepositRow{}
	newRow := convertToDepositRow(deposit, bank.ID)

	_ = db.Get(
		&depositRow,
		"SELECT * FROM deposit WHERE alias = ? AND bank_id = ?",
		deposit.Alias, bank.ID,
	)

	if depositRow.Alias == deposit.Alias {
		if depositRow.Rate != newRow.Rate {
			result := db.MustExec(
				`UPDATE deposit SET is_updated = TRUE, rate = ?, has_replenishment = ?,
				detail = ?, minimal_amount = ? WHERE alias = ? AND bank_id = ?`,
				newRow.Rate, newRow.HasReplenishment, newRow.Detail,
				newRow.MinimalAmount, newRow.Alias, newRow.BankID,
			)
			_, err := result.RowsAffected()
			return err
		}
		return nil
	}

	result := db.MustExec(`INSERT INTO deposit
		(alias, name, bank_id, minimal_amount, rate, has_replenishment, detail)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		newRow.Alias, newRow.Name, newRow.BankID, newRow.MinimalAmount,
		newRow.Rate, newRow.HasReplenishment, newRow.Detail,
	)
	_, err := result.LastInsertId()

	return err
}

// GetOrCreateBankForDeposit takes bank alias as a key.
// If the bank with given alias is existed then it should be updated.
// If there is no such bank, new bank is created.
func GetOrCreateBankForDeposit(deposit datastruct.Deposit) (datastruct.BankRow, error) {
	bank := datastruct.BankRow{}

	err := db.Get(&bank, "SELECT * FROM bank WHERE alias = ?", deposit.Organization.Alias)
	if err == nil && bank.Alias == deposit.Organization.Alias {
		return bank, nil
	}

	result := db.MustExec(
		"INSERT INTO bank (alias, name) VALUES (?, ?)",
		deposit.Organization.Alias,
		deposit.Organization.Name,
	)

	_, err = result.LastInsertId()
	if err != nil {
		return bank, err
	}

	return GetOrCreateBankForDeposit(deposit)
}

// LinkToDeposit compose link to deposit on sravni.ru site
func LinkToDeposit(id int) string {
	var depositAlias, bankAlias string

	row := db.QueryRow(
		`SELECT d.alias, b.alias FROM deposit d
		JOIN bank b ON d.bank_id = b.id
		WHERE d.id = ?`, id,
	)

	err := row.Scan(&depositAlias, &bankAlias)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(
		"datastructs://www.sravni.ru/bank/%s/vklad/%s/", bankAlias, depositAlias,
	)
}

// NewConnection returns a private database connection to sqlite3 database
func NewConnection() *sqlx.DB {
	return db
}
