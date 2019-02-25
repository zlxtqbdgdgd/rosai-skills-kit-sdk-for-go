package db

import (
	"encoding/json"
	"fmt"
	"time"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/examples/weather/model"

	"roobo.com/sailor/db/redis"
	"roobo.com/sailor/util"
)

func ConfRedis() error {
	s1, err1 := util.GetCfgVal("", "redis", "addr")
	s2, err2 := util.GetCfgVal("", "redis", "passwd")
	s3, err3 := util.GetCfgVal("5", "redis", "db")
	if err1 != nil || err2 != nil || err3 != nil {
		return util.NewErrf("get redis redis conf failed, %v,%v,%v", err1, err2, err3)
	}
	host := s1.(string)
	passwd := s2.(string)
	db := s3.(string)
	return redis.InitConf(host, passwd, db)
}

func genKey(city string, dates ...time.Time) string {
	key := fmt.Sprintf("weather.%s", city)
	for _, v := range dates {
		key += "." + v.Format("2006-01-02")
	}
	return key
}

func getExpired(date time.Time) int {
	expired := 7200 // 2 hours
	switch util.DaysBetween(time.Now(), date) {
	default:
		expired *= 24
	case 0:
	case 1:
		expired *= 3
	case 2:
		expired *= 6
	case 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 14, 15:
		expired *= 12
	}
	return expired
}

func StoreOneDayResult2Redis(city string, date time.Time, result *model.Result) error {
	expired, key := getExpired(date), genKey(city, date)
	raw, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return redis.SetValueAndExpire(key, raw, expired)
}

func RestoreOneDayResult4Redis(city string, date time.Time) (*model.Result, error) {
	raw, err := redis.GetValue(genKey(city, date))
	if err != nil {
		return nil, err
	}
	bytes, ok := raw.([]byte)
	if !ok {
		return nil, util.NewErrf("Result convert error, raw: %+v", raw)
	}
	result := &model.Result{}
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

func StoreDaysResults2Redis(city string, start, end time.Time, results model.Results) error {
	expired, key := getExpired(start), genKey(city, start, end)
	raw, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return redis.SetValueAndExpire(key, raw, expired)
}

func RestoreDaysResults4Redis(city string, start, end time.Time) (model.Results, error) {
	raw, err := redis.GetValue(genKey(city, start, end))
	if err != nil {
		return nil, err
	}
	bytes, ok := raw.([]byte)
	if !ok {
		return nil, util.NewErrf("Results from redis convert error, raw: %+v", raw)
	}
	results := model.Results{}
	if err := json.Unmarshal(bytes, &results); err != nil {
		return nil, err
	}
	return results, nil
}
