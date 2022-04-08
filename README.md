# cosmos-sdk-go-demo
Sample codes of using Cosmos DB SDK in Go. Currently, this is just a demonstration of ETag handling.

Prerequisite:
1. Set up a cosmos DB account beforehand. This can be done through Azure Portal, Azure Client or even the SDK used by this demo. 
2. Set local environmental variables: 
   1. AZURE_COSMOS_ENDPOINT should be set to the URI of teh cosmos DB account;
   2. AZURE_COSMOS_KEY should be set to the primary read-write key of the cosmos DB account;
   
All functionalities regarding ETag handling are put in demo() of main.go

