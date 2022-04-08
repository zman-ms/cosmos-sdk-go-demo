package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"os"
)

var (
	dbName        = "bookstore"
	containerName = "books"
	sampleBook    = Book{
		Id:    "000000001",
		Title: "Computer Science",
		Price: 100.00,
	}
)

func GetClient() *azcosmos.Client {
	cosmosDbEndpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	cosmosDbKey, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}
	cred, _ := azcosmos.NewKeyCredential(cosmosDbKey)
	client, err := azcosmos.NewClientWithKey(cosmosDbEndpoint, cred, nil)

	if err != nil {
		panic(err)
	}
	return client
}

func getContainer() *azcosmos.ContainerClient {
	client := GetClient()
	container, err := client.NewContainer(dbName, containerName)
	if err != nil {
		panic(err)
	}
	return container
}

func InitializeDatabaseAndContainer() {
	client := GetClient()

	databaseProperties := azcosmos.DatabaseProperties{ID: dbName}

	databaseResponse, err := client.CreateDatabase(context.Background(), databaseProperties, nil)
	if err != nil {
		panic(err)
	}

	database, err := client.NewDatabase(dbName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Database created. ActivityId %s\r\n", databaseResponse.ActivityID)

	containerProperties := azcosmos.ContainerProperties{
		ID: "books",
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/title"},
		},
	}

	throughput := azcosmos.NewManualThroughputProperties(400)

	resp, err := database.CreateContainer(context.Background(), containerProperties, &azcosmos.CreateContainerOptions{ThroughputProperties: &throughput})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Container created. ActivityId %s\r\n", resp.ActivityID)

	fmt.Printf("\r\nAdding sample book entry to DB...\r\n")
	container := getContainer()

	pk := azcosmos.NewPartitionKeyString(sampleBook.Title)

	marshalled, err := json.Marshal(sampleBook)
	if err != nil {
		panic(err)
	}

	itemResponse, err := container.CreateItem(context.Background(), pk, marshalled, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sample book entry added. ActivityId %s consuming %v RU\r\n", itemResponse.ActivityID, itemResponse.RequestCharge)
}
