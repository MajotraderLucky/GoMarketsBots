package main

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
)

func main() {
	var (
		apiKey    = "mlxVdHIyxdRQNG4cTiDy8GMSf1Q10C4iwZQr7Uj2kfBL95KbO99vqK1PTtZLT4h2"
		secretKey = "6a4ADHGMbMAAD34Aw9kqNcTgXWQWLu1bY2u9SmIgWxRgK0jGoMYhCAHXedh3oyMc"
	)
	futuresClient := binance.NewFuturesClient(apiKey, secretKey)

	res, err := futuresClient.NewDepthService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}