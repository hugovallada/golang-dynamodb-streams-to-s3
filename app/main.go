package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hugovallada/dynamodb-stream-processing-to-s3/src/serverless"
)

func main() {
	lambda.Start(serverless.HandleRequest)
}
