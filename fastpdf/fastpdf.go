package fastpdf

import (
	"errors"
	"net/http"
)

type PDFClient struct {
	APIVersion string
	BaseURL    string
	APIKey     string
	Headers    http.Header
	// autres champs si nécessaires
}

func NewPDFClient(apiKey string, baseURL string, apiVersion string) *PDFClient {
	return &PDFClient{
		APIVersion: apiVersion,
		BaseURL:    baseURL + "/" + apiVersion,
		APIKey:     apiKey,
		Headers:    http.Header{"Authorization": []string{apiKey}},
	}
}

func (client *PDFClient) ValidateToken() (bool, error) {
	url := client.BaseURL + "/token"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header = client.Headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("failed to validate token")
	}
	return true, nil
}
