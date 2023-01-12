package collector

import (
	d "github.com/aviscode/Coins/Pricing-service/data"
	"github.com/aviscode/Coins/Pricing-service/fetcher"
	"github.com/sirupsen/logrus"
	"time"
)

func NewCollector(store *d.Coins, fetcher *fetcher.ClientApi) {
	coins := make([]d.CoinData, 3, 3)
	for {
		data, err := fetcher.FetchData()
		if err != nil {
			logrus.Errorf("FetchData() failed with err: %+v", err)
		}
		for i, coin := range data {
			logrus.Infof("FetchCoin %s -> %.3f", coin.Symbol, coin.Quote["USD"].Price)
			coins[i] = d.CoinData{Symbol: coin.Symbol, UsdPrice: coin.Quote["USD"].Price}
		}
		store.UpdateCoins(coins)
		time.Sleep(time.Minute)
	}
}
