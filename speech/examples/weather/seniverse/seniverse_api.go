package seniverse

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/examples/weather/model"
	"roobo.com/sailor/db/mysql"
	"roobo.com/sailor/glog"
	"roobo.com/sailor/net"
	"roobo.com/sailor/util"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrNoResult = util.NewErr("No Result")

	apiHost         string
	apiKey          string
	apiWeatherDaily string
	apiWeatherNow   string
	apiAirDaily     string
	apiAirNow       string
	apiCityId       string

	dsn    string
	dbCity string
)

func ConfMysql() error {
	s1, err1 := util.GetCfgVal("", "mysql", "dsn")
	s2, err2 := util.GetCfgVal("", "mysql", "db_city")
	dsn = s1.(string)
	dbCity = s2.(string)
	if err1 != nil || err2 != nil {
		return util.NewErrf("%v,%v", err1, err2)
	}
	return nil
}

func ConfWeatherApi() error {
	s1, err1 := util.GetCfgVal("", "weather", "api", "host")
	s2, err2 := util.GetCfgVal("", "weather", "api", "key")
	s3, err3 := util.GetCfgVal("", "weather", "api", "weather_daily")
	s4, err4 := util.GetCfgVal("", "weather", "api", "weather_now")
	s5, err5 := util.GetCfgVal("", "weather", "api", "air_daily")
	s6, err6 := util.GetCfgVal("", "weather", "api", "air_now")
	s7, err7 := util.GetCfgVal("", "weather", "api", "city_id")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil ||
		err5 != nil || err6 != nil || err7 != nil {
		return util.NewErrf("%v,%v,%v,%v,%v,%v", err1, err2, err3, err4, err5, err6, err7)
	}
	apiHost = s1.(string)
	apiKey = s2.(string)
	apiWeatherDaily = s3.(string)
	apiWeatherNow = s4.(string)
	apiAirDaily = s5.(string)
	apiAirNow = s6.(string)
	apiCityId = s7.(string)
	return nil
}

type SeniverseWeather struct {
}

func GetCityID(city string) (string, error) {
	start := time.Now()
	defer func() {
		eclipse := float64(time.Since(start).Nanoseconds()) / 1e6
		glog.Infof("[GetCityID][mysql] cost: %6.3fms, city: %s", eclipse, city)
	}()
	dbt, err := mysql.GetClient(dsn, dbCity)
	if err != nil {
		return "", err
	}
	// case 1: city in database is unique
	// e.g.: 长沙
	stmt, err := dbt.Prepare("SELECT cid FROM cities_for_seniverse_weather WHERE name = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	rows, err := stmt.Query(city)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	count := 0
	var cid string
	for rows.Next() {
		count++
		err := rows.Scan(&cid)
		if err != nil {
			return "", err
		}
	}
	if count == 1 {
		return cid, nil
	} else if count > 1 {
		// case 2: city in database where name and province is same
		// e.g. "北京", "WX4FBXXFKE4F"
		stmt, err := dbt.Prepare("SELECT cid FROM cities_for_seniverse_weather " +
			"WHERE name = ? AND attr_province = ?")
		if err != nil {
			return "", err
		}
		defer stmt.Close()
		err = stmt.QueryRow(city, city).Scan(&cid)
		if err != nil {
			return "", err
		}
		return cid, nil
	} else {
		// case 3: city in database is multiple and supply attr_province or attr_city
		// e.g.: "北京朝阳", "WX4G17JWZEK7"
		stmt, err := dbt.Prepare("SELECT cid FROM cities_for_seniverse_weather " +
			"WHERE INSTR(?, name) > 1 AND (INSTR(?, attr_province) > 0 OR " +
			"(attr_city <> '' AND INSTR(?, attr_city) > 0));")
		if err != nil {
			return "", err
		}
		defer stmt.Close()
		err = stmt.QueryRow(city, city, city).Scan(&cid)
		if err != nil {
			return "", err
		}
		return cid, nil
	}
}

func GetCityName(lat, log float64) (string, error) {
	start := time.Now()
	defer func() {
		eclipse := float64(time.Since(start).Nanoseconds()) / 1e6
		glog.Infof("[GetCityIDByLatLog] cost: %6.3fms, [lat:%f,log:%f]", eclipse, lat, log)
	}()
	url := fmt.Sprintf("%s%s?key=%s&q=%f:%f", apiHost, apiCityId, apiKey, lat, log)
	raw, err := TimeHttpGet(url)
	if err != nil {
		return "", err
	}
	LogRespBody("GetCityID for "+fmt.Sprintf("lat: %f, log: %f", lat, log)+
		" response:", raw, 500)
	var loc LocResp
	if err := json.Unmarshal(raw, &loc); err != nil {
		return "", err
	}
	if len(loc.Results) == 0 {
		return "", util.NewErrf("city name for lat: %f, log: %f not found", lat, log)
	}
	sa := strings.Split(loc.Results[0].Path, ",")
	switch len(sa) {
	default:
		if sa[0] == sa[1] {
			return sa[2] + sa[0], nil
		}
		return sa[1] + sa[0], nil
	case 0:
		return "", util.NewErrf("city name for lat: %f, log: %f is empty", lat, log)
	case 1, 2:
		return sa[0], nil
	case 3:
		if sa[0] == sa[1] {
			return sa[0], nil
		}
		return sa[1] + sa[0], nil
	}
}

func searchNow(cid string) (*NowCond, *Location, error) {
	url := fmt.Sprintf("%s%s?key=%s&location=%s&language=%s&unit=c",
		apiHost, apiWeatherNow, apiKey, cid, "zh-Hans")
	raw, err := TimeHttpGet(url)
	if err != nil {
		return nil, nil, err
	}
	LogRespBody("WeatherNow response:", raw, 500)
	var apiResp NowCondResp
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		var apiStatus ApiStatusResp
		if err := json.Unmarshal(raw, &apiStatus); err != nil {
			return nil, nil, util.NewErrf("%s:%s", apiStatus.Code, apiStatus.Status)
		}
		return nil, nil, err
	}
	if len(apiResp.Results) == 0 {
		return nil, nil, util.NewErrf("ForcastNow url error: %s", ErrNoResult)
	}
	return apiResp.Results[0].Now, apiResp.Results[0].Location, nil
}

