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

type TrendData struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Stock struct {
	Name         string `json:"name"`
	Code         string `json:"code"`
	InternalCode string `json:"internal_code"`
	Type         string `json:"type"`
}

// f2: now price  f3: gains f5 成交量 f6: 成交额 f9 市盈 f12: internal_code f14 name f15 最高 f16 最低 f17今开 f18 昨收 f20 总市值 f21 流通市值 f23市净值
type MultiStock struct {
	Stock
	Price          float64 `json:"price"`
	Gains          float64 `json:"gains"`
	TrendVolume    float64 `json:"trend_volume"`    // 交易量
	TurnoverAmount float64 `json:"turnover_amount"` // 成交额
	High           float64 `json:"high"`            // 最高
	Low            float64 `json:"low"`             // 最低
	Open           float64 `json:"open"`            // 今开
	Close          float64 `json:"close"`           // 昨收
	TotalValue     float64 `json:"total_value"`     // 总市值
	Circulation    float64 `json:"circulation"`     // 流通值
	PBRatio        float64 `json:"pb_ratio"`        // 市净率
}

type StockWithDetail struct {
	Stock
	Gains          float64 `json:"gains"`           // 涨幅
	High           float64 `json:"high"`            // 最高
	Low            float64 `json:"low"`             // 最低
	Open           float64 `json:"open"`            // 今开
	Close          float64 `json:"close"`           // 昨收
	TrendVolume    float64 `json:"trend_volume"`    // 交易量
	TurnoverAmount float64 `json:"turnover_amount"` // 成交额
	QuantityRatio  float64 `json:"quantity_ratio"`  // 量比
	LimitUp        float64 `json:"limit_up"`        // 涨停
	LimitDown      float64 `json:"limit_down"`      // 跌停
	Circulation    float64 `json:"circulation"`     // 流通值
	TotalValue     float64 `json:"total_value"`     // 总市值
	PBRatio        float64 `json:"pb_ratio"`        // 市净率
	Turnover       int64   `json:"turnover"`        // 交易额
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
	Stock(code string) (*StockWithDetail, error)
	MultiStock(codes []string) ([]*MultiStock, error)
}
