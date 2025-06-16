package auth

import (
	"path/filepath"
	"testing"

	"golang.org/x/oauth2"
)

func TestNewGauth(t *testing.T) {
	ga := NewGauth("test_credentials.json", "test_scope")
	if ga.credentialsFile != "test_credentials.json" {
		t.Errorf("expected credentialsFile to be 'test_credentials.json', got %s", ga.credentialsFile)
	}
	if ga.scope != "test_scope" {
		t.Errorf("expected scope to be 'test_scope', got %s", ga.scope)
	}
	if ga.server != nil {
		t.Errorf("expected server to be nil")
	}
}

func TestSaveAndTokenFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "token.json")
	ga := Gauth{
		tokenFile: tokenFile,
	}

	// Create a dummy token
	testToken := &oauth2.Token{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
	}
	token = testToken

	// Save token
	err := ga.saveToken()
	if err != nil {
		t.Fatalf("saveToken failed: %v", err)
	}

	// Reset global token and read from file
	token = nil
	ga.tokenFile = tokenFile
	err = ga.tokenFromFile()
	if err != nil {
		t.Fatalf("tokenFromFile failed: %v", err)
	}
	if token == nil || token.AccessToken != "test-access-token" {
		t.Errorf("tokenFromFile did not load correct token")
	}
}

func TestTokenFromFile_FileNotFound(t *testing.T) {
	ga := Gauth{tokenFile: "nonexistent.json"}
	err := ga.tokenFromFile()
	if err == nil {
		t.Errorf("expected error for missing token file")
	}
}

func TestSaveToken_FileError(t *testing.T) {
	// Use a directory as the file path to force an error
	ga := Gauth{tokenFile: t.TempDir()}
	token = &oauth2.Token{}
	err := ga.saveToken()
	if err == nil {
		t.Errorf("expected error when saving token to a directory")
	}
}
