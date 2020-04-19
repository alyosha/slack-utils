package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/slack-go/slack"
)

const mockURL = "files-pri/T012345AB-F01234ABC/download/fake.csv"

func TestDownloadAndReadCSV(t *testing.T) {
	testCases := []struct {
		description     string
		csvDownloadResp []byte
		wantRows        [][]string
		wantErr         error
	}{
		{
			description:     "successful download returns valid rows",
			csvDownloadResp: []byte(mockCSVDownloadResp),
			wantRows: [][]string{
				[]string{"email"},
				[]string{"hoge@email.com"},
				[]string{"foo@email.com"},
				[]string{"bar@email.com"},
			},
			wantErr: nil,
		},
		{
			description: "invalid/empty CSV file err",
			wantErr:     ErrInvalidCSV,
		},
	}

	for _, tc := range testCases {
		mux := http.NewServeMux()
		mux.HandleFunc(fmt.Sprintf("/%v", mockURL), func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(tc.csvDownloadResp)
		})

		testServ := httptest.NewServer(mux)
		defer testServ.Close()

		client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))

		rows, err := DownloadAndReadCSV(client, fmt.Sprintf("%v/%v", testServ.URL, mockURL))

		if diff := pretty.Compare(tc.wantRows, rows); diff != "" {
			t.Fatalf("expected to receive rows: %v, got: %v", tc.wantRows, rows)
			return
		}

		if diff := pretty.Compare(tc.wantErr, err); diff != "" {
			t.Fatalf("wantErr does not match received err: %v, got: %v", tc.wantErr, err)
		}
	}
}
