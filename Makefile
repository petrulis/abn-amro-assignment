build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/history history/*
	env GOOS=linux go build -ldflags="-s -w" -o bin/send send/*

test:
	go test ./history
	go test ./send

deploy:

clean:
	rm -rf ./bin