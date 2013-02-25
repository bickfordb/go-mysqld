package main

import "mysqld"
import "os"
import "fmt"

func handleQuery(conn *mysqld.Conn, query string, rows chan map[string]interface{}, errors chan mysqld.Error) {
	if query == "baz" {
		rows <- map[string]interface{}{
			"column a": 1,
			"column b": "hey"}
	} else {
		errors <- mysqld.NotImplemented
	}
	defer close(rows)
	defer close(errors)
	return
}

func main() {
	server := mysqld.Server{}
	server.OnQuery = handleQuery
	err := server.Listen(3306)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
}
