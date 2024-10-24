package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"stori-challenge-v1/infrastructure/controller"
)

func main() {
	fmt.Println("Starting lambda")
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	lambda.Start(controller.HandleRequest)
}
