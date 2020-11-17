package services

import (
	"errors"
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
	return nil, errors.New("not implementation")
	// data, err := s.KLine(stockCode, t, start, end)
}

func (s *StockImpl) Trend(stockCode string, day int, showBefore bool) ([]*spiders.Trend, error) {
	return s.IStock.Trend(stockCode, day, showBefore)
}

func (s *StockImpl) Search(key string) ([]*spiders.Stock, error) {
	return s.IStock.Search(key)
}
