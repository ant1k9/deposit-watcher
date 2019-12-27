package main

import (
	"fmt"
	"log"
	"os/exec"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/ant1k9/deposit-watcher/internal/datastruct"
	"github.com/ant1k9/deposit-watcher/internal/db"
)

const (
	n = 25
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	page := 1
	reverse := true
	_, depositRows, depositIds := reloadDeposits(page, reverse)

	l := prepareDepositList()
	l.Rows = depositRows
	ui.Render(l)

	desc := prepareDepositDescription()

	uiEvents := ui.PollEvents()
	for {
		link := db.LinkToDeposit(depositIds[l.SelectedRow])
		desc.Text = db.GetDepositDescription(depositIds[l.SelectedRow])
		ui.Render(desc)

		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "g", "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		case "r":
			reverse = !reverse
			_, l.Rows, depositIds = reloadDeposits(page, reverse)
		case "<PageDown>":
			page++
			_, l.Rows, depositIds = reloadDeposits(page, reverse)

			if len(l.Rows) == 0 {
				page--
			}
			_, l.Rows, depositIds = reloadDeposits(page, reverse)
			l.ScrollTop()
		case "<PageUp>":
			if page > 1 {
				page--
			}
			_, l.Rows, depositIds = reloadDeposits(page, reverse)
			l.ScrollTop()
		case "<C-r>", "<Enter>":
			c := exec.Command("xdg-open", link)
			c.Start()
		case "<Delete>":
			db.DisableDeposit(depositIds[l.SelectedRow])
			if l.SelectedRow == 0 && page > 0 {
				page--
			}
			_, l.Rows, depositIds = reloadDeposits(page, reverse)
			if l.SelectedRow == len(l.Rows) {
				l.ScrollUp()
			}
		}

		ui.Render(l)
	}
}

func reloadDeposits(page int, reverse bool) ([]datastruct.DepositRowShort, []string, []int) {
	deposits := db.TopN(n, page, reverse)
	depositIds := make([]int, 0, n)
	depositRows := make([]string, 0, n)

	for _, d := range deposits {
		depositStr := fmt.Sprintf("[%s] %s (%.2f%%)", d.BankName, d.Name, d.Rate)
		depositRows = append(depositRows, depositStr)
		depositIds = append(depositIds, d.ID)
	}

	return deposits, depositRows, depositIds
}

func prepareDepositList() *widgets.List {
	l := widgets.NewList()
	l.Title = "Deposits"
	l.TextStyle = ui.NewStyle(ui.ColorCyan)
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	l.WrapText = false
	l.SetRect(0, 0, 100, n+2)
	return l
}

func prepareDepositDescription() *widgets.Paragraph {
	desc := widgets.NewParagraph()
	desc.SetRect(0, n+2, 100, n+20)
	return desc
}
