package main

import (
	"log"
	"os"
	"strings"

	"github.com/s-sokolko/revoked_certs_loader/internal/crl"
	"github.com/s-sokolko/revoked_certs_loader/internal/db"
	"github.com/s-sokolko/revoked_certs_loader/internal/vtigerapi"

	"github.com/joho/godotenv"
)

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file!")
	}
}

func retreiveAndMarkRevokedCerts() {
	urls := strings.Split(os.Getenv("CRL_DISTRIBUTION_POINTS"), " ")
	revokedWithDates := crl.LoadCRLItems(urls)
	serials := make([]string, len(revokedWithDates))
	i := 0
	for serial := range revokedWithDates {
		serials[i] = serial
		i++
	}
	connstring := os.Getenv("CONNSTRING")
	serialsToIds, err := db.LookupIdsBySerials(connstring, serials)
	if err != nil {
		log.Fatalln("Error looking up in the DB: ", err)
	}
	apiurl := os.Getenv("API_URL")
	user := os.Getenv("API_USER")
	key := os.Getenv("API_KEY")
	vtigerapi.UpdateCertificatesViaAPI(apiurl, user, key, serialsToIds, revokedWithDates)
}

func main() {
	loadConfig()
	retreiveAndMarkRevokedCerts()
	log.Println("Success!")
}
