package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

func main() {
	ctx := context.Background()

	sa := flag.String("s", "", "name of the azure storageaccount")

	flag.Parse()

	storageAccountName := *sa
	if storageAccountName == "" {
		panic("Flag storageAccountName missing.")
	}

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccountName)

	credential, _ := azidentity.NewDefaultAzureCredential(nil)
	client, _ := azblob.NewClient(url, credential, nil)

	containerCount, err := getContainerCount(ctx, client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d containers found.", containerCount)
}

// return the number of cointainers in a storage account
func getContainerCount(_ctx context.Context, client *azblob.Client) (int64, error) {
	ctx, cancel := context.WithTimeout(_ctx, time.Duration(time.Millisecond*5000))
	defer cancel()

	var count int64

	maxResults := int32(100)
	opts := &azblob.ListContainersOptions{
		Include:    service.ListContainersInclude{
			Metadata: false,
			Deleted:  false,
			System:   false,
		},
		MaxResults: &maxResults,
	}
	pager := client.NewListContainersPager(opts)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err == nil {
			count = count + int64(len(page.ContainerItems))
		} else {
			return 0, err
		}
	}

	return count, nil
}
