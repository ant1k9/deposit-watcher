package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //nolint

	"github.com/ant1k9/deposit-watcher/internal/datastruct"
	"github.com/ant1k9/deposit-watcher/internal/http/query"
)

const (
	// DBName is the name of the database with bank and deposits
	DBName = "deposits.db"
)

var (
	db *sqlx.DB
)

func init() {
	conn, err := sqlx.Connect("sqlite3", DBName)
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
		if depositRow.Rate != newRow.Rate && !depositRow.Off {
			now := time.Now().Format("2006-01-02")
			result := db.MustExec(
				`UPDATE deposit SET is_updated = TRUE, rate = ?, has_replenishment = ?,
				detail = ?, minimal_amount = ?, previous_rate = ?, updated_at = ?
				WHERE alias = ? AND bank_id = ?`,
				newRow.Rate, newRow.HasReplenishment, newRow.Detail,
				newRow.MinimalAmount, depositRow.Rate, now,
				newRow.Alias, newRow.BankID,
			)
			logChange("update", depositRow, newRow, bank)

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
	logChange("create", depositRow, newRow, bank)

	_, err := result.LastInsertId()
	return err
}

func logChange(operation string, depositRow, newRow datastruct.DepositRow, bank datastruct.BankRow) {
	switch operation {
	case "update":
		fmt.Printf(
			"\033[1mupdate [%s] %s (%f) -> (%f%%)\033[0m\n",
			bank.Name, depositRow.Name, depositRow.Rate, newRow.Rate,
		)
	case "create":
		fmt.Printf(
			"\033[1mcreate [%s] %s (%f%%)\033[0m\n",
			bank.Name, newRow.Name, newRow.Rate,
		)
	default:
	}
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
		"https://www.sravni.ru/bank/%s/vklad/%s/", bankAlias, depositAlias,
	)
}

// TopN returns top n deposits sorted by rate wigh pagination
func TopN(n, page int, desc bool) []datastruct.DepositRowShort {
	var deposits []datastruct.DepositRowShort

	order := "ASC"
	if desc {
		order = "DESC"
	}

	err := db.Select(
		&deposits,
		`SELECT d.id id, d.alias alias, d.name name, detail, rate, has_replenishment, b.name bank_name
			FROM deposit d JOIN bank b ON d.bank_id = b.id
			WHERE NOT off ORDER BY rate `+order+` LIMIT ? OFFSET ?`, n, (page-1)*n,
	)

	if err != nil {
		log.Fatal(err)
		return make([]datastruct.DepositRowShort, 0)
	}

	return deposits
}

// DisableDeposit is set flag off to deposit record and it won't be updated anymore
// It also won't be shown in application to an user
func DisableDeposit(id int) error {
	result := db.MustExec("UPDATE deposit SET off = TRUE WHERE id = ?", id)
	_, err := result.RowsAffected()
	return err
}

// GetDepositDescription load description for deposit from database.
// If description is empty it makes request to the site and save description in the database.
func GetDepositDescription(id int) string {
	desc, err := getDepositDescription(id)

	if err != nil {
		desc = query.GetDepositDescription(LinkToDeposit(id))
		db.MustExec(
			"INSERT INTO deposit_details (deposit_id, full_description) VALUES (?, ?)",
			id, desc,
		)
		return desc
	}

	return desc
}

func getDepositDescription(id int) (string, error) {
	var desc sql.NullString
	query := db.QueryRow(
		"SELECT full_description FROM deposit_details WHERE deposit_id = ?", id,
	)
	err := query.Scan(&desc)

	return desc.String, err
}

// NewConnection returns a private database connection to sqlite3 database
func NewConnection() *sqlx.DB {
	return db
}
