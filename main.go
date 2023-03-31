package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
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
	for _ = range time.Tick(time.Second * 30) {
		fmt.Println("----------------------")
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
		var nextLowFloat float64
		var lowSliceFloat64 []float64
		lowSliceFloat64 = append(lowSliceFloat64, startLowFloat)
		// fmt.Println(lowSliceFloat64)

		for i := 1; i < len(autoGeneratedKlines); i++ {
			nextLowString := autoGeneratedKlines[i].Low
			if s1, err := strconv.ParseFloat(nextLowString, 32); err == nil {
				nextLowFloat = s1
				lowSliceFloat64 = append(lowSliceFloat64, nextLowFloat)
			}
		}

		min := lowSliceFloat64[0]
		for _, number := range lowSliceFloat64 {
			if number < min {
				min = number
			}
		}

		fmt.Println("Lowest price    =", min)

		// Make high slice float64
		var nextHighFloat float64
		var highSliceFloat64 []float64

		for l := 0; l < len(autoGeneratedKlines); l++ {
			nextHighString := autoGeneratedKlines[l].High
			if s2, err := strconv.ParseFloat(nextHighString, 32); err == nil {
				nextHighFloat = s2
				highSliceFloat64 = append(highSliceFloat64, nextHighFloat)
			}
		}

		max := highSliceFloat64[0]
		for _, number := range highSliceFloat64 {
			if number > max {
				max = number
			}
		}

		fmt.Println("Highest price   =", max)
		fmt.Println("----------------------")

		longFib236 := max - ((max - min) * 0.236)
		fmt.Println("long Fibo 236 =", longFib236)
		longFib382 := max - ((max - min) * 0.382)
		fmt.Println("long Fibo 382 =", longFib382)
		longFib500 := max - ((max - min) * 0.500)
		fmt.Println("long Fibo 500 =", longFib500)
		longFib618 := max - ((max - min) * 0.618)
		fmt.Println("long Fibo 618 =", longFib618)
		longFib786 := max - ((max - min) * 0.786)
		fmt.Println("long Fibo 786 =", longFib786)

		longFib786String := fmt.Sprintf("%.0f", longFib786)
		fmt.Println(longFib786String)

		priceCorridor := max - min
		fmt.Println("----------------------")
		fmt.Println("Price corridor    =", priceCorridor)
		priceCorridorPercent := ((max - min) / max) * 100
		fmt.Print("Price corridor(%) = ", math.Round(priceCorridorPercent*100)/100, "%\n")
		fmt.Println("----------------------")

		accServ, err := futuresClient.NewGetAccountService().Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		accServVar, _ := json.Marshal(accServ)
		// fmt.Println(accServVar, reflect.TypeOf(accServVar))

		fileJson, err := json.Marshal(accServ)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("fileJson.json", fileJson, 0644)
		if err != nil {
			panic(err)
		}

		type AutoGeneratedPos struct {
			Assets []struct {
				Asset                  string `json:"asset"`
				InitialMargin          string `json:"initialMargin"`
				MaintMargin            string `json:"maintMargin"`
				MarginBalance          string `json:"marginBalance"`
				MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
				OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
				PositionInitialMargin  string `json:"positionInitialMargin"`
				UnrealizedProfit       string `json:"unrealizedProfit"`
				WalletBalance          string `json:"walletBalance"`
			} `json:"assets"`
			FeeTier                     int    `json:"feeTier"`
			CanTrade                    bool   `json:"canTrade"`
			CanDeposit                  bool   `json:"canDeposit"`
			CanWithdraw                 bool   `json:"canWithdraw"`
			UpdateTime                  int    `json:"updateTime"`
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
			Positions                   []struct {
				Isolated               bool   `json:"isolated"`
				Leverage               string `json:"leverage"`
				InitialMargin          string `json:"initialMargin"`
				MaintMargin            string `json:"maintMargin"`
				OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
				PositionInitialMargin  string `json:"positionInitialMargin"`
				Symbol                 string `json:"symbol"`
				UnrealizedProfit       string `json:"unrealizedProfit"`
				EntryPrice             string `json:"entryPrice"`
				MaxNotional            string `json:"maxNotional"`
				PositionSide           string `json:"positionSide"`
				PositionAmt            string `json:"positionAmt"`
				Notional               string `json:"notional"`
				IsolatedWallet         string `json:"isolatedWallet"`
				UpdateTime             int64  `json:"updateTime"`
			} `json:"positions"`
		}

		var autoGeneratedpos AutoGeneratedPos
		json.Unmarshal(accServVar, &autoGeneratedpos)

		var positionBTCindex int

		for k := 0; k < len(autoGeneratedpos.Positions); k++ {
			if autoGeneratedpos.Positions[k].Symbol == "BTCUSDT" {
				positionBTCindex = k
			}
		}
		fmt.Println("index position BTC -", positionBTCindex)
		fmt.Println("Unrealized profit =", autoGeneratedpos.TotalUnrealizedProfit)
		fmt.Println("The entry price position -", autoGeneratedpos.Positions[201].EntryPrice)
		fmt.Println("Item positions total -", len(autoGeneratedpos.Positions))
		fmt.Println("----------------------")

		var startTrade bool = false

		if priceCorridorPercent > 7 {
			fmt.Println(priceCorridorPercent, "> 7")
			fmt.Println("Corridor > 7 - you can trade")
			startTrade = true
		} else {
			fmt.Println("Corridor < 7 - you can't trade")
			startTrade = false
		}

		fmt.Println("Start trade =", startTrade)
		fmt.Println("----------------------")
	}
}
