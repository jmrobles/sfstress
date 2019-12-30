package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"

	_ "github.com/snowflakedb/gosnowflake"
)

func getEnvDefault(name string, def string) string {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	return v
}

func getConnStringSF() string {
	host := getEnvDefault("SF_HOST", "")
	user := getEnvDefault("SF_USER", "")
	password := getEnvDefault("SF_PASSWORD", "")
	database := getEnvDefault("SF_DB", "")
	warehouse := getEnvDefault("SF_WAREHOUSE", "")
	return fmt.Sprintf("%s:%s@%s/%s?warehouse=%s", user, password, host, database, warehouse)
}

func getFilesInFolder(folder string) ([]string, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ret = append(ret, f.Name())
	}
	return ret, nil
}

func getRandomSQLquery(basePath string, sqlFiles []string) (string, string, error) {

	rndFile := sqlFiles[rand.Int31n(int32(len(sqlFiles)))]
	sql, err := ioutil.ReadFile(path.Join(basePath, rndFile))
	if err != nil {
		return "", "", err
	}
	return string(sql), rndFile, nil
}

func executeSQL(sqlQuery string) (time.Duration, error) {
	// Get new connection
	db, err := sql.Open("snowflake", getConnStringSF())
	if err != nil {
		return 0, err
	}
	defer db.Close()
	// Change session to not use cached result
	_, err = db.Exec("ALTER SESSION SET USE_CACHED_RESULT = False")
	if err != nil {
		return 0, err
	}
	start := time.Now()
	// Execute SQL
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return time.Since(start), err
	}
	// Ignore rows
	rows.Close()
	// Close connection
	return time.Since(start), nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
