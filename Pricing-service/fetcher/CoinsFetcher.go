package fetcher

import (
	"errors"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
)

var (
	ErrEmptyApiKey = errors.New("error API KEY cannot be empty")
)

type ClientApi struct {
	*cmc.Client
}

func NewFetcher(apiKey string) (*ClientApi, error) {
	if apiKey == "" {
		return nil, ErrEmptyApiKey
	}
	return &ClientApi{cmc.NewClient(&cmc.Config{ProAPIKey: apiKey})}, nil
}

func (c *ClientApi) FetchData() ([]*cmc.Listing, error) {
	return c.Cryptocurrency.LatestListings(&cmc.ListingOptions{
		Start:   1,
		Limit:   3,
		Convert: "USD",
	})
}
