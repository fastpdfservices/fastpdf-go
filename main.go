package main

import (
	"fastpdf/fastpdf"
	"fmt"
)

func main() {
	client := fastpdf.NewPDFClient("f447756e0bc543e1aac304d0fcc2a800", "https://data.fastpdfservice.com", "v1")
	isValid, err := client.ValidateToken()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Token is valid:", isValid)
}
