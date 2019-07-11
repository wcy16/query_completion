package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type QueryCompletion struct {
	Trie    *Trie
	RMQ     *RMQ
	Queries []string
	K       int
}

// for rpc call
func (qc *QueryCompletion) Search(query *string, suggests *[]string) error {
	low, high := qc.Trie.Find([]rune(*query))
	topk := qc.RMQ.TopK(low, high, qc.K)

	//*suggests = make([]string, len(topk))
	//for i, qid := range topk {
	//	(*suggests)[i] = qc.Queries[qid]
	//}

	for _, qid := range topk {
		*suggests = append(*suggests, qc.Queries[qid])
	}

	return nil
}

// build trie and rmq for sorted queries and scores
func (qc *QueryCompletion) BuildIndex(queries []string, scores []int, topk int) {
	if len(queries) != len(scores) {
		log.Fatal("different length of queries and scores.")
	}
	qc.Queries = queries
	qc.RMQ = CreateRMQ(scores)
	qc.Trie = NewTrie(-1, -1)

	qc.K = topk

	for id, query := range queries {
		qc.Trie.Insert([]rune(query), id)
	}
}

// save rmq
func (qc *QueryCompletion) Save(path string) {
	if data, err := json.Marshal(qc); err == nil {
		if fp, err := os.Create(path); err == nil {
			_, _ = fp.Write(data)
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
}

// load rmq
func (qc *QueryCompletion) Load(path string) {
	if data, err := ioutil.ReadFile(path); err == nil {
		err = json.Unmarshal(data, qc)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}
