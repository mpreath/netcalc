all: go-test build

build: netcalc netcalc-api

netcalc:
	go build -o bin/netcalc ./cmd/netcalc/

netcalc-api:
	go build -o bin/netcalc-api ./cmd/netcalc-api/

go-test:
	go test ./...

run-api:
	go run ./cmd/netcalc-api/

clean:
	rm bin/*
	rmdir bin/