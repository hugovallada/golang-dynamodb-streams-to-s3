package serverless

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/hugovallada/dynamodb-stream-processing-to-s3/src/bucket"
	"github.com/hugovallada/dynamodb-stream-processing-to-s3/src/models"
)

func HandleRequest(ctx context.Context, event events.DynamoDBEvent) {
	payload := processStream(ctx, event)
	if payload.SendToS3 {
		bucket.SendToS3(ctx, payload)
	}
}

func processStream(ctx context.Context, event events.DynamoDBEvent) models.Payload {
	images := generateListOfImages(event.Records)
	data := generateJSONStringFromList(images)
	payload := generatePayload(data)
	return payload
}

func generatePayload(data []byte) models.Payload {
	return models.Payload{Data: data, SendToS3: len(data) > 2}
}

func generateJSONStringFromList(images []models.Image) []byte {
	data, err := json.Marshal(images)
	if err != nil {
		panic(err)
	}
	return data
}

func generateListOfImages(records []events.DynamoDBEventRecord) []models.Image {
	images := make([]models.Image, 0)
	for _, record := range records {
		if isRemove(record) {
			jsonForImage := mountJSON(record.Change)
			images = append(images, jsonForImage)
		}
	}
	return images
}

func isRemove(record events.DynamoDBEventRecord) bool {
	return record.EventName == "REMOVE"
}

func mountJSON(record events.DynamoDBStreamRecord) map[string]any {
	image := make(models.Image)
	for key, value := range record.OldImage {
		nestedMap := make(map[string]any)
		marshalledJson, err := value.MarshalJSON()
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(marshalledJson, &nestedMap)
		if err != nil {
			panic(err)
		}
		for _, val := range nestedMap {
			image[key] = val
		}
	}
	return image
}
