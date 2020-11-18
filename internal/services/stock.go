package services

import (
	"stock/internal/entities"
	"stock/pkg/spiders"
	"time"
)

type StockImpl struct {
	spiders.IStock
}

func NewService(s spiders.IStock) *StockImpl {
	return &StockImpl{s}
}

func (s *StockImpl) KLine(stockCode string, t spiders.Type, start, end time.Time) (*entities.KLine, error) {
	data, err := s.IStock.KLine(stockCode, t, start, end)
	if err != nil {
		return nil, err
	}
	kline := &entities.KLine{
		Labels: make([]string, len(data)),
		KLine:  make([][]float64, len(data)),
	}
	for i, item := range data {
		if t == spiders.OneHour || t == spiders.ThirtyMinutes || t == spiders.FifteenMinutes || t == spiders.FiveMinutes {
			kline.Labels[i] = item.Time.Format("15:04")
		} else {
			kline.Labels[i] = item.Time.Format("2006-01-02")
		}

		kline.KLine[i] = []float64{
			item.Open,
			item.Close,
			item.High,
			item.Low,
		}
	}
	return kline, nil
}

func (s *StockImpl) Trend(stockCode string, day int, showBefore bool) ([]*spiders.Trend, error) {
	return s.IStock.Trend(stockCode, day, showBefore)
}

func (s *StockImpl) Search(key string) ([]*spiders.Stock, error) {
	return s.IStock.Search(key)
}

func (s *StockImpl) Stock(code string) (*spiders.StockWithDetail, error) {
	return s.IStock.Stock(code)
}
