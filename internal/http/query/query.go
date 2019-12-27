package query

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	errutils "github.com/ant1k9/deposit-watcher/internal/errors"
)

var (
	productDescriptionRe = regexp.MustCompile(`"productDescription":"([^"]+)"`)
)

// GetDepositDescription fetch deposit page and parse the description from it
func GetDepositDescription(link string) string {
	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodGet, link, nil)
	errutils.FailOnErr(err)

	res, err := client.Do(req)
	errutils.FailOnErr(err)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	errutils.FailOnErr(err)

	el := doc.Find("#__INITIAL_STATE__").First()
	match := productDescriptionRe.FindStringSubmatch(el.Text())
	if len(match) == 2 {
		return strings.ReplaceAll(match[1], "\\n", "\n")
	}

	return ""
}
