go-mysqld
=========

Library for authoring daemons which speak the MySQL protocol

Example
-------

```go

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
```


```bash

[bran@bathysphere ~]$ mysql -e 'baz'
+----------+----------+
| column a | column b |
+----------+----------+
|        1 |      hey |
+----------+----------+
[bran@bathysphere ~]$ mysql -e 'hey'
ERROR 1 (S1000) at line 1: Not Implemented

```

