build:
	GOOS=linux go build -ldflags="-X 'main.mode=dev'" -o up ./cmd/upkeep/

build-dbg:
	GOOS=linux go build -ldflags="-X 'main.mode=dbg'" -o up ./cmd/upkeep/

build-prod:
	GOOS=linux go build -ldflags="-s -w" -o bin/linux/upkeep ./cmd/upkeep/
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/osx/upkeep ./cmd/upkeep/

install:
	GOOS=linux go install -ldflags="-s -w" ./cmd/upkeep/