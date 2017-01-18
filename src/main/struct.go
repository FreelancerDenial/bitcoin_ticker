package main

type btce_ticker struct {
	BTCEuro *btce_euro     `json:"btc_eur"`
	BTCUsd  *btce_usd      `json:"btc_usd"`
	EuroUsd *btce_euro_usd `json:"eur_usd"`
}

type btce_euro struct {
	EuroCurrency float64 `json:"avg"`
}

type btce_usd struct {
	UsdCurrency float64 `json:"avg"`
}

type btce_euro_usd struct {
	EuroUsdCurrency float64 `json:"avg"`
}

type btc_charts struct {
	Currency string  `json:"currency"`
	AvgPrice float64 `json:"avg"`
}

type fixer_currency struct {
	BaseMarker string      `json:"base"`
	Currency   *fixer_euro `json:"rates"`
}

type fixer_euro struct {
	UsdCurrency float64 `json:"USD"`
}
