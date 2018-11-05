build:
	go build -ldflags="-s -w" -o bin/history api/history/*
	go build -ldflags="-s -w" -o bin/send api/send/*
	go build -ldflags="-s -w" -o bin/sender jobs/sender/*

test:
	go test ./api/*
	go test ./jobs/*

clean:
	rm -rf ./bin