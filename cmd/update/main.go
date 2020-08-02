package main

import (
	"log"

	"github.com/ant1k9/deposit-watcher/internal/db"
	"github.com/ant1k9/deposit-watcher/internal/http"
)

func main() {
	total := http.GetTotal(1)
	existedIds := make([]int, 0, total)

	for i := 1; i < total/http.PageSize+1; i++ {
		for _, deposit := range http.GetDeposits(i) {
			bank, err := db.GetOrCreateBankForDeposit(deposit)
			if err != nil {
				log.Println(err)
				continue
			}

			id, err := db.CreateOrUpdateDeposit(deposit, bank)
			if err != nil {
				log.Println(err)
			}
			existedIds = append(existedIds, id)
		}
	}

	db.SetActiveDeposits(existedIds)
}
