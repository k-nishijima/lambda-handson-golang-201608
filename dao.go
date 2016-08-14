package lambdaHandson

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type DynamoDBClient struct{}

type AddValueRequest struct {
	Stage   string `json:"stage" valid:"required"`
	Email   string `json:"email" valid:"email,length(1|512),required"`
	Message string `json:"message" valid:"length(1|1024),required"`
}

type ContactItem struct {
	Timestamp string `json:"timestamp"`
	Email     string `json:"email"`
	Message   string `json:"message"`
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("configuration file load failed / " + err.Error())
	}
}

func (self DynamoDBClient) ValidateRequest(request interface{}) error {
	_, err := govalidator.ValidateStruct(request)
	return err
}

func (dao DynamoDBClient) svc() *dynamodb.DynamoDB {
	var cfg *aws.Config
	profile := viper.GetString("dynamo.profile")
	if profile == "" {
		cfg = aws.NewConfig().WithRegion(viper.GetString("dynamo.region")).WithEndpoint(viper.GetString("dynamo.endpoint"))
	} else {
		cred := credentials.NewSharedCredentials("", viper.GetString("dynamo.profile"))
		cfg = aws.NewConfig().WithRegion(viper.GetString("dynamo.region")).WithCredentials(cred).WithEndpoint(viper.GetString("dynamo.endpoint"))
	}

	return dynamodb.New(session.New(cfg), aws.NewConfig().WithRegion(viper.GetString("dynamo.region")))
}

func (dao DynamoDBClient) Put(request AddValueRequest) error {
	svc := dao.svc()
	params := &dynamodb.PutItemInput{
		TableName: aws.String(viper.GetString(request.Stage + ".contactTable")),
		Item: map[string]*dynamodb.AttributeValue{
			"Timestamp": {
				S: aws.String(time.Now().String()),
			},
			"Email": {
				S: aws.String(request.Email),
			},
			"Message": {
				S: aws.String(request.Message),
			},
		},
	}
	_, err := svc.PutItem(params)

	if err != nil {
		return errors.Wrap(err, "DynamoDB access error")
	}

	return nil
}

func (dao DynamoDBClient) GetItems(request AddValueRequest) ([]ContactItem, error) {
	svc := dao.svc()

	params := &dynamodb.ScanInput{
		TableName: aws.String(viper.GetString(request.Stage + ".contactTable")),
		AttributesToGet: []*string{
			aws.String("Timestamp"),
			aws.String("Email"),
			aws.String("Message"),
		},
		Limit: aws.Int64(500),
	}
	resp, err := svc.Scan(params)

	if err != nil {
		return nil, errors.Wrap(err, "DynamoDB access error")
	}

	var values []ContactItem
	for i := range resp.Items {
		var v ContactItem
		for key, item := range resp.Items[i] {
			switch key {
			case "Timestamp":
				v.Timestamp = *item.S
			case "Email":
				v.Email = *item.S
			case "Message":
				v.Message = *item.S
			}
		}
		values = append(values, v)
	}

	return values, nil
}
