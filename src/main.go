package main

import "mysqld"
import "os"
import "fmt"

func main() {
	server := mysqld.NewServer()
	err := server.Listen(":3306")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listening: %s\n", err.Error())
		return
	}
	for {
		query := <-server.Queries
		if query.Statement == "baz" {
			query.WriteRow(map[string]interface{}{"column a": "hey", "column b": 1})
			query.WriteRow(map[string]interface{}{"column a": "you", "column b": 2})
			query.Finish(nil)
		} else {
			query.Finish(mysqld.NotImplemented)
		}
	}
}
