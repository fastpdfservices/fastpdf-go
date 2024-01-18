/*
errors.go

MIT License

FastPDF Service/Fast Track Technologies
*/

package fastpdf

import (
	"net/http"
	"fmt"
	"io"
)

type PDFError struct {
    StatusCode int
    Status     string
    Response   string
    Message    string
}

func (e *PDFError) Error() string {
    return fmt.Sprintf("%s. Status Code: %d, Response: %s", e.Message, e.StatusCode, e.Response)
}

func NewPDFError(resp *http.Response) error {
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        bodyBytes = []byte("failed to read response body")
    }
    return &PDFError{
        StatusCode: resp.StatusCode,
        Response:   string(bodyBytes),
        Message:    "Server returned an HTTP error",
    }
}


