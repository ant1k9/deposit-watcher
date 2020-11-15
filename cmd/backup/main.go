package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	errutils "github.com/ant1k9/deposit-watcher/internal/errors"
)

const (
	dropboxFileUploadURL = "https://content.dropboxapi.com/2/files/upload"
	dropboxDirectory     = "deposits"
	defaultTimeout       = 5
	dbName               = "deposits.db"
)

func main() {
	token := os.Getenv("DB_BACKUP_TOKEN")
	if token == "" {
		log.Fatal("not DB_BACKUP_TOKEN env variable")
	}

	client := http.Client{Timeout: defaultTimeout * time.Second}
	req, err := http.NewRequest(
		http.MethodPost,
		dropboxFileUploadURL,
		bytes.NewBuffer(zippedBackup()),
	)
	errutils.FailOnErr(err)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", getPayload())

	response, err := client.Do(req)
	errutils.FailOnErr(err)

	if response.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(response.Body)
		log.Fatal(response.StatusCode, string(data))
	}
}

func getPayload() string {
	path := fmt.Sprintf(
		"/%s/%s.db.zip",
		dropboxDirectory,
		time.Now().Format("20060102_150405"),
	)

	payload, _ := json.Marshal(struct {
		Path           string `json:"path"`
		Mode           string `json:"mode"`
		Autorename     bool   `json:"autorename"`
		Mute           bool   `json:"mute"`
		StrictConflict bool   `json:"strict_conflict"`
	}{
		Path:           path,
		Mode:           "add",
		Autorename:     true,
		Mute:           false,
		StrictConflict: false,
	})
	return string(payload)
}

func zippedBackup() []byte {
	file, err := ioutil.ReadFile(dbName)
	errutils.FailOnErr(err)

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	f, err := w.Create(dbName)
	errutils.FailOnErr(err)

	_, err = f.Write([]byte(file))
	errutils.FailOnErr(err)

	errutils.FailOnErr(w.Close())

	return buf.Bytes()
}
