# Query Completion using Trie and Range Maximum Queries

[![Go Report Card](https://goreportcard.com/badge/github.com/wcy16/query_completion)](https://goreportcard.com/report/github.com/wcy16/query_completion) [![Build Status](https://travis-ci.org/wcy16/query_completion.svg?branch=master)](https://travis-ci.org/wcy16/query_completion) 

## Usage

1. put data under data/queries.csv

2. start the server

3. client uses RPC call to get top K related queries. Example:

```go
    client, err := rpc.DialHTTP("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }
 
    args := "a"
    var reply []string
    err = client.Call("QueryCompletion.Search", &args, &reply)
    if err != nil {
        log.Fatal("error:", err)
    }
    fmt.Println(reply)
```
