package contract

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/entities/contract"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/repository/adapter"
	"time"
)

type Controller struct {
	repository adapter.Interface
}

type Interface interface {
	ListOne(id uuid.UUID) (entity contract.Contract, err error)
	ListAll() (entities []contract.Contract, err error)
	Create(entity *contract.Contract) (uuid.UUID, error)
	Update(id uuid.UUID, entity *contract.Contract) error
	Remove(id uuid.UUID) error
}

func NewController(repository adapter.Interface) Interface {
	return &Controller{
		repository: repository,
	}
}

func (c *Controller) ListOne(id uuid.UUID) (entity contract.Contract, err error) {
	entity.ID = id
	response, err := c.repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		return entity, err
	}
	return contract.ParseDynamoAttributeToStruct(response.Item)
}

func (c *Controller) ListAll() (entities []contract.Contract, err error) {
	entities = []contract.Contract{}
	var entity contract.Contract

	filter := expression.Name("name").NotEqual(expression.Value(""))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()

	if err != nil {
		return entities, err
	}

	response, err := c.repository.FindAll(condition, entity.TableName())
	if err != nil {
		return entities, err
	}

	if response != nil {
		for _, value := range response.Items {
			entity, err := contract.ParseDynamoAttributeToStruct(value)
			if err != nil {
				return entities, err
			}
			entities = append(entities, entity)
		}
	}
	return entities, nil
}
func (c *Controller) Create(entity *contract.Contract) (uuid.UUID, error) {
	entity.CreatedAt = time.Now()
	_, err := c.repository.CreateOrUpdate(entity.GetMap(), entity.TableName())
	return entity.ID, err
}

func (c *Controller) Update(id uuid.UUID, entity *contract.Contract) error {
	found, err := c.ListOne(id)
	if err != nil {
		return err
	}

	found.ID = id
	found.Name = entity.Name
	found.UpdatedAt = time.Now()
	_, err = c.repository.CreateOrUpdate(found.GetMap(), entity.TableName())
	return err
}

func (c *Controller) Remove(id uuid.UUID) error {
	entity, err := c.ListOne(id)
	if err != nil {
		return err
	}

	_, err = c.repository.Delete(entity.GetFilterId(), entity.TableName())
	return err
}
