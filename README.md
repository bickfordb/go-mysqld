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
```


```bash

[bran@bathysphere mysqld (master)]$ mysql -e 'baz'
+----------+----------+
| column a | column b |
+----------+----------+
|      hey |        1 |
|      you |        2 |
+----------+----------+
[bran@bathysphere mysqld (master)]$ mysql -e 'hey'
ERROR 1 (S1000) at line 1: Not Implemented

```

