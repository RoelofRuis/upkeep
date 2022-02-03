build-dev:
	GOOS=linux go build -ldflags="-X 'main.mode=dev'" -o ts

build:
	GOOS=linux go build -ldflags="-s -w" -o ts