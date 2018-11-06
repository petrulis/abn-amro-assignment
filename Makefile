build:
	go build -ldflags="-s -w" -o bin/history api/history/*
	go build -ldflags="-s -w" -o bin/send api/send/*
	go build -ldflags="-s -w" -o bin/scanner jobs/scanner/*
	go build -ldflags="-s -w" -o bin/sender jobs/sender/*

test:
	go test ./api/...
	go test ./dynamodbdriver/...
	go test ./jobs/...
	go test ./model/...
	go test ./validator/...

clean:
	rm -rf ./bin