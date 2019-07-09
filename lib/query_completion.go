package lib

import "log"

type QueryCompletion struct {
	trie    *Trie
	rmq     *RMQ
	queries []string
	k       int
}

// for rpc call
func (qc *QueryCompletion) Search(query *string, suggests *[]string) error {
	low, high := qc.trie.Find([]rune(*query))
	topk := qc.rmq.TopK(low, high, qc.k)

	//*suggests = make([]string, len(topk))
	//for i, qid := range topk {
	//	(*suggests)[i] = qc.queries[qid]
	//}

	for _, qid := range topk {
		*suggests = append(*suggests, qc.queries[qid])
	}

	return nil
}

// build trie and rmq for sorted queries and scores
func (qc *QueryCompletion) BuildIndex(queries []string, scores []int, topk int) {
	if len(queries) != len(scores) {
		log.Fatal("different length of queries and scores.")
	}
	qc.queries = queries
	qc.rmq = CreateRMQ(scores)
	qc.trie = NewTrie(-1, -1)

	qc.k = topk

	for id, query := range queries {
		qc.trie.Insert([]rune(query), id)
	}
}
