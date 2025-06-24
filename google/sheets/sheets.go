package sheets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/username/go-util/google/auth"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// SheetsService provides access to Google Sheets API.
type SheetsService struct {
	srv *sheets.Service
}

// NewSheetsService creates a new Sheets service with an authenticated client.
func NewSheetsService(ctx context.Context, client *http.Client) (*SheetsService, error) {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Sheets service: %w", err)
	}
	return &SheetsService{srv: srv}, nil
}

// CreateSpreadsheet creates a new spreadsheet.
func (s *SheetsService) CreateSpreadsheet(ctx context.Context, title string) (string, error) {
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{Title: title},
	}
	created, err := s.srv.Spreadsheets.Create(spreadsheet).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create spreadsheet: %w", err)
	}
	return created.SpreadsheetId, nil
}