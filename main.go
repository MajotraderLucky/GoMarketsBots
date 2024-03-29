package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
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

		klines, err := futuresClient.NewKlinesService().Symbol("BTCUSDT").Interval("1h").Do(context.Background())
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
		fmt.Println("1h open :", autoGeneratedKlines[498].Open)
		fmt.Println("1h close:", autoGeneratedKlines[498].Close)
		fmt.Println("1h high :", autoGeneratedKlines[498].High)
		fmt.Println("1h low  :", autoGeneratedKlines[498].Low)
		fmt.Println("----------------------")

		tStart := MillisecondsToTime(autoGeneratedKlines[0].CloseTime)
		fmt.Println("Start history:")
		fmt.Println(tStart)
		fmt.Println("1h open :", autoGeneratedKlines[0].Open)
		fmt.Println("1h close:", autoGeneratedKlines[0].Close)
		fmt.Println("1h high :", autoGeneratedKlines[0].High)
		fmt.Println("1h low  :", autoGeneratedKlines[0].Low)
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

		accountStart := 18.149229049682617 + 7.53667852 + 11.86 - 10
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
		fmt.Println()

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
		err = os.WriteFile("fileJson.json", fileJson, 0644)
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
		fmt.Println("The entry price position -", autoGeneratedpos.Positions[positionBTCindex].EntryPrice)
		fmt.Println("Position size", autoGeneratedpos.Positions[positionBTCindex].PositionAmt)
		fmt.Println("Item positions total -", len(autoGeneratedpos.Positions))
		fmt.Println("----------------------")

		var startTrade bool = false

		if priceCorridorPercent > 10 {
			fmt.Println(priceCorridorPercent, "> 10")
			fmt.Println("Corridor > 10 - you can trade")
			startTrade = true
		} else {
			fmt.Println("Corridor < 10 - you can't trade")
			startTrade = false
		}

		fmt.Println("Start trade =", startTrade)
		fmt.Println("----------------------")

		var askPriceFloat float64

		if askPriceFloat, err = strconv.ParseFloat(autoGenerated.Asks[0].Price, 32); err != nil {
			fmt.Println(err)
		}
		var priceAbove236 bool = false
		if (askPriceFloat > longFib236) && (askPriceFloat < max) {
			priceAbove236 = true
		} else {
			priceAbove236 = false
		}
		fmt.Println("Price above 236 fibo =", priceAbove236)

		var priceAbove382 bool = false
		if (askPriceFloat > longFib382) && (askPriceFloat < longFib236) {
			priceAbove382 = true
		} else {
			priceAbove382 = false
		}
		fmt.Println("Price above 382 fibo =", priceAbove382)

		var priceAbove500 bool = false
		if (askPriceFloat > longFib500) && (askPriceFloat < longFib382) {
			priceAbove500 = true
		} else {
			priceAbove500 = false
		}
		fmt.Println("Price above 500 fibo =", priceAbove500)

		var priceAbove618 bool = false
		if (askPriceFloat > longFib618) && (askPriceFloat < longFib500) {
			priceAbove618 = true
		} else {
			priceAbove618 = false
		}
		fmt.Println("Price above 618 fibo =", priceAbove618)

		var priceAbove786 bool = false
		if (askPriceFloat > longFib786) && (askPriceFloat < longFib618) {
			priceAbove786 = true
		} else {
			priceAbove786 = false
		}
		fmt.Println("Price above 786 fibo =", priceAbove786)
		fmt.Println("----------------------")

		var positionSizeFloat float64
		if positionSizeFloat, err = strconv.ParseFloat(autoGeneratedpos.Positions[positionBTCindex].PositionAmt, 32); err != nil {
			fmt.Println(err)
		}
		openPosition := false
		if positionSizeFloat != 0 {
			openPosition = true
		} else {
			openPosition = false
		}
		// Level 382 open orders
		var startTradeTo382 = false
		if priceAbove382 == true && startTrade == true && openPosition == false {
			startTradeTo382 = true
		} else {
			startTradeTo382 = false
		}
		fmt.Println("Start trade to level 382 =", startTradeTo382)

		openOrders, err := futuresClient.NewListOpenOrdersService().Symbol("BTCUSDT").
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, o := range openOrders {
			fmt.Println(o)
			fmt.Println(len(openOrders), "orders have been opened")
		}
		if len(openOrders) == 0 && startTradeTo382 {
			fmt.Println(len(openOrders), "orders have been opened")
			longFib382String := fmt.Sprintf("%.0f", longFib382)
			limitOrder, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").
				Side(futures.SideTypeBuy).Type(futures.OrderTypeLimit).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.001").
				Price(longFib382String).Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(limitOrder)
		}

		if len(openOrders) == 1 && startTradeTo382 == true {
			longFib500String := fmt.Sprintf("%.0f", longFib500)
			stopOrder, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").
				Side(futures.SideTypeSell).Type(futures.OrderTypeStopMarket).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.001").StopPrice(longFib500String).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(stopOrder)
		}

		var takeProfit382Condition bool = positionSizeFloat != 0 && len(openOrders) == 1 && priceAbove382 == true

		if takeProfit382Condition {
			longFib236String := fmt.Sprintf("%.0f", longFib236)
			takeProfitOrder, err := futuresClient.NewCreateOrderService().
				Symbol("BTCUSDT").Side(futures.SideTypeSell).Type(futures.OrderTypeTakeProfitMarket).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.001").StopPrice(longFib236String).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(takeProfitOrder)
		}

		// Level 500 open orders
		var startTradeTo500 = false
		if priceAbove500 == true && startTrade == true && openPosition == false {
			startTradeTo500 = true
		} else {
			startTradeTo500 = false
		}
		fmt.Println("Start trade to level 500 =", startTradeTo500)

		if len(openOrders) == 0 && startTradeTo500 == true {
			fmt.Println(len(openOrders), "orders have been opened")
			longFib500String := fmt.Sprintf("%.0f", longFib500)
			limitOrder, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").
				Side(futures.SideTypeBuy).Type(futures.OrderTypeLimit).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.002").
				Price(longFib500String).Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(limitOrder)
		}

		if len(openOrders) == 1 && startTradeTo500 == true {
			longFib618String := fmt.Sprintf("%.0f", longFib618)
			stopOrder, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").
				Side(futures.SideTypeSell).Type(futures.OrderTypeStopMarket).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.002").StopPrice(longFib618String).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(stopOrder)
		}

		var takeProfit500Condition bool = positionSizeFloat != 0 && len(openOrders) == 1 && priceAbove500 == true

		if takeProfit500Condition {
			longFib382String := fmt.Sprintf("%.0f", longFib382)
			takeProfitOrder, err := futuresClient.NewCreateOrderService().
				Symbol("BTCUSDT").Side(futures.SideTypeSell).Type(futures.OrderTypeTakeProfitMarket).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.002").StopPrice(longFib382String).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(takeProfitOrder)
		}

		// Level 618 open orders
		var startTradeTo618 = false
		if priceAbove618 == true && startTrade == true && openPosition == false {
			startTradeTo618 = true
		} else {
			startTradeTo618 = false
		}
		fmt.Println("Start trade to level 618 =", startTradeTo618)

		if len(openOrders) == 0 && startTradeTo618 == true {
			fmt.Println(len(openOrders), "orders have been opened")
			longFib618String := fmt.Sprintf("%.0f", longFib618)
			limitOrder, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").
				Side(futures.SideTypeBuy).Type(futures.OrderTypeLimit).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.003").
				Price(longFib618String).Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(limitOrder)
		}

		if len(openOrders) == 1 && startTradeTo618 == true {
			longFib786String := fmt.Sprintf("%.0f", longFib786)
			stopOrder, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").
				Side(futures.SideTypeSell).Type(futures.OrderTypeStopMarket).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.003").StopPrice(longFib786String).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(stopOrder)
		}

		var takeProfit618Condition bool = positionSizeFloat != 0 && len(openOrders) == 1 && priceAbove618 == true

		if takeProfit618Condition {
			longFib500String := fmt.Sprintf("%.0f", longFib500)
			takeProfitOrder, err := futuresClient.NewCreateOrderService().
				Symbol("BTCUSDT").Side(futures.SideTypeSell).Type(futures.OrderTypeTakeProfitMarket).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.003").StopPrice(longFib500String).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(takeProfitOrder)
		}

		// Level 786 open orders
		var startTradeTo786 = false
		if priceAbove786 == true && startTrade == true && openPosition == false {
			startTradeTo786 = true
		} else {
			startTradeTo786 = false
		}
		fmt.Println("Start trade to level 786 =", startTradeTo786)

		if len(openOrders) == 0 && startTradeTo786 == true {
			fmt.Println(len(openOrders), "orders have been opened")
			longFib786String := fmt.Sprintf("%.0f", longFib786)
			limitOrder, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").
				Side(futures.SideTypeBuy).Type(futures.OrderTypeLimit).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.005").
				Price(longFib786String).Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(limitOrder)
		}

		if len(openOrders) == 1 && startTradeTo786 == true {
			minString := fmt.Sprintf("%.0f", min)
			stopOrder, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").
				Side(futures.SideTypeSell).Type(futures.OrderTypeStopMarket).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.005").StopPrice(minString).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(stopOrder)
		}

		var takeProfit786Condition bool = positionSizeFloat != 0 && len(openOrders) == 1 && priceAbove786 == true

		if takeProfit786Condition {
			longFib618String := fmt.Sprintf("%.0f", longFib618)
			takeProfitOrder, err := futuresClient.NewCreateOrderService().
				Symbol("BTCUSDT").Side(futures.SideTypeSell).Type(futures.OrderTypeTakeProfitMarket).
				TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.005").StopPrice(longFib618String).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(takeProfitOrder)
		}

		var openOrderSize []string

		for _, o := range openOrders {
			openOrderSize = append(openOrderSize, o.OrigQuantity)
		}
		fmt.Println(openOrderSize)
		for x, y := range openOrders {
			if y != openOrders[x] {
				fmt.Println("The order sizes are not equal")
				err = futuresClient.NewCancelAllOpenOrdersService().
					Symbol("BTCUSDT").Do(context.Background())
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		fmt.Println("*******************************************")
	}
}
