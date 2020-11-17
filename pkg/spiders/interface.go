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
	Time    time.Time `json:"time" time_format:"15:04"`
	Price   float64   `json:"price"`
	Volume  int64     `json:"volume"`
	Incrace float64   `json:"incrace"`
}

type Stock struct {
	Name         string `json:"name"`
	Code         string `json:"code"`
	InternalCode string `json:"internal_code"`
	Type         string `json:"type"`
}

type Type string

const (
	FiveMinutes    Type = "5min"
	FifteenMinutes Type = "15min"
	ThirtyMinutes  Type = "30min"
	OneHour        Type = "1h"
	OneDay         Type = "1d"
	OneWeek        Type = "1w"
	OneMonth       Type = "1m"
)

type IStock interface {
	KLine(stockCode string, t Type, start, end time.Time) ([]*KLine, error)
	Trend(stockCode string, day int, showBefore bool) ([]*Trend, error)
	Search(key string) ([]*Stock, error)
}