func (sw *SeniverseWeather) SearchNow(city string) (*model.Result, error) {
	cid, err := GetCityID(city)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	cond, loc, err := searchNow(cid)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	result := model.NewResult(1).WithFocus("weather").WithWeather(cond.Text).
		WithCity(loc.Path).WithDate(time.Now().Format("2006-01-02")).
		WithTemperature(cond.Temperature).
		WithHumidity(cond.Humidity).
		WithWindDir(cond.WindDirection).WithWindLevel(cond.WindScale)
	return result, nil
}

func forcastOneDayWeather(cid string, tm time.Time) (*DailyCond, *Location, error) {
	days := util.DaysBetween(time.Now(), tm)
	url := fmt.Sprintf("%s%s?key=%s&location=%s&language=%s&unit=c&start=%d&days=%d",
		apiHost, apiWeatherDaily, apiKey, cid, "zh-Hans", days, 1)
	raw, err := TimeHttpGet(url)
	if err != nil {
		return nil, nil, err
	}
	LogRespBody("WeatherDaily response:", raw, 500)
	var apiResp DailyCondResp
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		var apiStatus ApiStatusResp
		if err := json.Unmarshal(raw, &apiStatus); err != nil {
			return nil, nil, util.NewErrf("%s:%s", apiStatus.Code, apiStatus.Status)
		}
		return nil, nil, err
	}
	if len(apiResp.Results) == 0 || len(apiResp.Results[0].Daily) == 0 {
		return nil, nil, util.NewErrf("forcastOneDayWeather url error: %s", ErrNoResult)
	}
	cond := apiResp.Results[0].Daily[0]
	loc := apiResp.Results[0].Location
	return cond, loc, nil
}

func (sw *SeniverseWeather) ForcastOneDayWeather(city string, tm time.Time) (
	*model.Result, error) {
	cid, err := GetCityID(city)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	glog.Infof("GetCityID for city: %s return id: %s", city, cid)
	cond, loc, err := forcastOneDayWeather(cid, tm)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	glog.Infof("forcastOneDayWeather for city: %s date: %v result: %+v", city, tm, cond)
	// get pm25
	aqiCond, _, err := forcastOneDayAqi(cid, tm)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	// get humidity
	nowCond, _, err := searchNow(cid)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	// make result
	wtText := cond.TextDay
	if len(cond.TextNight) > 0 && cond.TextNight != cond.TextDay {
		wtText += "转" + cond.TextNight
	}
	//wtText += "，空气质量为" + aqiCond.Quality
	result := model.NewResult(1).WithWeather(wtText).
		WithCity(loc.Path).WithDate(cond.Date).
		WithTemperature(nowCond.Temperature).WithMaxTemp(cond.High).WithMinTemp(cond.Low).
		WithWindDir(cond.WindDirection).WithWindLevel(cond.WindScale).
		WithWindDay(cond.WindDirection).WithWindDayLevel(cond.WindScale).
		WithWindNight(cond.WindDirection).WithWindNightLevel(cond.WindScale).
		WithPm25(aqiCond.Pm25).WithHumidity(nowCond.Humidity)
	return result, nil
}

