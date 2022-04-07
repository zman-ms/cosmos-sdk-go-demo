package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func main() {
	//CreateDatabaseAndContainer()
	demo()
}

func demo() {
	book := Book{
		Id:    "000000001",
		Title: "Computer Science",
		Price: 100.00,
	}
	//addBook(book)
	receivedBook1, etag1 := readBookInfo(book)
	updateBookPrice(receivedBook1, 124.00, etag1)
	updateBookPrice(receivedBook1, 126.00, etag1)
	receivedBook2, etag2 := readBookInfo(book)
	updateBookPrice(receivedBook2, 126.00, etag2)
}

func addBook(book Book) {
	fmt.Printf("\r\nAdding book entry to DB...\r\n")
	container := getContainer()

	pk := azcosmos.NewPartitionKeyString(book.Title)

	marshalled, err := json.Marshal(book)
	if err != nil {
		panic(err)
	}

	itemResponse, err := container.CreateItem(context.Background(), pk, marshalled, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Book entry added. ActivityId %s consuming %v RU\r\n", itemResponse.ActivityID, itemResponse.RequestCharge)
}

func readBookInfo(book Book) (Book, azcore.ETag) {
	fmt.Printf("\r\nReading book info...\r\n")
	container := getContainer()
	pk := azcosmos.NewPartitionKeyString(book.Title)
	itemResponse, err := container.ReadItem(context.Background(), pk, string(book.Id), nil)

	var receivedBook Book
	err = json.Unmarshal(itemResponse.Value, &receivedBook)
	if err != nil {
		panic(err)
	}

	jsonStr, err := json.Marshal(receivedBook)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Printf("Book info read=%s\r\nEtag=%s\r\nActivityId %s consuming %v RU\r\n", string(jsonStr), itemResponse.ETag, itemResponse.ActivityID, itemResponse.RequestCharge)
	return receivedBook, itemResponse.ETag
}

func updateBookPrice(book Book, newPrice float32, etag azcore.ETag) {
	fmt.Printf("\r\nUpdating book price...\r\n")
	container := getContainer()
	book.Price = newPrice
	marshalledBook, err := json.Marshal(book)
	if err != nil {
		panic(err)
	}

	// Replace with Etag
	pk := azcosmos.NewPartitionKeyString(book.Title)
	itemResponse, err := container.ReplaceItem(context.Background(), pk, string(book.Id), marshalledBook, &azcosmos.ItemOptions{IfMatchEtag: &etag})
	if err != nil {
		fmt.Printf("Update rejected.\r\nError is\r\n%s\r\n", err)

	} else {
		fmt.Printf("Book price updated.\r\nEtag=%s\r\nActivityId %s consuming %v RU\r\n", itemResponse.ETag, itemResponse.ActivityID, itemResponse.RequestCharge)
	}
}
