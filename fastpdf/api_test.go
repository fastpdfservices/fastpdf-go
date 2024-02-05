/*
api_test.go

MIT License

FastPDF Service/Fast Track Technologies
*/

package fastpdf

import (
	"os"
	"testing"
)

// Reuse the same PDFClient for all tests
var client *PDFClient

// TestNewPDFClient tests the NewPDFClient function
func TestNewPDFClient(t *testing.T) {
	apiKey := "api-key"
	baseURL := "http://127.0.0.1:5011"
	//baseURL := "https://data.fastpdfservice.com"
	apiVersion := "v1"

	client = NewPDFClient(apiKey, SetBaseURL(baseURL), SetAPIVersion(apiVersion))

	if client.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, got %s", apiKey, client.APIKey)
	}
	if client.BaseURL != baseURL+"/"+apiVersion {
		t.Errorf("Expected BaseURL to be %s, got %s", baseURL+"/"+apiVersion, client.BaseURL)
	}
	if client.APIVersion != apiVersion {
		t.Errorf("Expected APIVersion to be %s, got %s", apiVersion, client.APIVersion)
	}
}

// TestValidateToken tests the ValidateToken method of PDFClient
func TestValidateToken(t *testing.T) {
	isValid, err := client.ValidateToken()
	if err != nil {
		t.Errorf("ValidateToken returned an unexpected error: %v", err)
	}
	if !isValid {
		t.Error("Expected ValidateToken to return true, got false")
	}
}

func TestSplit(t *testing.T) {
	// Perform the split operation
	splits := []int{3, 7}
	response, err := client.Split("../input/sample-multipage.pdf", splits)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	// Check that response is not empty
	if len(response) == 0 {
		t.Fatal("Expected non-empty response, got empty response")
	}

	// Save the result
	_, err = client.Save(response, "../output/split.pdf")
	if err != nil {
		t.Fatalf("Failed to save split result: %v", err)
	}

	// Optionally, check if the file was actually created (not typically necessary)
	if _, err := os.Stat("../output/split.pdf"); os.IsNotExist(err) {
		t.Fatal("Failed to create output file 'output/split.pdf'")
	}
}
func TestSplitZip(t *testing.T) {
	// Perform the SplitZip operation
	splits := [][]int{{0, 1}, {1, 2}}
	zipFileName := "../output/split_zip.zip"
	zipFile, err := client.SplitZip("../input/sample-multipage.pdf", splits)
	if err != nil {
		t.Fatalf("SplitZip failed: %v", err)
	}

	// Save the zip file
	_, err = client.Save(zipFile, zipFileName)
	if err != nil {
		t.Fatalf("Failed to save SplitZip result: %v", err)
	}

	// Optionally, check if the zip file was actually created
	if _, err := os.Stat(zipFileName); os.IsNotExist(err) {
		t.Fatal("Failed to create output file 'split_zip.zip'")
	}
}
