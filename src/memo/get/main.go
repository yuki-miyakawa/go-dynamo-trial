package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ResponseBody struct {
	MemberId string `json:"memberId"`
	Content  string `json:"content"`
}

func HandlerRequest(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var responseBody ResponseBody
	// AWS SDKの設定をロード
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("config load error: %v", err)
		return nil, err
	}

	// S3のクライアントを作成
	client := s3.NewFromConfig(cfg)
	data, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("dev-go-dynamo-trial"),
		Key:    aws.String("memo/" + event.PathParameters["memberId"] + ".json"),
	})
	if err != nil {
		log.Printf("s3 get object error: %v", err)
		return nil, err
	}
	defer data.Body.Close()

	// data.BodyをresponseBodyに入れる
	err = json.NewDecoder(data.Body).Decode(&responseBody)
	if err != nil {
		log.Printf("json decode error: %v", err)
		return nil, err
	}

	body, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(HandlerRequest)
}
