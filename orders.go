// longFib500String := fmt.Sprintf("%.0f", longFib500)
// minString := fmt.Sprintf("%.0f", min)

// limitOrder, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").
// 	Side(futures.SideTypeBuy).Type(futures.OrderTypeLimit).
// 	TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.001").
// 	Price(longFib786String).Do(context.Background())
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
// fmt.Println(limitOrder)

// stopOrder, err := futuresClient.NewCreateOrderService().Symbol("BTCUSDT").
// 	Side(futures.SideTypeSell).Type(futures.OrderTypeStopMarket).
// 	TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.001").StopPrice(minString).
// 	Do(context.Background())
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
// fmt.Println(stopOrder)

// takeProfitOrder, err := futuresClient.NewCreateOrderService().
// 	Symbol("BTCUSDT").Side(futures.SideTypeSell).Type(futures.OrderTypeTakeProfitMarket).
// 	TimeInForce(futures.TimeInForceTypeGTC).Quantity("0.001").StopPrice(longFib500String).
// 	Do(context.Background())
// if err != nil {
// 	fmt.Println(err)
// }
// fmt.Println(takeProfitOrder)

// _, err = futuresClient.NewCancelAllOpenOrdersService().
// 	Symbol("BTCUSDT").Do(context.Background())
// if err != nil {
// 	fmt.Println(err)
// 	return
// }

// err = futuresClient.NewCancelAllOpenOrdersService().
// 	Symbol("BTCUSDT").Do(context.Background())
// if err != nil {
// 	fmt.Println(err)
// }

// openOrders, err := futuresClient.NewListOpenOrdersService().Symbol("BTCUSDT").
// 	Do(context.Background())
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
// for _, o := range openOrders {
// 	fmt.Println(o.OrderID, reflect.TypeOf(o.OrderID))
// 	fmt.Println(o.Price)
// }

// fmt.Println(futuresClient.NewGetAccountService().Do(context.Background()))