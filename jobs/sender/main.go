package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler() error {
	fmt.Println("Running...")
	return nil
}

func main() {
	lambda.Start(Handler)
}
