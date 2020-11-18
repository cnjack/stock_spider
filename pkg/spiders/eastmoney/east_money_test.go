package eastmoney_test

import (
	"stock/pkg/spiders"
	"stock/pkg/spiders/eastmoney"
	"testing"
	"time"

	jsontime "github.com/liamylian/jsontime/v2/v2"
	"github.com/stretchr/testify/assert"
)

var json = jsontime.ConfigWithCustomTimeFormat

func init() {
	jsontime.SetDefaultTimeFormat("2006-01-02 15:04", time.Local)
}

func TestEastMoneyProvider_KLine(t *testing.T) {
	spider := &eastmoney.EastMoneyProvider{}
	end := time.Now()
	data, err := spider.KLine("90.BK0729", spiders.OneHour, end.AddDate(0, 0, -10), end)
	if assert.NoError(t, err) {
		out, _ := json.Marshal(data)
		t.Log(string(out))
	}
}

func TestEastMoneyProvider_Trend(t *testing.T) {
	spider := &eastmoney.EastMoneyProvider{}
	data, err := spider.Trend("1.600350", 2, true)
	if assert.NoError(t, err) {
		out, _ := json.Marshal(data)
		t.Log(string(out))
	}
}

func TestEastMoneyProvider_Search(t *testing.T) {
	spider := &eastmoney.EastMoneyProvider{}
	data, err := spider.Search("600350")
	if assert.NoError(t, err) {
		out, _ := json.Marshal(data)
		t.Log(string(out))
	}
}

func TestEastMoneyProvider_Stock(t *testing.T) {
	spider := &eastmoney.EastMoneyProvider{}
	data, err := spider.Stock("0.300059")
	if assert.NoError(t, err) {
		out, _ := json.Marshal(data)
		t.Log(string(out))
	}
}
