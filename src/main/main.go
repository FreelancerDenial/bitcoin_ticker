package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	BtceAdrr      = "https://btc-e.com/api/3/ticker/eur_usd-btc_usd-btc_eur"
	BtcChurtsAdrr = "http://api.bitcoincharts.com/v1/markets.json"
)

func BtceUpd() (btc_euro, btc_usd, eur_usd float64, ok bool) {
	var BtceList *btce_ticker

	client := &http.Client{}
	responce, err := client.Get(BtceAdrr)
	if err != nil {
		log.Printf("Getting err on sending msg to API: %v", err)
		return 0, 0, 0, false
	}
	defer responce.Body.Close()

	err = json.NewDecoder(responce.Body).Decode(&BtceList)
	if err != nil {
		log.Printf("Getting err on reading HTML responce from API: %v", err)
		return 0, 0, 0, false
	}

	return BtceList.BTCEuro.EuroCurrency, BtceList.BTCUsd.UsdCurrency, BtceList.EuroUsd.EuroUsdCurrency, true
}

func BtcChurtsUpd() (btc_euro, btc_usd, eur_usd float64, ok bool) {
	var BtcChurtsList []*btc_charts

	client := &http.Client{}
	responce, err := client.Get(BtcChurtsAdrr)
	if err != nil {
		log.Printf("Getting err on sending msg to API: %v", err)
		return 0, 0, 0, false
	}
	defer responce.Body.Close()

	err = json.NewDecoder(responce.Body).Decode(&BtcChurtsList)
	if err != nil {
		log.Printf("Getting err on reading HTML responce from API: %v", err)
		return 0, 0, 0, false
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

	return EurSum / EurCount, UsdSum / UsdCount, 0, true

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Start working!")

	a, b, c, flag := BtceUpd()
	fmt.Println(a, b, c, flag)

	a, b, c, flag = BtcChurtsUpd()
	fmt.Println(a, b, c, flag)
}
