package bucket

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hugovallada/dynamodb-stream-processing-to-s3/src/models"
)

func SendToS3(ctx context.Context, payload models.Payload) {
	bucket := createBucket(ctx)
	putObject := createPutObject(payload)
	_, err := bucket.PutObject(ctx, &putObject)
	
	if err != nil {
		panic(err)
	}

}

func createPutObject(payload models.Payload) s3.PutObjectInput {
	return s3.PutObjectInput{
		Bucket: aws.String("testebucket-hlvl"),
		Key:    aws.String(fmt.Sprintf("%d.json", time.Now().UnixMilli())),
		Body:   bytes.NewReader(payload.Data),
	}
}

func createBucket(ctx context.Context) s3.Client {
	return *s3.NewFromConfig(createConfiguration(ctx))
}

func createConfiguration(ctx context.Context) aws.Config {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("Config not available")
	}
	return cfg
}
