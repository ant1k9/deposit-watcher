package query

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	productDescriptionRe = regexp.MustCompile(`"productDescription":"([^"]+)"`)
)

// GetDepositDescription fetch deposit page and parse the description from it
func GetDepositDescription(link string) string {
	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	res, err := client.Do(req)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	el := doc.Find("#__INITIAL_STATE__").First()
	match := productDescriptionRe.FindStringSubmatch(el.Text())
	if len(match) == 2 {
		return strings.ReplaceAll(match[1], "\\n", "\n")
	}

	return ""
}
