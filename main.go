package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/shomali11/xsql"

	// postgres driver
	_ "github.com/lib/pq"
)

func getenvInt(key string, def int) int {
	s := os.Getenv(key)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func main() {
	query := os.Getenv("STAT_QUERY")           // they sql query to execute
	postURL := os.Getenv("STAT_URL")           // URL to HTTP POST query output
	interval := getenvInt("STAT_INTERVAL", 10) // interval in minutes, defaults to 10
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PW"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL"),
	)

	// connect to postgres database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// initialize HTTP client
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	// create new ticker on configured interval
	ticker := time.NewTicker(time.Duration(interval) * time.Minute)

	// execute forever
	for {
		// wait for ticker
		<-ticker.C
		// query db and send result every interval
		rows, err := db.Query(query)
		if err != nil {
			log.Println(err)
			continue
		}
		results, err := xsql.Pretty(rows)
		if err != nil {
			log.Println(err)
			continue
		}

		// marshal body
		body, err := json.Marshal(map[string]interface{}{
			"text": fmt.Sprintf("```%s```", results),
		})

		if err != nil {
			log.Println(err)
			continue
		}

		// create new request
		request, err := http.NewRequest("POST", postURL, bytes.NewBuffer(body))

		if err != nil {
			log.Println(err)
			continue
		}

		// set headers
		request.Header.Set("Content-Type", "application/json")

		// do request
		resp, err := client.Do(request)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println(resp)
	}
}
