package crl

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func fetchCRLData(url string) (*[]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errorText := fmt.Sprintf("Wrong HTTP status code of %d", resp.StatusCode)
		return nil, errors.New(errorText)
	}

	defer resp.Body.Close()

	buffer, err := io.ReadAll(resp.Body)

	return &buffer, err
}

func parseCRL(buffer *[]byte) (map[string]string, error) {
	crl, err := x509.ParseCRL(*buffer)

	if err != nil {
		return nil, err
	}

	items := make(map[string]string)
	for _, cert := range crl.TBSCertList.RevokedCertificates {
		//serialText := fmt.Sprintf("%034x", cert.SerialNumber)
		revokedDateStr := cert.RevocationTime.Format("2006-01-02")
		serialText := cert.SerialNumber.Text(16)
		items[serialText] = revokedDateStr
	}
	return items, err
}

func LoadCRLItems(urls []string) map[string]string {
	result := make(map[string]string)
	for _, url := range urls {
		bufferPtr, err := fetchCRLData(url)
		if err != nil {
			log.Println("Error fetching URL: ", err)
			continue
		}
		items, err := parseCRL(bufferPtr)
		if err != nil {
			log.Println("Error parsing CRL: ", err)
			continue
		}
		//copy to result
		for k, v := range items {
			result[k] = v
		}
	}
	return result
}
