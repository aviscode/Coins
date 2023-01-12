package data

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const dataFileName = "coinsPrices.json"

type Coins struct {
	Kv *sync.Map `json:"coins"`
}

func NewCoinsDataStore() *Coins {
	return &Coins{Kv: &sync.Map{}}
}

type CoinData struct {
	Symbol   string  `json:"symbol"`
	UsdPrice float64 `json:"usd_price"`
}

func (c *Coins) UpdateCoins(coins []CoinData) {
	for _, coin := range coins {
		c.Kv.Store(coin.Symbol, coin.UsdPrice)
	}
	c.writeToFile(dataFileName)
}

func (c *Coins) GetSymbolPrice(symbol string) (float64, error) {
	if v, ok := c.Kv.Load(symbol); ok {
		return v.(float64), nil
	}
	return 0, fmt.Errorf("error loading symbol: %v", symbol)
}

func (c *Coins) writeToFile(path string) error {
	jsonData := make(map[string]interface{})
	c.Kv.Range(func(k interface{}, v interface{}) bool {
		jsonData[k.(string)] = v
		return true
	})
	data, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}
