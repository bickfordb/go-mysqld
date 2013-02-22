package main

import "mysqld"
import "os"
import "fmt"

func F(conn *mysqld.Conn, query string, rows chan map[string]interface{}) {
  row := map[string]interface{}{
    "foo": 1,
    "bar": "hey"}
  rows <- row
  close(rows)
  return
}

func main() {
  server := mysqld.Server{}
  server.OnQuery = F
  err := server.Listen(3306)
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
  }
}

