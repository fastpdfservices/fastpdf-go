/*
api.go

MIT License

FastPDF Service/Fast Track Technologies
*/

// Package fastpdf provides a client for the FastPDF Service.
// It allows for interactions with the FastPDF API for tasks such as
// validating tokens, converting, and processing PDFs.
package fastpdf

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// PDFClient struct represents a client for PDF services. It provides methods
// for interacting with the FastPDF service, such as validating API tokens and
// processing PDFs. The client maintains configuration such as the API version,
// base URL, API key, custom headers, and supported image formats.
type PDFClient struct {
	APIVersion            string
	BaseURL               string
	APIKey                string
	Headers               http.Header
	SupportedImageFormats []string
}

// NewPDFClient creates a new instance of PDFClient with default values.
// It accepts an API key and a variadic list of options to customize the client.
// By default, the client is configured with the API version 'v1' and the
// standard base URL of the FastPDF service. Optional configuration functions
// like SetBaseURL and SetAPIVersion can be used to modify these defaults.
//
// Example usage:
//
//	client := NewPDFClient("your-api-key")
func NewPDFClient(apiKey string, options ...func(*PDFClient)) *PDFClient {
	client := &PDFClient{
		APIVersion: "v1",
		BaseURL:    "https://data.fastpdfservice.com",
		APIKey:     apiKey,
		Headers:    http.Header{"Authorization": []string{apiKey}},
		SupportedImageFormats: []string{
			"jpeg", "png", "gif", "bmp", "tiff", "webp", "svg", "ico", "pdf",
			"psd", "ai", "eps", "cr2", "nef", "sr2", "orf", "rw2", "dng",
			"arw", "heic",
		},
	}

	for _, option := range options {
		option(client)
	}

	client.BaseURL = fmt.Sprintf("%s/%s", client.BaseURL, client.APIVersion)

	return client
}

// SetBaseURL is an optional function used to configure a custom base URL
// for the PDFClient. This function returns a function that accepts a *PDFClient
// and modifies its BaseURL field.
//
// Example usage:
//
//	client := NewPDFClient("your-api-key", SetBaseURL("other-url"))
func SetBaseURL(url string) func(*PDFClient) {
	return func(c *PDFClient) {
		c.BaseURL = url
	}
}

// SetAPIVersion is an optional function used to set a specific API version
// for the PDFClient. This function returns a function that accepts a *PDFClient
// and modifies its APIVersion field.
//
// Example usage:
//
//	client := NewPDFClient("your-api-key", SetAPIVersion("v2"))
func SetAPIVersion(version string) func(*PDFClient) {
	return func(c *PDFClient) {
		c.APIVersion = version
	}
}

func readFile(file interface{}) (filename string, fileContent []byte, contentType string, err error) {
	switch f := file.(type) {
	case string:
		filename = filepath.Base(f)
		fileContent, err = os.ReadFile(f)
		if err != nil {
			return
		}
		contentType = mime.TypeByExtension(filepath.Ext(f))

	case []byte:
		fileContent = f
		contentType = http.DetectContentType(f)

	default:
		err = fmt.Errorf("unsupported file type: %T", file)
		return
	}

	if contentType == "" {
		// Default to "application/octet-stream" if content type is not detected
		contentType = "application/octet-stream"
	}

	return
}

// ValidateToken checks if the API token is valid. It sends a request to the
// FastPDF service and returns true if the token is valid. If the service
// returns an error or an invalid status code, it returns false and an error.
//
// Example usage:
//
//	client := NewPDFClient("your-api-key")
//	isValid, err := client.ValidateToken()
//	if err != nil {
//	    // Handle error
//	}
//	if isValid {
//	    // Proceed with using the client
//	}
func (client *PDFClient) ValidateToken() (bool, error) {
	url := client.BaseURL + "/token"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header = client.Headers

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return false, NewPDFError(response)
	}
	return true, nil
}

func (c *PDFClient) makeMultipartPostRequest(url string, fileField string, fileData []byte, filename string, extraFields map[string]string) ([]byte, error) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	// Create file part
	part, err := writer.CreateFormFile(fileField, filename)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(fileData)
	if err != nil {
		return nil, err
	}

	// Add additional fields
	for key, value := range extraFields {
		err := writer.WriteField(key, value)
		if err != nil {
			return nil, err
		}
	}

	// Close the writer before making the request
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	for key, values := range c.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewPDFError(resp)
	}

	return io.ReadAll(resp.Body)
}

// Save writes the given content to a file at the specified file path.
// If the file path is not provided, it returns the content as a byte slice.
//
// Example usage:
//
//	client := NewPDFClient("your-api-key")
//	pdfContent, err := client.URLToPDF("https://www.example.com")
//	if err != nil {
//	    // Handle error
//	}
//
//	// Save to file
//	err = client.Save(pdfContent, "path/to/save/your.pdf")
//	if err != nil {
//	    // Handle error
//	}
//
//	// Get as byte slice
//	contentBytes, err := client.Save(pdfContent, "")
//	if err != nil {
//	    // Handle error
//	}
func (c *PDFClient) Save(content []byte, filePath string) ([]byte, error) {
	if filePath != "" {
		// Write content to file
		err := os.WriteFile(filePath, content, 0644)
		if err != nil {
			return nil, err
		}
		return nil, nil
	} else {
		// Return content as a byte slice
		return content, nil
	}
}

