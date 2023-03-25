package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// Convert milliseconds to time.Time
func MillisecondsToTime(milliseconds int64) time.Time {
	return time.Unix(0, milliseconds*int64(time.Millisecond))
}

func main() {
	apiKey, exists := os.LookupEnv("BINANCE_API_KEY")
	if exists {
		fmt.Println("apiKey exist")
	}

	secretKey, exexists := os.LookupEnv("BINANCE_SECRET_KEY")
	if exexists {
		fmt.Println("secretKey exist")
	}

	futuresClient := binance.NewFuturesClient(apiKey, secretKey)
	res, err := futuresClient.NewDepthService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(res)

	depthVar, _ := json.Marshal(res)
	// fmt.Println(string(depthVar))

	type AutoGenerated struct {
		LastUpdateID int64 `json:"lastUpdateId"`
		E            int64 `json:"E"`
		T            int64 `json:"T"`
		Bids         []struct {
			Price    string `json:"Price"`
			Quantity string `json:"Quantity"`
		} `json:"bids"`
		Asks []struct {
			Price    string `json:"Price"`
			Quantity string `json:"Quantity"`
		} `json:"asks"`
	}

	var autoGenerated AutoGenerated
	json.Unmarshal(depthVar, &autoGenerated)
	fmt.Println("----------------------")
	fmt.Println("----------------------")
	fmt.Println("ASK:", autoGenerated.Asks[0].Price, "-", autoGenerated.Asks[0].Quantity)
	fmt.Println("BID:", autoGenerated.Bids[0].Price, "-", autoGenerated.Bids[0].Quantity)
	fmt.Println("----------------------")

	klines, err := futuresClient.NewKlinesService().Symbol("BTCUSDT").Interval("15m").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	klinesVar, _ := json.Marshal(klines)

	type AutoGeneratedKlines []struct {
		OpenTime                 int64  `json:"openTime"`
		Open                     string `json:"open"`
		High                     string `json:"high"`
		Low                      string `json:"low"`
		Close                    string `json:"close"`
		Volume                   string `json:"volume"`
		CloseTime                int64  `json:"closeTime"`
		QuoteAssetVolume         string `json:"quoteAssetVolume"`
		TradeNum                 int    `json:"tradeNum"`
		TakerBuyBaseAssetVolume  string `json:"takerBuyBaseAssetVolume"`
		TakerBuyQuoteAssetVolume string `json:"takerBuyQuoteAssetVolume"`
	}

	var autoGeneratedKlines AutoGeneratedKlines
	json.Unmarshal(klinesVar, &autoGeneratedKlines)
	t := MillisecondsToTime(autoGeneratedKlines[498].CloseTime)
	fmt.Println("Last kline:")
	fmt.Println(t)
	fmt.Println("15min open :", autoGeneratedKlines[498].Open)
	fmt.Println("15min close:", autoGeneratedKlines[498].Close)
	fmt.Println("15min high :", autoGeneratedKlines[498].High)
	fmt.Println("15min low  :", autoGeneratedKlines[498].Low)
	fmt.Println("----------------------")

	tStart := MillisecondsToTime(autoGeneratedKlines[0].CloseTime)
	fmt.Println("Start history:")
	fmt.Println(tStart)
	fmt.Println("15min open :", autoGeneratedKlines[0].Open)
	fmt.Println("15min close:", autoGeneratedKlines[0].Close)
	fmt.Println("15min high :", autoGeneratedKlines[0].High)
	fmt.Println("15min low  :", autoGeneratedKlines[0].Low)
	fmt.Println("----------------------")

	resAcc, err := futuresClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(resAcc)

	accVar, _ := json.Marshal(resAcc)
	// fmt.Println(accVar)

	type Account struct {
		FeeTier                     int    `json:"feeTier"`
		CanTrade                    bool   `json:"canTrade"`
		CanDeposit                  bool   `json:"canDeposit"`
		CanWithdraw                 bool   `json:"canWithdraw"`
		UpdateTime                  int64  `json:"updateTime"`
		TotalInitialMargin          string `json:"totalInitialMargin"`
		TotalMaintMargin            string `json:"totalMaintMargin"`
		TotalWalletBalance          string `json:"totalWalletBalance"`
		TotalUnrealizedProfit       string `json:"totalUnrealizedProfit"`
		TotalMarginBalance          string `json:"totalMarginBalance"`
		TotalPositionInitialMargin  string `json:"totalPositionInitialMargin"`
		TotalOpenOrderInitialMargin string `json:"totalOpenOrderInitialMargin"`
		TotalCrossWalletBalance     string `json:"totalCrossWalletBalance"`
		TotalCrossUnPnl             string `json:"totalCrossUnPnl"`
		AvailableBalance            string `json:"availableBalance"`
		MaxWithdrawAmount           string `json:"maxWithdrawAmount"`
	}

	var account Account
	json.Unmarshal(accVar, &account)
	fmt.Println("----------------------")

	accountStart := 18.149229049682617 + 7.53667852
	accountNowString := account.AvailableBalance
	if accountNowFloat, err := strconv.ParseFloat(accountNowString, 32); err == nil {
		fmt.Println(accountStart, "- start")
		fmt.Println(accountNowFloat, "- now")
		fmt.Print("proffit($) = ", accountNowFloat-accountStart, "$", "\n")
		if accountNowFloat < accountStart {
			fmt.Print("proffit(%) = -", (accountNowFloat/accountStart)*100, "%")
		} else {
			fmt.Print("proffit(%) = ", (accountNowFloat/accountStart)*100, "%")
		}
	}
	fmt.Println("\n")
	startLowString := autoGeneratedKlines[0].Low
	var startLowFloat float64
	if s, err := strconv.ParseFloat(startLowString, 32); err == nil {
		startLowFloat = s
	}
	fmt.Println("Start kline low =", startLowFloat)

	// Make low slice float64
	lowSliceFloat64 := make([]float64, len(autoGeneratedKlines))
	fmt.Println(lowSliceFloat64)
}
