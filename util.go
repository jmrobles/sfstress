package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/snowflakedb/gosnowflake"
)

type sqlEngine int

const (
	sqlEngineSnowflake sqlEngine = iota
	sqlEngineSQLServer
	sqlEngineODBC
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
	schema := getEnvDefault("SF_SCHEMA", "")
	if schema == "" {
		return fmt.Sprintf("%s:%s@%s/%s?warehouse=%s", user, password, host, database, warehouse)
	}
	return fmt.Sprintf("%s:%s@%s/%s/%s?warehouse=%s", user, password, host, database, schema, warehouse)

}

func getConnStringSQLServer() string {
	server := getEnvDefault("MSSQL_SERVER", "localhost")
	port := getEnvDefault("MSSQL_PORT", "1433")
	user := getEnvDefault("MSSQL_USER", "SA")
	password := getEnvDefault("MSSQL_PASSWORD", "")
	database := getEnvDefault("MSSQL_DB", "")
	return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", server, user, password, port, database)
}

func getConnStringODBC() string {
	dsn := getEnvDefault("ODBC_DSN", "")
	password := getEnvDefault("ODBC_PASSWORD", "")
	return fmt.Sprintf("DSN=%s;PWD=%s", dsn, password)
	//return fmt.Sprintf("DRIVER={SnowflakeDSIIDriver};server=axesor.west-europe.azure.snowflakecomputing.com;database=STAGE_CO_DB;warehouse=TEST_WH;UID=CIVICA_POC;PWD=CIVI19ca_;AUTOCOMMIT=FALSE")
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

func executeSQL(sqlQuery string, mode sqlEngine) (time.Duration, error) {
	// Get new connection
	var connStr = ""
	var engine = ""
	if mode == sqlEngineSQLServer {
		connStr = getConnStringSQLServer()
		engine = "mssql"
	} else if mode == sqlEngineODBC {
		connStr = getConnStringODBC()
		engine = "odbc"
	} else {
		connStr = getConnStringSF()
		engine = "snowflake"
		// log.Printf("DSN : %s", connStr)
	}
	db, err := sql.Open(engine, connStr)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	// Change session to not use cached result
	if mode == sqlEngineSnowflake || mode == sqlEngineODBC {
		_, err = db.Exec("ALTER SESSION SET USE_CACHED_RESULT = False")
		if err != nil {
			return 0, err
		}
	} else if mode == sqlEngineSQLServer {
		_, err = db.Exec("DBCC DROPCLEANBUFFERS")
		if err != nil {
			return 0, err
		}
		_, err = db.Exec("DBCC FREEPROCCACHE")
		if err != nil {
			return 0, err
		}
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
