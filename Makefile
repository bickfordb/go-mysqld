all:

sql: src/sql/parser.go
	GOPATH=$(realpath .) go run src/main.go

%.go: %.y
	GOPATH=$(realpath .) go tool yacc -o /tmp/x.go $+
	cp /tmp/x.go $@
	rm /tmp/x.go

main:
	GOPATH=$(realpath .) go run src/main.go

fmt:
	GOPATH=$(realpath .) go fmt mysqld sql
	GOPATH=$(realpath .) gofmt src/main.go >/tmp/fmt.go
	cp /tmp/fmt.go src/main.go


test:
	GOPATH=$(realpath .) go test -test.v mysqld


