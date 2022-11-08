all:
	go get
	go get github.com/stretchr/testify
	go get github.com/uptrace/bun/dialect/sqlitedialect
	go test
	go test -short -race
