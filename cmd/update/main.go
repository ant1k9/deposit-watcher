package main

import (
	"log"

	"github.com/ant1k9/deposit-watcher/internal/db"
	"github.com/ant1k9/deposit-watcher/internal/http"
)

func main() {
	total := http.GetTotal(1)

	for i := 1; i < total/http.PageSize+1; i++ {
		for _, deposit := range http.GetDeposits(i) {
			bank, err := db.GetOrCreateBankForDeposit(deposit)
			if err != nil {
				log.Println(err)
				continue
			}

			err = db.CreateOrUpdateDeposit(deposit, bank)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
