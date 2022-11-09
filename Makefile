all:
	go get
	go get -t github.com/stretchr/testify
	go get -t  github.com/uptrace/bun/dialect/sqlitedialect
	go get -t  github.com/uptrace/bun/driver/sqliteshim
	go test
	go test -short -race
