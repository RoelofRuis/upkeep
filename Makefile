build-dev:
	GOOS=linux go build -ldflags="-X 'main.mode=dev'" -o upkeep

build:
	GOOS=linux go build -ldflags="-s -w" -o upkeep