package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type RequestBody struct {
	Content string `json:"content"`
}

type S3Payload struct {
	MemberId string `json:"memberId"`
	Content  string `json:"content"`
}

func HandlerRequest(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("event.Body: %v", event.Body)
	log.Printf("event.PathParameters: %v", event.PathParameters)
	var requestBody RequestBody
	err := json.Unmarshal([]byte(event.Body), &requestBody)
	if err != nil {
		log.Printf("json unmarshal error: %v", err)
		return nil, err
	}

	s3Payload := S3Payload{
		MemberId: event.PathParameters["memberId"],
		Content:  requestBody.Content,
	}

	s3PayloadJson, err := json.Marshal(s3Payload)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		return nil, err
	}

	body := string(s3PayloadJson)

	// AWS SDKの設定をロード
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("config load error: %v", err)
		return nil, err
	}

	// S3のクライアントを作成
	client := s3.NewFromConfig(cfg)

	// bodyにevent.PathParameters["memberId"]とevent.Bodyを結合したものを入れる

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("dev-go-dynamo-trial"),
		Key:    aws.String("memo/" + event.PathParameters["memberId"] + ".json"),
		Body:   strings.NewReader(body),
	})
	if err != nil {
		log.Printf("s3 put object error: %v", err)
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "ok",
	}, nil
}

func main() {
	lambda.Start(HandlerRequest)

}
