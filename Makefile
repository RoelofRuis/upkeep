build:
	GOOS=linux go build -ldflags="-X 'main.mode=dev'" -o upkeep

build-dbg:
	GOOS=linux go build -ldflags="-X 'main.mode=dbg'" -o upkeep

build-prod:
	GOOS=linux go build -ldflags="-s -w" -o bin/linux/upkeep
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/osx/upkeep

install:
	GOOS=linux go install -ldflags="-s -w"