// Split divides a PDF file at the specified page numbers.
//
// This method takes the path to the PDF file or the file itself as a byte slice,
// along with a slice of page numbers indicating where to split the file. It returns
// the content of the split PDF files as a byte slice.
//
// If the request to the FastPDF service fails, it returns a PDFException error.
//
// Example usage:
//
//	client := NewPDFClient("your-api-key")
//	splitPDFContent, err := client.Split("path/to/your.pdf", []int{3, 6})
//	if err != nil {
//	    // Handle error
//	}
//	// Use splitPDFContent
func (c *PDFClient) Split(file interface{}, splits []int) ([]byte, error) {
	filename, fileContent, _, err := readFile(file)
	if err != nil {
		return nil, err
	}

	splitsData, err := json.Marshal(splits)
	if err != nil {
		return nil, err
	}

	return c.makeMultipartPostRequest(
		c.BaseURL+"/pdf/split",
		"file",
		fileContent,
		filename,
		map[string]string{"splits": string(splitsData)},
	)
}

// SplitZip splits a PDF file into multiple parts based on the provided splits.
//
// This method takes the path to the PDF file or the file itself as a byte slice,
// along with a 2D slice representing the page ranges to split the PDF.
// Each inner slice contains two integers representing the start and end page of a split
//
// If the request to the FastPDF service fails, it returns a PDFException error.
//
// Example usage:
//
// client := NewPDFClient("your-api-key")
// zipContent, err := client.SplitZip("path/to/your.pdf", [][]int{{0, 1}, {1, 2}})
//
//	if err != nil {
//	    // Handle error
//	}
//
// outputZipPath := "output.zip"
// err = os.WriteFile(outputZipPath, zipContent, 0644)
//
//	if err != nil {
//	    // Handle error
//	}
//
// fmt.Println("ZIP file created successfully:", outputZipPath)
func (c *PDFClient) SplitZip(file interface{}, splits [][]int) ([]byte, error) {
	filename, fileContent, _, err := readFile(file)
	if err != nil {
		return nil, err
	}

	splitsData, err := json.Marshal(splits)
	if err != nil {
		return nil, err
	}

	return c.makeMultipartPostRequest(
		c.BaseURL+"/pdf/split-zip",
		"file",
		fileContent,
		filename,
		map[string]string{"splits": string(splitsData)},
	)
}

// Extract extracts files from a zip archive.
// It takes a byte slice representing the zip archive and an output path where the extracted files will be saved.
//
// If the request to the FastPDF service fails, it returns a PDFException error.
//
// Example Usage:
//
// client := fastpdf.NewPDFClient("your-api-key")
// zipBytes, err := ioutil.ReadFile(zipFilePath)
//
//	if err != nil {
//	     // Handle error
//	}
//
// fmt.Println("Files extracted successfully to:", outputPath)
// }
// Appeler la fonction Extract
// err = client.Extract(zipBytes, outputPath)
//
//	if err != nil {
//		    // Handle error
//	}
func (c *PDFClient) Extract(zipBytes []byte, outputPath string) error {
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		fPath := filepath.Join(outputPath, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
		} else {
			f, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			rc, err := file.Open()
			if err != nil {
				f.Close()
				return err
			}
			_, err = io.Copy(f, rc)
			f.Close()
			rc.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// EditMetadata edits the metadata of a PDF file.
// It returns the content of the modified PDF file as a byte slice.
//
// If the request to the FastPDF service fails, it returns a PDFException error.
//
// Example Usage:
//
// client := fastpdf.NewPDFClient("f447756e0bc543e1aac304d0fcc2a800")
//
//	metadata := map[string]string{
//			"/Title":    "John Doe",
//			"/Author":   "Updated Title",
//			"/Subject":  "Soujet",
//			"/Producer": "John Doe",
//			"/Creator":  "IbraTV",
//	}
//
// pdfBytes, err := os.ReadFile(pdfFilePath)
//
//	if err != nil {
//		    // Handle error
//		}
//		updatedPdfBytes, err := client.EditMetadata(pdfBytes, metadata)
//		if err != nil {
//		    // Handle error
//		}
func (c *PDFClient) EditMetadata(file interface{}, metadata map[string]string) ([]byte, error) {
	filename, fileContent, _, err := readFile(file)
	if err != nil {
		return nil, err
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return c.makeMultipartPostRequest(
		c.BaseURL+"/pdf/metadata",
		"file",
		fileContent,
		filename,
		map[string]string{"metadata": string(metadataJSON)},
	)
}
