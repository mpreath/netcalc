all: build

build:
	go build ./cmd/netcalc/
	go build ./cmd/netcalc-server/

clean:
	rm netcalc
	rm netcalc-server