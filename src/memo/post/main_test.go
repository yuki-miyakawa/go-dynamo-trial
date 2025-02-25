package main

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// AWS SDKのクライアントのインターフェースを定義します。
// 実際のAWSサービスとの通信を行わず、モックの挙動を定義することができます。
type MockS3Client struct {
	PutObjectFunc func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func (m *MockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m.PutObjectFunc(ctx, params)
}

func TestHandlerRequest(t *testing.T) {
	tests := []struct {
		name          string
		event         events.APIGatewayProxyRequest
		want          int
		err           error
		putObjectFunc func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	}{
		{
			name: "successful invocation",
			event: events.APIGatewayProxyRequest{
				Body: `{"Content":"test content"}`,
				PathParameters: map[string]string{
					"memberId": "12345",
				},
			},
			want: http.StatusOK,
			err:  nil,
			putObjectFunc: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				// バケット名とオブジェクトキーが期待通りかチェックします。
				if *params.Bucket != "dev-go-dynamo-trial" {
					return nil, errors.New("wrong bucket name")
				}
				if *params.Key != "memo/12345.json" {
					return nil, errors.New("wrong key name")
				}
				return &s3.PutObjectOutput{}, nil // 成功したとみなし、空のレスポンスを返します。
			},
		},

		// テストケースを追加...
	}

	// 各テストケースを繰り返し実行します。
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// LoadDefaultConfigとS3クライアントのモックを作成します。
			oldLoadDefaultConfig := loadDefaultConfig
			oldNewFromConfig := newFromConfig
			defer func() {
				loadDefaultConfig = oldLoadDefaultConfig
				newFromConfig = oldNewFromConfig
			}()
			loadDefaultConfig = func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
				return aws.Config{}, nil
			}
			newFromConfig = func(cfg aws.Config, optFns ...func(*s3.Options)) *s3.Client {
				return &s3.Client{APIOptions: []func(*s3.Options){}}
			}

			// カスタムモック関数でS3クライアントを置き換えます。
			client = &MockS3Client{
				PutObjectFunc: tt.putObjectFunc,
			}

			// HandlerRequest関数を実行します。
			resp, err := HandlerRequest(context.Background(), tt.event)
			if !errors.Is(err, tt.err) {
				t.Errorf("HandlerRequest() error = %v, wantErr %v", err, tt.err)
				return
			}
			if resp.StatusCode != tt.want {
				t.Errorf("HandlerRequest() = %v, want %v", resp.StatusCode, tt.want)
			}
		})
	}
}
