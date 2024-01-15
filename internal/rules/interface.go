package rules

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"io"
)

type Interface interface {
	ConvertIoReaderToStruct(data io.Reader, model interface{}) (body interface{}, er error)
	GetMock() interface{}
	Migrate(connection *dynamodb.DynamoDB) error
	Validate(model interface{}) error
}
