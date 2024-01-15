package contract

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/entities"
	"time"
)

type Contract struct {
	entities.Base
	Name string `json:"Name"`
}

func InterfaceToModel(data interface{}) (instance *Contract, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}
	return instance, json.Unmarshal(bytes, &instance)
}

func (c *Contract) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"id": c.ID.String()}
}

func (c *Contract) TableName() string {
	return "contracts"
}

func (c *Contract) Bytes() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Contract) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"_id":       c.ID.String(),
		"name":      c.Name,
		"createdAt": c.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt": c.UpdatedAt.Format(entities.GetTimeFormat()),
	}
}

func ParseDynamoAttributeToStruct(response map[string]*dynamodb.AttributeValue) (c Contract, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return c, errors.New("Item not found")
	}
	for key, value := range response {
		if key == "_id" {
			c.ID, err = uuid.Parse(*value.S)
			if c.ID == uuid.Nil {
				err = errors.New("Item not found")
			}
			if key == "name" {
				c.Name = *value.S
			}
			if key == "createdAt" {
				c.CreatedAt, err = time.Parse(entities.GetTimeFormat(), *value.S)
			}
			if key == "updatedAt" {
				c.UpdatedAt, err = time.Parse(entities.GetTimeFormat(), *value.S)
			}
			if err != nil {
				return c, err
			}
		}
	}

	return c, nil
}
