sql: src/sql/parser.go
	GOPATH=$(realpath .) go run src/sqlexample.go

%.go: %.y
	GOPATH=$(realpath .) go tool yacc -o /tmp/x.go $+
	cp /tmp/x.go $@

all:
	GOPATH=$(realpath .) go run src/main.go


test:
	GOPATH=$(realpath .) go test -test.v mysqld


