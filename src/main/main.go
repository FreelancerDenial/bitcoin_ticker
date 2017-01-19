package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	BtceAddr      = "https://btc-e.com/api/3/ticker/eur_usd-btc_usd-btc_eur"
	BtcChurtsAddr = "http://api.bitcoincharts.com/v1/markets.json"
	FixerAddr     = "http://api.fixer.io/latest?base=EUR"
	SourceCount   = 2
)

func BtceUpd() (btc_euro, btc_usd, eur_usd float64, BTCAlive, CurrAlive int64) {
	var BtceList *btce_ticker

	client := &http.Client{}
	responce, err := client.Get(BtceAddr)
	if err != nil {
		log.Printf("Getting err on sending msg to API: %v", err)
		return 0, 0, 0, 0, 0
	}
	defer responce.Body.Close()

	err = json.NewDecoder(responce.Body).Decode(&BtceList)
	if err != nil {
		log.Printf("Getting err on reading HTML responce from API: %v", err)
		return 0, 0, 0, 0, 0
	}

	return BtceList.BTCEuro.EuroCurrency, BtceList.BTCUsd.UsdCurrency, BtceList.EuroUsd.EuroUsdCurrency, 1, 1
}

func BtcChurtsUpd() (btc_euro, btc_usd float64, BTCAlive int64) {
	var BtcChurtsList []*btc_charts

	client := &http.Client{}
	responce, err := client.Get(BtcChurtsAddr)
	if err != nil {
		log.Printf("Getting err on sending msg to API: %v", err)
		return 0, 0, 0
	}
	defer responce.Body.Close()

	err = json.NewDecoder(responce.Body).Decode(&BtcChurtsList)
	if err != nil {
		log.Printf("Getting err on reading HTML responce from API: %v", err)
		return 0, 0, 0
	}

	var UsdSum, EurSum float64
	var UsdCount, EurCount float64
	for _, currency := range BtcChurtsList {
		switch currency.Currency {
		case "USD":
			if currency.AvgPrice != 0 {
				UsdSum += currency.AvgPrice
				UsdCount++
			}
		case "EUR":
			if currency.AvgPrice != 0 {
				EurSum += currency.AvgPrice
				EurCount++
			}
		}
	}

	return EurSum / EurCount, UsdSum / UsdCount, 1
}

func FixerCurrency() (eur_usd float64, CurrAlive int64) {
	var CurrencyList *fixer_currency

	client := &http.Client{}
	responce, err := client.Get(FixerAddr)
	if err != nil {
		log.Printf("Getting err on sending msg to API: %v", err)
		return 0, 0
	}
	defer responce.Body.Close()

	err = json.NewDecoder(responce.Body).Decode(&CurrencyList)
	if err != nil {
		log.Printf("Getting err on reading HTML responce from API: %v", err)
		return 0, 0
	}

	return CurrencyList.Currency.UsdCurrency, 1
}

func Ticker() {
	var BTCEuroCurr, BTCUsdCurr, EuroUsdCurr float64

	BTCEuroSourceOne, BTCUsdSourceOne, EuroUsdSourceOne, CountAliveBTCOne, CountAliveCurrOne := BtceUpd()

	BTCEuroSourceTwo, BTCUsdSourceTwo, CountAliveBTCTwo := BtcChurtsUpd()

	EuroUsdSourceTwo, CountAliveCurrTwo := FixerCurrency()

	if BTCEuroSourceOne >= BTCEuroSourceTwo {
		BTCEuroCurr = BTCEuroSourceOne
	} else {
		BTCEuroCurr = BTCEuroSourceTwo
	}

	if BTCUsdSourceOne >= BTCUsdSourceTwo {
		BTCUsdCurr = BTCUsdSourceOne
	} else {
		BTCUsdCurr = BTCUsdSourceTwo
	}

	if EuroUsdSourceOne >= EuroUsdSourceTwo {
		EuroUsdCurr = EuroUsdSourceOne
	} else {
		EuroUsdCurr = EuroUsdSourceTwo
	}

	fmt.Printf("Greater exchange rate:\nBTC/USD: %v EUR/USD: %v BTC/EUR: %v", BTCUsdCurr, EuroUsdCurr, BTCEuroCurr)
	fmt.Printf("\nActive sources: BTC/USD (%v of %v) EUR/USD (%v of %v) BTC/EUR (%v of %v)", CountAliveBTCOne+CountAliveBTCTwo, SourceCount, CountAliveCurrOne+CountAliveCurrTwo, SourceCount, CountAliveBTCOne+CountAliveBTCTwo, SourceCount)
	fmt.Println("\n ---------------------------------")

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	tick := time.NewTicker(time.Second * 5)
	for _ = range tick.C {
		Ticker()
	}
}
