package sheets

import (
	"context"
	"os"
	"testing"

	"github.com/username/go-util/google/auth"
	"google.golang.org/api/sheets/v4"
)

func TestCreateSpreadsheet(t *testing.T) {
	ctx := context.Background()

	if _, err := os.Stat("credentials.json"); os.IsNotExist(err) {
		t.Skip("credentials.json not found; skipping test")
	}

	// Authenticate using google/auth
	ga, err := auth.NewGoogleAuth(ctx, "credentials.json", sheets.SpreadsheetsScope)
	if err != nil {
		t.Fatalf("NewGoogleAuth failed: %v", err)
	}

	client, err := ga.NewClient(ctx, "token.json")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Create Sheets service
	srv, err := NewSheetsService(ctx, client)
	if err != nil {
		t.Fatalf("NewSheetsService failed: %v", err)
	}

	spreadsheetID, err := srv.CreateSpreadsheet(ctx, "TestSpreadsheet")
	if err != nil {
		t.Fatalf("CreateSpreadsheet failed: %v", err)
	}
	if spreadsheetID == "" {
		t.Error("Expected non-empty spreadsheet ID")
	}
}