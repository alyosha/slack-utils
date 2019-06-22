package utils

import (
	"bytes"
	"encoding/csv"
	"errors"

	"github.com/nlopes/slack"
)

var ErrInvalidCSV = errors.New("received invalid/empty CSV file")

func DownloadAndReadCSV(userClient *slack.Client, callback *slack.InteractionCallback) ([][]string, error) {
	b := bytes.Buffer{}
	err := userClient.GetFile(callback.Message.Files[0].URLPrivateDownload, &b)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(&b)
	rows, err := r.ReadAll()

	if len(rows) == 0 {
		return nil, ErrInvalidCSV
	}

	return rows, nil
}
