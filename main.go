package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type report struct {
	thID     int
	duration time.Duration
	hasError bool
}

func main() {
	// log.Printf("Snowflake Stress Tool")
	mode := flag.String("mode", "stress", "\"stress\" or \"bulk\"")
	bulkFile := flag.String("bulkFile", "", "bulk SQL file to use")
	pathSQL := flag.String("pathSQL", "sql", "path where the SQL files to run are located")
	backend := flag.String("backend", "snowflake", "database backend: \"snowflake\", \"sqlserver\" or \"odbc\"")
	duration := flag.Int("duration", 300, "duration in seconds of stress")
	concurrent := flag.Int("concurrent", 50, "number of concurrent queries")
	// log.Printf("loading SQL queries from: \"%s\"", *pathSQL)
	flag.Parse()
	// Check if folder exists
	if _, err := os.Stat(*pathSQL); os.IsNotExist(err) {
		log.Fatalf("SQL queries folder \"%s\" not exists or not permission to access it", *pathSQL)
	}
	if *duration < *concurrent {
		log.Fatalf("Duration can't be less than concurrent: %d < %d", *duration, *concurrent)
	}

	backendSel := sqlEngineSnowflake
	if *backend == "sqlserver" {
		backendSel = sqlEngineSQLServer
	} else if *backend == "odbc" {
		backendSel = sqlEngineODBC
	}
	if *mode == "bulk" {

		bulk(backendSel, *bulkFile)
		return
	}

	quit := make(chan interface{})
	wg := new(sync.WaitGroup)
	wg.Add(1)
	goWorking := true
	resultsCh := make(chan report)
	go func() {
		// log.Printf("** Timeout")
		<-time.After(time.Duration(*duration) * time.Second)
		goWorking = false
		log.Printf("** Timeout done")
		close(quit)
		close(resultsCh)
		wg.Done()
	}()
	sqls, err := getFilesInFolder(*pathSQL)
	if err != nil {
		log.Fatalf("Can't get files in folder: %s", err)
	}
	if len(sqls) == 0 {
		log.Fatalf("No tests in SQL path!")
	}
	log.Printf("Started for %d seconds", *duration)
	reports := make([]report, 0)
	go func() {
		for goWorking {
			packet := <-resultsCh
			reports = append(reports, packet)
		}
	}()

	for i := 0; i < *concurrent; i++ {
		go queryThread(i, quit, *pathSQL, sqls, resultsCh, backendSel)
		time.Sleep(1 * time.Second)
	}
	// Wait
	wg.Wait()
	// Print results
	processResults(reports)
	log.Printf("Bye!")

}

func queryThread(thNum int, quit chan interface{}, basePath string, sqls []string, resultCh chan report, backend sqlEngine) {
	// log.Printf("[%04d] Thread start", thNum)
	for {
		select {
		case <-quit:
			log.Printf("[%04d] quitting", thNum)
			return
		default:
			// Get random SQL File
			sql, sqlFile, err := getRandomSQLquery(basePath, sqls)
			if err != nil {
				log.Printf("[%04d] can't load file. try again", thNum)
			}
			log.Printf("[%04d] Thread running SQL: %s", thNum, sqlFile)
			duration, err := executeSQL(sql, backend)
			if err != nil {
				log.Printf("[%04d] can't execute SQL: %s", thNum, err)
				resultCh <- report{thID: thNum, duration: duration, hasError: true}
			} else {
				log.Printf("[%04d] Thread end running SQL", thNum)
				resultCh <- report{thID: thNum, duration: duration, hasError: false}
			}
		}
	}
}

func processResults(reports []report) {
	// Process Results
	log.Printf("Stats")
	ok := 0
	fail := 0
	var totalTime time.Duration
	totalRuns := len(reports)
	if totalRuns == 0 {
		log.Printf("No executions!")
		return
	}
	for _, item := range reports {
		if item.hasError {
			fail++
		} else {
			ok++
		}
		totalTime += item.duration
	}

	log.Printf("*** OK: %d", ok)
	log.Printf("*** Fail: %d", fail)
	log.Printf("*** Total test: %d", totalRuns)
	avg := time.Duration(float32(totalTime) / float32(totalRuns))
	log.Printf("*** Avg time per query: %s", avg)
}

func bulk(backend sqlEngine, pathSQL string) error {
	sql, err := getSQLContent(pathSQL)
	if err != nil {
		log.Fatalf("Can't get SQL content")
	}
	durationQueries, durationTotal, rows, err := executeSQLBulk(backend, sql)
	if err != nil {
		log.Fatalf("Can't execute SQL bulk: %s", err)
	}
	fmt.Printf("%.2f,%.2f,%.2f,%d\n", durationTotal.Seconds()-durationQueries.Seconds(), durationQueries.Seconds(), durationTotal.Seconds(), rows)
	return nil

}
