package util

import (
	"encoding/json"
	"math/big"
	"time"
)

func MarshalJsonIgnoreErr(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func UnMarshalJsonIgnoreErr(data string, v interface{}) {
	json.Unmarshal([]byte(data), &v)
}

func FmtTime(t time.Time, fmt string) string {
	return t.Format(fmt)
}

const (
	TimeZoneShanghai = "Asia/Shanghai"
	AesKeyStr        = "irita-1234567890"

	DateFmtYYYYMMDD       = "2006-01-02"
	DateFmtYYYYMMDDHHmmss = "2006-01-02 15:04:05"
)
const (
	_ Unit = iota
	Day
	Hour
	Min
	Sec
)

type Unit int

var CstZone = time.FixedZone("CST", 8*3600) // 东八

func TruncateTime(t time.Time, unit Unit) time.Time {
	switch unit {
	case Day:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case Hour:
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	case Min:
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	case Sec:
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	}
	panic("not exist unit")
}

func ParseDuration(num int, unit Unit) time.Duration {
	switch unit {
	case Day:
		return time.Duration(num*24) * time.Hour
	case Hour:
		return time.Duration(num) * time.Hour
	case Min:
		return time.Duration(num) * time.Minute
	case Sec:
		return time.Duration(num) * time.Second
	}
	panic("not exist unit")
}

func RunTimer(num int, uint Unit, fn func()) {
	go func() {
		// run once right now
		fn()
		for {
			now := time.Now()
			next := now.Add(ParseDuration(num, uint))
			next = TruncateTime(next, uint)
			t := time.NewTimer(next.Sub(now))
			select {
			case <-t.C:
				fn()
			}
		}
	}()
}

func DistinctStringSlice(slice []string) []string {
	var res []string
	elementExistMap := make(map[string]bool)
	if len(slice) > 0 {
		for _, v := range slice {
			if !elementExistMap[v] {
				res = append(res, v)
				elementExistMap[v] = true
			}
		}
	}

	return res
}

func BigFloatMul(numstr string, num int64) string {
	n, ok := new(big.Rat).SetString(numstr)
	if ok {
		m := new(big.Rat).SetInt64(num)
		m.Mul(n, m)
		return m.RatString()
	}
	return ""
}
