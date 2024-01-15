package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/matthiasmohr/ed4-pricechanger-go/config"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/repository/adapter"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/repository/instance"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/routes"
	"github.com/matthiasmohr/ed4-pricechanger-go/internal/rules"
	RulesContract "github.com/matthiasmohr/ed4-pricechanger-go/internal/rules/contract"
	"github.com/matthiasmohr/ed4-pricechanger-go/utils/logger"
	"log"
	"net/http"
)

func main() {
	configs := config.GetConfig()
	connection := instance.GetConnection()
	repository := adapter.NewAdapter(connection)

	logger.INFO("waiting for the service to start...", nil)

	errors := migrate(connection)
	if len(errors) > 0 {
		for _, err := range errors {
			logger.PANIC("Error on migration:...", err)
		}
	}

	logger.PANIC("", checkTables(connection))

	port := fmt.Sprintf(":%v", configs.Port)
	router := routes.NewRouter().SetRouters(repository)
	logger.INFO("service is running on port: ", port)

	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func migrate(connection *dynamodb.DynamoDB) []error {
	var errors []error
	callMigrateAndAppendError(&errors, connection, &RulesContract.Rules{})
	return errors
}

func callMigrateAndAppendError(errors *[]error, connection *dynamodb.DynamoDB, rule rules.Interface) {
	err := rule.Migrate(connection)
	if err != nil {
		*errors = append(*errors, err)
	}
}
func checkTables(connection *dynamodb.DynamoDB) error {
	response, err := connection.ListTables(&dynamodb.ListTablesInput{})
	if response != nil {
		if len(response.TableNames) == 0 {
			logger.INFO("Tables not found: ", nil)
		}
		for _, tableName := range response.TableNames {
			logger.INFO("Table found: ", *tableName)
		}
	}
	return err
}
