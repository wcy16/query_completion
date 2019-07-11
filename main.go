package main

import (
	"./lib"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

var help = `Parameters
	build [-s | save]
	load
`

// establish server for rpc call
func main() {
	if len(os.Args) == 1 {
		fmt.Print(help)
		return
	}

	var qc lib.QueryCompletion

	if os.Args[1] == "build" {
		queries := make([]string, 0, 10)
		scores := make([]int, 0, 10)

		// read data
		file, err := os.Open("./data/queries.csv") // For read access.
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

		// save to file
		if len(os.Args) == 3 && os.Args[2] == "-s" {
			qc.Save("qc.dat")
			log.Println("Saved.")
		}
	} else if os.Args[1] == "load" {
		qc.Load("qc.dat")
		log.Println("Loaded.")
	} else {
		fmt.Print(help)
		return
	}

	err := rpc.Register(&qc)
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
