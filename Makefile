build:
	GOOS=linux go build -ldflags="-X 'main.mode=dev'" -o upkeep

build-dbg:
	GOOS=linux go build -ldflags="-X 'main.mode=dbg'" -o upkeep

build-prod:
	GOOS=linux go build -ldflags="-s -w" -o upkeep

install:
	GOOS=linux go install -ldflags="-s -w"