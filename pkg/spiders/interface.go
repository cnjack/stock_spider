package spiders

import "time"

type KLine struct {
	Open  float64   `json:"open"`
	Close float64   `json:"close"`
	High  float64   `json:"high"`
	Low   float64   `json:"low"`
	Time  time.Time `json:"time"`
	Type  Type      `json:"type"`
}

type Trend struct {
	Time    time.Time `json:"time"`
	Price   float64   `json:"price"`
	Volume  int64     `json:"volume"`
	Incrace float64   `json:"incrace"`
}

type Type string

const (
	FiveMinutes    Type = "5min"
	FifteenMinutes Type = "15min"
	ThirtyMinutes  Type = "30min"
	OneDay         Type = "1d"
	OneWeek        Type = "1w"
	OneMonth       Type = "1m"
	OneHour        Type = "1h"
)

type Stock interface {
	Kline(stockCode string, t Type, start, end time.Time) ([]*KLine, error)
	Trend(stockCode string) ([]*Trend, error)
}