func forcastDaysWeather(cid string, start, end time.Time) ([]*DailyCond, *Location, error) {
	days := util.DaysBetween(time.Now(), start)
	url := fmt.Sprintf("%s%s?key=%s&location=%s&language=%s&unit=c&start=%d&days=%d",
		apiHost, apiWeatherDaily, apiKey, cid, "zh-Hans", days, 1+util.DaysBetween(start, end))
	raw, err := TimeHttpGet(url)
	if err != nil {
		return nil, nil, err
	}
	LogRespBody("WeatherDaily response:", raw, 500)
	var apiResp DailyCondResp
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		var apiStatus ApiStatusResp
		if err := json.Unmarshal(raw, &apiStatus); err != nil {
			return nil, nil, util.NewErrf("%s:%s", apiStatus.Code, apiStatus.Status)
		}
		return nil, nil, err
	}
	if len(apiResp.Results) == 0 || len(apiResp.Results[0].Daily) == 0 {
		return nil, nil, util.NewErrf("forcastDaysWeather url error: %s", ErrNoResult)
	}
	conds := apiResp.Results[0].Daily
	loc := apiResp.Results[0].Location
	return conds, loc, nil
}

func (sw *SeniverseWeather) ForcastDaysWeather(city string, start, end time.Time) (
	model.Results, error) {
	cid, err := GetCityID(city)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	glog.Infof("GetCityID for city: %s return id: %s", city, cid)
	// days' weather
	conds, loc, err := forcastDaysWeather(cid, start, end)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	//glog.Infof("forcastDaysWeather for city[%s] start[%v] end[%v] result: %+v",
	//	city, start, end, conds)
	// get pm25
	aqiConds, _, err := forcastDaysAqi(cid, start, end)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	results := model.NewResults()
	for i, v := range conds {
		var pm25 string
		if len(aqiConds) > i {
			pm25 = aqiConds[i].Pm25
		}
		wtText := v.TextDay
		if len(v.TextNight) > 0 && v.TextNight != v.TextDay {
			wtText += "转" + v.TextNight
		}
		//wtText += "，空气质量为" + quality
		result := model.NewResult(i + 1).WithWeather(wtText).
			WithCity(loc.Path).WithDate(v.Date).
			WithMaxTemp(v.High).WithMinTemp(v.Low).
			WithWindDir(v.WindDirection).WithWindLevel(v.WindScale).
			WithWindDay(v.WindDirection).WithWindDayLevel(v.WindScale).
			WithWindNight(v.WindDirection).WithWindNightLevel(v.WindScale).
			WithPm25(pm25)
		results = results.Append(result)
	}
	return results, nil
}

func forcastOneDayAqi(cid string, tm time.Time) (*DailyAirCond, *Location, error) {
	days := util.DaysBetween(time.Now(), tm)
	url := fmt.Sprintf("%s%s?key=%s&location=%s&language=%s&start=%d&days=%d",
		apiHost, apiAirDaily, apiKey, cid, "zh-Hans", days, 1)
	raw, err := TimeHttpGet(url)
	if err != nil {
		return nil, nil, err
	}
	LogRespBody("AirDaily response:", raw, 500)
	var apiResp DailyAirResp
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		var apiStatus ApiStatusResp
		if err := json.Unmarshal(raw, &apiStatus); err != nil {
			return nil, nil, util.NewErrf("%s:%s", apiStatus.Code, apiStatus.Status)
		}
		return nil, nil, err
	}
	if len(apiResp.Results) == 0 || len(apiResp.Results[0].Daily) == 0 {
		return nil, nil, util.NewErrf("forcastOneDayAqi url error: %s", ErrNoResult)
	}
	cond := apiResp.Results[0].Daily[0]
	loc := apiResp.Results[0].Location
	return cond, loc, nil
}

func forcastDaysAqi(cid string, start, end time.Time) ([]*DailyAirCond, *Location, error) {
	days := util.DaysBetween(time.Now(), start)
	url := fmt.Sprintf("%s%s?key=%s&location=%s&language=%s&start=%d&days=%d",
		apiHost, apiAirDaily, apiKey, cid, "zh-Hans", days, 1+util.DaysBetween(start, end))
	raw, err := TimeHttpGet(url)
	if err != nil {
		return nil, nil, err
	}
	LogRespBody("WeatherDaily response:", raw, 500)
	var apiResp DailyAirResp
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		var apiStatus ApiStatusResp
		if err := json.Unmarshal(raw, &apiStatus); err != nil {
			return nil, nil, util.NewErrf("%s:%s", apiStatus.Code, apiStatus.Status)
		}
		return nil, nil, err
	}
	if len(apiResp.Results) == 0 || len(apiResp.Results[0].Daily) == 0 {
		return nil, nil, util.NewErrf("forcastDaysWeather url error: %s", ErrNoResult)
	}
	conds := apiResp.Results[0].Daily
	loc := apiResp.Results[0].Location
	return conds, loc, nil
}

func TimeHttpGet(url string) ([]byte, error) {
	start := time.Now()
	defer func() {
		eclipse := float64(time.Since(start).Nanoseconds()) / 1e6
		glog.Infof("http Get cost: %6.3fms, url: %s", eclipse, url)
	}()
	return net.Get(url, 2000)
}

func LogRespBody(prefix string, raw []byte, max int) {
	l := max
	suffix := "... "
	if len(raw) < max {
		l = len(raw)
		suffix = ""
	}
	glog.Info(prefix, string(raw[:l]), suffix, ", len=", len(raw))
}
