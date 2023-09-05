package main

import "time"

type TickData struct {
	Mode            string
	InstrumentToken uint32
	IsTradable      bool
	IsIndex         bool

	// Timestamp represents Exchange timestamp
	Timestamp          time.Time
	LastTradeTime      time.Time
	LastPrice          float64
	LastTradedQuantity uint32
	TotalBuyQuantity   uint32
	TotalSellQuantity  uint32
	VolumeTraded       uint32
	TotalBuy           uint32
	TotalSell          uint32
	AverageTradePrice  float64
	OI                 uint32
	OIDayHigh          uint32
	OIDayLow           uint32
	NetChange          float64
}
