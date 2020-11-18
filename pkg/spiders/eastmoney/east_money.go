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
	easyMoneyAPI       = "http://push2.eastmoney.com/api/"
	easyMoneySearchAPI = "http://searchapi.eastmoney.com/api/"
	timeFormat         = "20060102"
	minTimeFormat      = "2006-01-02 15:04"
	kLineTimeFormat    = "2006-01-02"
)

type EastMoneyProvider struct {
	httpClient *http.Client
}

var _ spiders.IStock = new(EastMoneyProvider)

type EastMoneyKLine struct {
	Data struct {
		KLines []string `json:"klines"`
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
	u := fmt.Sprintf("%s%s?%s", easyMoneyAPI, "qt/stock/kline/get", param.Encode())
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

type EastMoneyTrends struct {
	Data struct {
		Close  float64  `json:"preClose"`
		Trends []string `json:"trends"`
	} `json:"data"`
}

func (p *EastMoneyProvider) Trend(stockCode string, day int, showBefore bool) ([]*spiders.Trend, error) {
	if p.httpClient == nil {
		p.httpClient = httpClient
	}
	param := url.Values{}
	param.Set("secid", stockCode)
	param.Set("fields1", "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13")
	param.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58")
	iscr := "0"
	if showBefore {
		iscr = "1"
	}
	param.Set("iscr", iscr)
	param.Set("ndays", strconv.Itoa(day))
	u := fmt.Sprintf("%s%s?%s", easyMoneyAPI, "qt/stock/trends2/get", param.Encode())
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

type EastMoneyStockSearch struct {
	Data []struct {
		Name             string `json:"Name"`
		Code             string `json:"Code"`
		MktNum           string `json:"MktNum"`
		SecurityTypeName string `json:"SecurityTypeName"`
	} `json:"data"`
}

func (p *EastMoneyProvider) Search(key string) ([]*spiders.Stock, error) {
	if p.httpClient == nil {
		p.httpClient = httpClient
	}
	param := url.Values{}
	param.Set("and14", fmt.Sprintf("MultiMatch/Name,Code,PinYin/%s/true", key))
	param.Set("type", "14")
	param.Set("appid", "el1902262")
	param.Set("token", "CCSDCZSDCXYMYZYYSYYXSMDDSMDHHDJT")
	param.Set("returnfields14", "Name,Code,MktNum,SecurityTypeName")
	param.Set("pageIndex14", "1")
	param.Set("pageSize14", "20")
	u := fmt.Sprintf("%s%s?%s", easyMoneySearchAPI, "Info/Search", param.Encode())
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
	var ed = new(EastMoneyStockSearch)
	err = decoder.Decode(ed)
	if err != nil {
		return nil, err
	}
	if len(ed.Data) == 0 {
		return make([]*spiders.Stock, 0), nil
	}
	stocks := make([]*spiders.Stock, len(ed.Data))
	for i := range ed.Data {
		stocks[i] = &spiders.Stock{
			Name:         ed.Data[i].Name,
			Code:         ed.Data[i].Code,
			InternalCode: fmt.Sprintf("%s.%s", ed.Data[i].MktNum, ed.Data[i].Code),
			Type:         ed.Data[i].SecurityTypeName,
		}
	}
	return stocks, nil
}

type EastMoneyStock struct {
	Data struct {
		F43  int64   `json:"F43"`
		F44  int64   `json:"F44"`
		F45  int64   `json:"F45"`
		F46  int64   `json:"F46"`
		F47  int64   `json:"F47"`
		F48  float64 `json:"F48"`
		F50  int64   `json:"F50"`
		F51  int64   `json:"F51"`
		F52  int64   `json:"F52"`
		F57  string  `json:"F57"`
		F58  string  `json:"F58"`
		F60  int64   `json:"F60"`
		F107 int     `json:"F107"`
		F117 float64 `json:"F117"`
		F116 float64 `json:"F116"`
		F128 string  `json:"F128"`
		F167 int64   `json:"F167"`
		F168 int64   `json:"F168"`
	} `json:"Data"`
}

func intToFloat64(in int64) float64 {
	return float64(in) / 100
}

// f43 涨幅  f44 最高 f45 最低 f46 今开 f60 昨收 f47 成交量 f48 成交额 f50 量比 f51 涨停 f52 跌停 f57 code f58 name:
// f117 流通值 f116 总市值 f167 市净率 f168 换手  f128 板块 f107 start

func (s *EastMoneyStock) ToStockWithDetail() *spiders.StockWithDetail {
	return &spiders.StockWithDetail{
		Stock: spiders.Stock{
			Name:         s.Data.F58,
			Code:         s.Data.F57,
			InternalCode: strconv.Itoa(s.Data.F107) + s.Data.F57,
			Type:         s.Data.F128,
		},
		Gains:          intToFloat64(s.Data.F43),
		High:           intToFloat64(s.Data.F44),
		Low:            intToFloat64(s.Data.F45),
		Open:           intToFloat64(s.Data.F46),
		Close:          intToFloat64(s.Data.F60),
		TrendVolume:    intToFloat64(s.Data.F47),
		TurnoverAmount: s.Data.F48,
		QuantityRatio:  intToFloat64(s.Data.F50),
		LimitUp:        intToFloat64(s.Data.F51),
		LimitDown:      intToFloat64(s.Data.F52),
		Circulation:    s.Data.F117,
		TotalValue:     s.Data.F116,
		PBRatio:        intToFloat64(s.Data.F167),
		Turnover:       s.Data.F168,
	}
}

func (p *EastMoneyProvider) Stock(code string) (*spiders.StockWithDetail, error) {
	if p.httpClient == nil {
		p.httpClient = httpClient
	}
	param := url.Values{}
	param.Set("secid", code)
	param.Set("fields", "f43,f44,f45,f46,f47,f48,f50,f51,f52,f57,f58,f60,f107,f110,f116,f117,f128,f167,f168")
	u := fmt.Sprintf("%s%s?%s", easyMoneyAPI, "qt/stock/get", param.Encode())
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
	var s = new(EastMoneyStock)
	err = decoder.Decode(s)
	if err != nil {
		return nil, err
	}

	return s.ToStockWithDetail(), nil
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
