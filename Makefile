all: build

build:
	go build ./cmd/netcalc/
	go build ./cmd/netcalc-api/

clean:
	rm netcalc
	rm netcalc-api