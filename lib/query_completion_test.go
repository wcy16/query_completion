package lib

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
)

func BenchmarkQueryCompletion(b *testing.B) {
	// build index
	queries := make([]string, 0, 10)
	scores := make([]int, 0, 10)

	// read data
	file, err := os.Open("../data/queries.csv") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		score, err := strconv.Atoi(record[1])
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		queries = append(queries, record[0])
		scores = append(scores, score)
	}

	_ = file.Close()
	log.Println("Read done.")

	log.Println("Total queries", len(queries))

	// build index
	var qc QueryCompletion
	qc.BuildIndex(queries, scores, 10)
	log.Println("Query completion index built.")

	// build queries
	testQueries := make([]string, 0, len(queries))
	for _, q := range queries {
		if len(q) < 4 {
			testQueries = append(testQueries, q)
		} else {
			testQueries = append(testQueries, q[0:4])
		}
	}

	// begin
	log.Println("Begin.")
	var results []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, q := range testQueries {
			results = results[:0]
			_ = qc.Search(&q, &results)
		}
	}
}
