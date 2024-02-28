package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func LookupIdsBySerials(connstring string, serials []string) (map[string]int, error) {
	db, err := sql.Open("mysql", connstring)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	questionMarks := make([]string, len(serials))
	params := make([]any, len(serials))
	for i := range questionMarks {
		questionMarks[i] = "?"
		params[i] = serials[i]
	}
	questionMarksJoined := strings.Join(questionMarks, ", ")
	sql := fmt.Sprintf(`SELECT 
						certificatesid AS id, LOWER(certificate_serial) AS serial 
						FROM vtiger_certificates 
						WHERE COALESCE(cert_revoked, 0) = 0 AND LOWER(certificate_serial) IN (%s)`,
		questionMarksJoined)

	results, err := db.Query(sql, params...)
	if err != nil {
		return nil, err
	}

	serialToIds := make(map[string]int)
	for results.Next() {
		var id int
		var serial string
		err := results.Scan(&id, &serial)
		if err != nil {
			return nil, err
		}
		serialToIds[serial] = id
	}
	err = results.Close()
	if err != nil {
		log.Println("Error closing dataset ", err)
	}
	return serialToIds, nil
}
