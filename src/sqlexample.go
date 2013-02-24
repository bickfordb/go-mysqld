package main

import "fmt"
import "sql"

func run(s string) {
  stmt, err := sql.Parse(s)
  if err != nil {
    println("error:", err.Error())
  } else {
    fmt.Printf("stmt: %+v\n", stmt)
  }
}

func main() {
  run("select 1,2")
}

