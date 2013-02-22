all:
	GOPATH=$(realpath .) go run src/main.go


test:
	GOPATH=$(realpath .) go test -test.v mysqld


