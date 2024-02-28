package db

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const MAX_PLACEHOLDERS = 60000

func LookupIdsBySerials(connstring string, serials []string) (map[string]int, error) {
	db, err := sql.Open("mysql", connstring)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	placeholderCount := len(serials)
	chunks := placeholderCount / MAX_PLACEHOLDERS // avoiding mysql "too many prepared params error"
	if placeholderCount%MAX_PLACEHOLDERS > 0 {
		chunks++
	}
	result := make(map[string]int)
	for i := 0; i < chunks; i++ {
		start := i * MAX_PLACEHOLDERS
		end := int(math.Min(float64(start+MAX_PLACEHOLDERS), float64(placeholderCount)))
		serialsChunk := serials[start:end]
		serialsToIds, err := performChunkedQuery(db, serialsChunk)
		if err != nil {
			log.Println("Error performing SQL query for chunk: ", i)
			continue
		}
		for k, v := range serialsToIds {
			result[k] = v
		}
	}
	return result, nil
}

func performChunkedQuery(db *sql.DB, serials []string) (map[string]int, error) {
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
						WHERE cert_revoked=0 AND LOWER(certificate_serial) IN (%s)`,
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
