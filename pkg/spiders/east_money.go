package spiders

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

var _ IStock = new(EastMoneyProvider)

type EastMoneyKLine struct {
	Data struct {
		KLines []string `json:"klines"`
	} `json:"data"`
}

func (p *EastMoneyProvider) getKLTFromType(t Type) string {
	switch t {
	case FiveMinutes:
		return "5"
	case FifteenMinutes:
		return "15"
	case ThirtyMinutes:
		return "30"
	case OneHour:
		return "60"
	case OneDay:
		return "101"
	case OneWeek:
		return "102"
	case OneMonth:
		return "103"
	default:
		return "101"
	}
}

// https://blog.csdn.net/weixin_40929065/article/details/101053773
func (p *EastMoneyProvider) KLine(stockCode string, t Type, start, end time.Time) ([]*KLine, error) {
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
		return make([]*KLine, 0), nil
	}
	kline := make([]*KLine, len(ed.Data.KLines))
	for i := range ed.Data.KLines {
		line := strings.Split(ed.Data.KLines[i], ",")
		if len(line) != 8 {
			return nil, fmt.Errorf("invalid data line [%s]", ed.Data.KLines[i])
		}
		timeLayout := kLineTimeFormat
		if t == FifteenMinutes || t == FiveMinutes || t == ThirtyMinutes || t == OneHour {
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
		kline[i] = &KLine{
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

func (p *EastMoneyProvider) Trend(stockCode string, day int, showBefore bool) ([]*Trend, error) {
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
		return make([]*Trend, 0), nil
	}
	trends := make([]*Trend, len(ed.Data.Trends))
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

		trends[i] = &Trend{
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

func (p *EastMoneyProvider) Search(key string) ([]*Stock, error) {
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
		return make([]*Stock, 0), nil
	}
	stocks := make([]*Stock, len(ed.Data))
	for i := range ed.Data {
		stocks[i] = &Stock{
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

func (s *EastMoneyStock) ToStockWithDetail() *StockWithDetail {
	return &StockWithDetail{
		Stock: Stock{
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

func (p *EastMoneyProvider) Stock(code string) (*StockWithDetail, error) {
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

type EastMoneyMultiStockItem struct {
	F2  int64   `json:"F2"`
	F3  int64   `json:"F3"`
	F5  int64   `json:"F5"`
	F6  float64 `json:"F6"`
	F9  int64   `json:"F9"`
	F12 string  `json:"F12"`
	F13 int     `json:"F13"`
	F14 string  `json:"F14"`
	F15 int64   `json:"F15"`
	F16 int64   `json:"F16"`
	F17 int64   `json:"F17"`
	F18 int64   `json:"F18"`
	F20 int64   `json:"F20"`
	F21 int64   `json:"F21"`
	F23 int64   `json:"F23"`
}

// f2: now price  f3: gains f5 成交量 f6: 成交额 f9 市盈 f12: internal_code f13 market numb f14 name f15 最高 f16 最低 f17今开 f18 昨收 f20 总市值 f21 流通市值 f23 市净值
// https://blog.csdn.net/qq_38704184/article/details/101292802

func (ms *EastMoneyMultiStockItem) ToMultiStock() *MultiStock {
	return &MultiStock{
		Stock: Stock{
			Name:         ms.F14,
			Code:         ms.F12,
			InternalCode: strconv.Itoa(ms.F13) + "." + ms.F12,
		},
		Price:          intToFloat64(ms.F2),
		Gains:          intToFloat64(ms.F3),
		TrendVolume:    intToFloat64(ms.F5),
		TurnoverAmount: ms.F6,
		High:           intToFloat64(ms.F15),
		Low:            intToFloat64(ms.F16),
		Open:           intToFloat64(ms.F17),
		Close:          intToFloat64(ms.F18),
		TotalValue:     intToFloat64(ms.F20),
		Circulation:    intToFloat64(ms.F21),
		PBRatio:        intToFloat64(ms.F23),
	}
}

type EastMoneyMultiStock struct {
	Data struct {
		Diff map[string]*EastMoneyMultiStockItem `json:"diff"`
	} `json:"data"`
}

func (p *EastMoneyProvider) MultiStock(codes []string) ([]*MultiStock, error) {
	if p.httpClient == nil {
		p.httpClient = httpClient
	}
	param := url.Values{}
	param.Set("pi", "0")
	param.Set("fs", fmt.Sprintf("i:%s", strings.Join(codes, ",i:")))
	param.Set("fields", "f2,f3,f5,f6,f9,f12,f13,f14,f15,f16,f17,f18,f19,f20,f21,f22,f23")
	u := fmt.Sprintf("%s%s?%s", easyMoneyAPI, "qt/clist/get", param.Encode())
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
	var s = new(EastMoneyMultiStock)
	err = decoder.Decode(s)
	if err != nil {
		return nil, err
	}
	ms := make([]*MultiStock, 0)
	for key := range s.Data.Diff {
		ms = append(ms, s.Data.Diff[key].ToMultiStock())
	}
	return ms, nil
}
