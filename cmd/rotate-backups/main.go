package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	errutils "github.com/ant1k9/deposit-watcher/internal/errors"
)

const (
	dropboxFileListURL   = "https://api.dropboxapi.com/2/files/list_folder"
	dropboxFileDeleteURL = "https://api.dropboxapi.com/2/files/delete"
	dropboxDirectory     = "/deposits"
	defaultTimeout       = 5
	hardLimit            = 30
)

type Entry struct {
	Name string `json:"name"`
}

type EntryInfo struct {
	Entries []Entry `json:"entries"`
}

func main() {
	token := os.Getenv("DB_BACKUP_TOKEN")
	if token == "" {
		log.Fatal("not DB_BACKUP_TOKEN env variable")
	}

	entries := getEntries(token)

	if len(entries.Entries) > hardLimit {
		sort.Slice(entries.Entries, func(i, j int) bool {
			return entries.Entries[i].Name < entries.Entries[j].Name
		})
		for i := 0; i < len(entries.Entries)-hardLimit; i++ {
			deleteEntry(entries.Entries[i].Name, token)
		}
	}
}

func deleteEntry(name, token string) {
	payload := pathPayload(fmt.Sprintf("%s/%s", dropboxDirectory, name))
	doRequest(dropboxFileDeleteURL, token, payload)
}

func doRequest(url, token string, payload []byte) *http.Response {
	client := http.Client{Timeout: defaultTimeout * time.Second}
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(payload),
	)
	errutils.FailOnErr(err)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	errutils.FailOnErr(err)

	return response
}

func getEntries(token string) EntryInfo {
	payload := pathPayload(dropboxDirectory)
	response := doRequest(dropboxFileListURL, token, payload)

	data, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		log.Fatal(response.StatusCode, string(data))
	}

	entries := EntryInfo{}
	errutils.FailOnErr(json.Unmarshal(data, &entries))

	return entries
}

func pathPayload(path string) []byte {
	payload, _ := json.Marshal(struct {
		Path string `json:"path"`
	}{Path: path})
	return payload
}
