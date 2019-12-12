package main

import (
	"github.com/ant1k9/deposit-watcher/internal/db"
)

func main() {
	conn := db.NewConnection()
	conn.MustExec(db.InitialMigration)
}
