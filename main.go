package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Image map[string]any

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event events.DynamoDBEvent) {
	jsons := []Image{}

	for _, record := range event.Records {
		if record.EventName == "INSERT" {
			continue
		}

		image := make(Image)
		for key, value := range record.Change.OldImage {
			dict := map[string]any{}
			marshalledJson, _ := value.MarshalJSON()
			json.Unmarshal(marshalledJson, &dict)
			for _, v := range dict {
				image[key] = v
			}
		}

		jsons = append(jsons, image)
	}

	listaJsons, _ := json.Marshal(jsons)

	timeNow := time.Now().UnixMilli()

	bucket := s3.PutObjectInput{
		Bucket: aws.String("testebucket-hlvl"),
		Key:    aws.String(fmt.Sprintf("%d.json", timeNow)),
		Body:   bytes.NewReader(listaJsons),
	}

	if len(jsons) > 0 {
		fmt.Println("Sending")

		cfg, _ := config.LoadDefaultConfig(ctx)
		c := s3.NewFromConfig(cfg)
		resp, err := c.PutObject(ctx, &bucket)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(resp.ResultMetadata)
	}
}
