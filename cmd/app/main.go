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

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "Deposits"

	_, depositRows, depositIds := reloadDeposits()
	l.Rows = depositRows

	l.TextStyle = ui.NewStyle(ui.ColorCyan)
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	l.WrapText = false
	l.SetRect(0, 0, 100, 22)

	p := widgets.NewParagraph()
	p.SetRect(0, 22, 100, 26)

	desc := widgets.NewParagraph()
	desc.SetRect(0, 26, 100, 45)

	ui.Render(l, p)

	uiEvents := ui.PollEvents()
	for {
		p.Text = db.LinkToDeposit(depositIds[l.SelectedRow])
		desc.Text = db.GetDepositDescription(depositIds[l.SelectedRow])
		ui.Render(p, desc)

		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-r>", "<Enter>":
			c := exec.Command("xdg-open", p.Text)
			c.Start()
		case "<Delete>":
			db.DisableDeposit(depositIds[l.SelectedRow])
			_, l.Rows, depositIds = reloadDeposits()
			ui.Render(l)
		}

		ui.Render(l)
	}
}

func reloadDeposits() ([]datastruct.DepositRowShort, []string, []int) {
	n := 20
	deposits := db.TopN(n)
	depositIds := make([]int, 0, n)
	depositRows := make([]string, 0, n)

	for _, d := range deposits {
		depositStr := fmt.Sprintf("[%s] %s (%.2f%%)", d.BankName, d.Name, d.Rate)
		depositRows = append(depositRows, depositStr)
		depositIds = append(depositIds, d.ID)
	}

	return deposits, depositRows, depositIds
}
