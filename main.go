package main

import (
	"./lib"
	"encoding/csv"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

// establish server for rpc call
func main() {
	var qc lib.QueryCompletion
	queries := make([]string, 0, 10)
	scores := make([]int, 0, 10)

	// read data
	file, err := os.Open("src/query_completion/data/queries.csv") // For read access.
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

	// build index
	qc.BuildIndex(queries, scores, 10)
	log.Println("Query completion index built.")

	err = rpc.Register(&qc)
	if err != nil {
		log.Fatal(err)
	}

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	_ = http.Serve(l, nil)
}
