package eastmoney

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"stock/pkg/spiders"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return errors.New("disable redirect")
	},
}

const (
	easyMoneyAPI    = "http://push2his.eastmoney.com/api"
	timeFormat      = "20060102"
	minTimeFormat   = "2006-01-02 15:04"
	kLineTimeFormat = "2006-01-02"
)

type EastMoneyProvider struct {
	httpClient *http.Client
}

type EastMoneyKLine struct {
	Data struct {
		KLines []string `json:"klines"`
	} `json:"data"`
}

type EastMoneyTrends struct {
	Data struct {
		Close  float64  `json:"preClose"`
		Trends []string `json:"trends"`
	} `json:"data"`
}

// https://blog.csdn.net/weixin_40929065/article/details/101053773
func (p *EastMoneyProvider) KLine(stockCode string, t spiders.Type, start, end time.Time) ([]*spiders.KLine, error) {
	if p.httpClient == nil {
		p.httpClient = httpClient
	}
	param := url.Values{}
	param.Set("secid", stockCode)
	param.Set("fields1", "f1,f2,f3,f4,f5")
	param.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58")
	param.Set("klt", p.getKLTFromType(t))
	param.Set("fqt", "0")
	param.Set("beg", start.Format(timeFormat))
	param.Set("end", end.Format(timeFormat))
	u := fmt.Sprintf("%s%s?%s", easyMoneyAPI, "/qt/stock/kline/get", param.Encode())
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	ed := new(EastMoneyKLine)
	err = decoder.Decode(ed)
	if err != nil {
		return nil, err
	}
	if len(ed.Data.KLines) == 0 {
		return make([]*spiders.KLine, 0), nil
	}
	kline := make([]*spiders.KLine, len(ed.Data.KLines))
	for i := range ed.Data.KLines {
		line := strings.Split(ed.Data.KLines[i], ",")
		if len(line) != 8 {
			return nil, fmt.Errorf("invalid data line [%s]", ed.Data.KLines[i])
		}
		timeLayout := kLineTimeFormat
		if t == spiders.FifteenMinutes || t == spiders.FiveMinutes || t == spiders.ThirtyMinutes || t == spiders.OneHour {
			timeLayout = minTimeFormat
		}
		klineTime, err := time.ParseInLocation(timeLayout, line[0], time.Local)
		if err != nil {
			return nil, fmt.Errorf("invalid time line [%s]", ed.Data.KLines[i])
		}
		open, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid open data line [%s]", ed.Data.KLines[i])
		}
		closePrice, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid close data line [%s]", ed.Data.KLines[i])
		}
		high, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid high data line [%s]", ed.Data.KLines[i])
		}
		low, err := strconv.ParseFloat(line[4], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid low data line [%s]", ed.Data.KLines[i])
		}
		kline[i] = &spiders.KLine{
			Open:  open,
			Close: closePrice,
			High:  high,
			Low:   low,
			Time:  klineTime,
			Type:  t,
		}
	}
	return kline, nil
}

func (p *EastMoneyProvider) Trend(stockCode string) ([]*spiders.Trend, error) {
	if p.httpClient == nil {
		p.httpClient = httpClient
	}
	param := url.Values{}
	param.Set("secid", stockCode)
	param.Set("fields1", "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13")
	param.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58")
	param.Set("iscr", "0")
	param.Set("ndays", "1")
	u := fmt.Sprintf("%s%s?%s", easyMoneyAPI, "/qt/stock/trends2/get", param.Encode())
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var ed = new(EastMoneyTrends)
	err = decoder.Decode(ed)
	if err != nil {
		return nil, err
	}
	if len(ed.Data.Trends) == 0 {
		return make([]*spiders.Trend, 0), nil
	}
	trends := make([]*spiders.Trend, len(ed.Data.Trends))
	for i := range ed.Data.Trends {
		line := strings.Split(ed.Data.Trends[i], ",")
		if len(line) != 8 {
			return nil, fmt.Errorf("invalid data line [%s]", ed.Data.Trends[i])
		}
		timeLayout := minTimeFormat
		trendTime, err := time.ParseInLocation(timeLayout, line[0], time.Local)
		if err != nil {
			return nil, fmt.Errorf("invalid time line [%s]", ed.Data.Trends[i])
		}
		price, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid open data line [%s]", ed.Data.Trends[i])
		}
		volume, err := strconv.ParseInt(line[5], 10, 0)
		if err != nil {
			return nil, fmt.Errorf("invalid open data line [%s]", ed.Data.Trends[i])
		}

		trends[i] = &spiders.Trend{
			Time:    trendTime,
			Price:   price,
			Volume:  volume,
			Incrace: (price - ed.Data.Close) / ed.Data.Close,
		}
	}
	return trends, nil
}

func (p *EastMoneyProvider) getKLTFromType(t spiders.Type) string {
	switch t {
	case spiders.FiveMinutes:
		return "5"
	case spiders.FifteenMinutes:
		return "15"
	case spiders.ThirtyMinutes:
		return "30"
	case spiders.OneHour:
		return "60"
	case spiders.OneDay:
		return "101"
	case spiders.OneWeek:
		return "102"
	case spiders.OneMonth:
		return "103"
	default:
		return "101"
	}
}
