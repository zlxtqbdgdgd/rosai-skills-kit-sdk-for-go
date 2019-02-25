package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/examples/weather/db"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/examples/weather/model"
	snv "roobo.com/rosai-skills-kit-sdk-for-go/speech/examples/weather/seniverse"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
	sp "roobo.com/rosai-skills-kit-sdk-for-go/speech/speechlet"
	"roobo.com/sailor/glog"
	"roobo.com/sailor/util"
)

const (
	IntentSearchOneDay = "SearchOneDay"
	IntentSearchDays   = "SearchDays"
	SlotDate           = "date"
	SlotDuration       = "duration"
	SlotCity           = "city"
	SlotFocus          = "focus"

	WeatherFocus  = "天气"
	TempFocus     = "温度"
	AqiFocus      = "空气"
	HumidityFocus = "湿度"
	WindFocus     = "风向"
)

type Weather struct {
}

var (
	weatherApi WeatherApi = &snv.SeniverseWeather{}

	focusMap = map[string]string{
		WeatherFocus: "weather",
		TempFocus:    "temp",
		AqiFocus:     "pm25",
	}
)

type CondDays struct {
	cond string
	days []string
}

type WeatherApi interface {
	ForcastOneDayWeather(city string, tm time.Time) (*model.Result, error)
	ForcastDaysWeather(city string, tStart, tEnd time.Time) (model.Results, error)
}

func InitConf() error {
	err1 := snv.ConfWeatherApi()
	err2 := snv.ConfMysql()
	err3 := db.ConfRedis()
	if err1 != nil || err2 != nil || err3 != nil {
		util.NewErrf("init conf error: %v,%v,%v", err1, err2, err3)
	}
	return nil
}

func (wt *Weather) OnSessionStarted(re *sp.RequestEnvelope) error {
	glog.Infof("OnSessionStarted requestId=%s", re.Request.GetRequestId())
	return nil
}

func (wt *Weather) OnSessionEnded(re *sp.RequestEnvelope) error {
	glog.Infof("OnSessionEnded requestId=%s", re.Request.GetRequestId())
	return nil
}

func (wt *Weather) OnLaunch(re *sp.RequestEnvelope) (*sp.Response, error) {
	glog.Infof("OnLaunch requestId=%s", re.Request.GetRequestId())
	return getWelcomeResponse(), nil
}

func (wt *Weather) OnIntent(re *sp.RequestEnvelope) (*sp.Response, *sp.Context, error) {
	request, ok := re.Request.(*sp.IntentRequest)
	if !ok {
		glog.Infof("re.Request type: %T, value: %+v", re.Request, re.Request)
		return nil, nil, errors.New("OnIntent assert requestEnvelope for IntentRequest failed")
	}
	intent := request.Intent
	intentName := ""
	if intent != nil {
		intentName = intent.Name
	}
	glog.Infof("OnIntent requestId=%s, intent: %s", request.RequestId, intentName)
	ctx := re.Context
	switch intentName {
	case IntentSearchOneDay:
		return handleSearchOneDayIntent(intent, ctx)
	case IntentSearchDays:
		return handleSearchDaysIntent(intent, ctx)
	case "ROSAI.HelpIntent":
		return getHelpResponse()
	default:
		tip := fmt.Sprintf("Intent(%s) is unsupported. Please try something else.", intentName)
		return sp.NewAskResponse(tip), nil, nil
	}
}

func handleSearchOneDayIntent(intent *slu.Intent, inCtx *sp.Context) (
	resp *sp.Response, ctx *sp.Context, err error) {
	defer func() {
		ctx = inCtx
	}()
	var city, date, focus string
	if city = intent.GetSlot(SlotCity).GetStringValue(); city == "" {
		if city = inCtx.GetStringValue(SlotCity); city == "" {
			glog.Infof("get city[%s] from context", city)
			if city = getCityFromSysInfo(inCtx); city == "" {
				return sp.NewAskResponse("你要查询哪个城市的天气"), nil, nil
			} else {
				glog.Infof("get city[%s] from context system info", city)
				//inCtx.SetStringValue(SlotCity, city)
			}
		}
	}
	if date = intent.GetSlot(SlotDate).GetStringValue(); date == "" {
		if date = inCtx.GetStringValue(SlotDate); date == "" {
			date = time.Now().Format("2006-01-02")
			//return sp.NewAskResponse("你要查询哪一天的天气"), nil, nil
		}
	}
	if focus = intent.GetSlot(SlotFocus).GetStringValue(); focus == "" {
		if focus = inCtx.GetStringValue(SlotFocus); focus == "" {
			glog.Warning("focus is null")
			return nil, nil, util.NewErr("focus for one day is null")
		}
	}
	glog.Infof("SearchOneDay slots city: %s, date: %s, focus: %s", city, date, focus)
	return getFinalOneDayResponse(city, date, focus)
}

func handleSearchDaysIntent(intent *slu.Intent, inCtx *sp.Context) (
	resp *sp.Response, ctx *sp.Context, err error) {
	defer func() {
		ctx = inCtx
	}()
	var city, duration, focus string
	if city = intent.GetSlot(SlotCity).GetStringValue(); city == "" {
		if city = inCtx.GetStringValue(SlotCity); city == "" {
			glog.Infof("get city[%s] from context", city)
			if city = getCityFromSysInfo(inCtx); city == "" {
				return sp.NewAskResponse("你要查询哪个城市的天气"), nil, nil
			} else {
				glog.Infof("get city[%s] from context system info", city)
				//inCtx.SetStringValue(SlotCity, city)
			}
		}
	}
	if duration = intent.GetSlot(SlotDuration).GetStringValue(); duration == "" {
		if duration = inCtx.GetStringValue(SlotDuration); duration == "" {
			return sp.NewAskResponse("你要查询哪段时间的天气"), nil, nil
		}
	}
	if focus = intent.GetSlot(SlotFocus).GetStringValue(); focus == "" {
		if focus = inCtx.GetStringValue(SlotFocus); focus == "" {
			glog.Warning("focus is null")
			return nil, nil, util.NewErr("focus for days is null")
		}
	}
	glog.Infof("SearchDays slots city: %s, duration: %s, focus: %s", city, duration, focus)
	return getFinalDaysResponse(city, duration, focus)
}

func getFinalOneDayResponse(city, date, focus string) (
	*sp.Response, *sp.Context, error) {
	tDate, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return nil, nil, err
	}
	var (
		result *model.Result
		text   string
	)
	result, err = db.RestoreOneDayResult4Redis(city, tDate)
	if err != nil {
		glog.Infof("restoreOneDayResult4Redis city[%s] date[%s] result[%+v] error: %s",
			city, tDate.Format("2006-01-02"), result, err)
		result, err = weatherApi.ForcastOneDayWeather(city, tDate)
		if err != nil {
			return nil, nil, err
		}
		if err := db.StoreOneDayResult2Redis(city, tDate, result); err != nil {
			glog.Warningf("storeOneDayResult2Redis city[%s] date[%v] error: %s",
				city, tDate.Format("2006-01-02"), err)
		}
	} else {
		glog.Infof("restoreOneDayResult4Redis city[%s] date[%s] successed",
			city, tDate.Format("2006-01-02"))
	}
	switch focus {
	default:
		return nil, nil, util.NewErrf("unrecognize focus[%s]", focus)
	case WeatherFocus:
		text, err = makeWeatherSpeechText(result)
	case AqiFocus:
		text, err = makeAqiSpeechText(result)
	case TempFocus:
		text, err = makeTempSpeechText(result)
	case HumidityFocus:
		text, err = makeHumiditySpeechText(result)
	case WindFocus:
		text, err = makeWindSpeechText(result)
	}
	if err != nil {
		return nil, nil, err
	}
	result.SetFocus(focusMap[focus])
	resp := sp.NewTellResponse(text).WithResults(
		sp.NewResult().WithOutputPlainTextSpeech(text).WithData(result))
	return resp, nil, nil
}

func getFinalDaysResponse(city, duration, focus string) (
	*sp.Response, *sp.Context, error) {
	sa := strings.Split(duration, "/")
	if len(sa) != 2 {
		return nil, nil, util.NewErrf("parse duration[%s] error", duration)
	}
	start, err1 := time.ParseInLocation("2006-01-02", sa[0], time.Local)
	end, err2 := time.ParseInLocation("2006-01-02", sa[1], time.Local)
	if err1 != nil || err2 != nil {
		return nil, nil, util.NewErrf("%v %v", err1, err2)
	}
	var (
		results model.Results
		text    string
		err     error
	)
	results, err = db.RestoreDaysResults4Redis(city, start, end)
	if err != nil {
		glog.Infof("restoreDaysResults4Redis city[%s] start[%s] end[%s] error: %s",
			city, start.Format("2006-01-02"), end.Format("2006-01-02"), err)
		results, err = weatherApi.ForcastDaysWeather(city, start, end)
		if err != nil {
			return nil, nil, err
		}
		if err := db.StoreDaysResults2Redis(city, start, end, results); err != nil {
			glog.Warningf("storeOneDayResult2Redis city[%s] start[%s] end[%s]"+
				" results[%+v] error: %s", city, start.Format("2006-01-02"),
				end.Format("2006-01-02"), results, err)
		}
	} else {
		glog.Infof("restoreDaysResult4Redis city[%s] date[%s-%s] result[%d] successed",
			city, start.Format("2006-01-02"), end.Format("2006-01-02"), len(results))
	}
	switch focus {
	default:
		return nil, nil, util.NewErrf("unrecognize focus[%s]", focus)
	case WeatherFocus:
		for _, v := range results {
			v.SetFocus("weather")
		}
		text, err = makeWeatherSpeechText(results...)
	case AqiFocus:
		for _, v := range results {
			v.SetFocus("pm25")
		}
		text, err = makeAqiSpeechText(results...)
	case TempFocus:
		for _, v := range results {
			v.SetFocus("pm25")
		}
		text, err = makeTempSpeechText(results...)
	case HumidityFocus:
		text, err = makeHumiditySpeechText(results...)
	case WindFocus:
		text, err = makeWindSpeechText(results...)
	}
	if err != nil {
		return nil, nil, err
	}
	resp := sp.NewTellResponse(text).
		WithResults(sp.NewResult().WithOutputPlainTextSpeech(text).WithData(results))
	return resp, nil, nil
}

func getHelpResponse() (*sp.Response, *sp.Context, error) {
	text := "您可以对我说“北京今天天气”"
	return sp.NewAskResponse(text), nil, nil
}

func getWelcomeResponse() *sp.Response {
	text := "欢迎使用roobo天气服务，您可以对我说“北京今天天气”"
	return sp.NewAskResponse(text)
}

func makeAqiSpeechText(results ...*model.Result) (string, error) {
	if len(results) == 0 {
		return "", snv.ErrNoResult
	}
	var (
		err                    error
		city, date, pm25, text string
	)
	city = parseCityHint(results[0].City)
	if len(results) == 1 {
		r := results[0]
		date = util.GetOralDate("2006-01-02", r.Date)
		pm25 = fmt.Sprintf("PM2.5是%s", r.Pm25)
		iPm25, err := strconv.Atoi(r.Pm25)
		if err != nil {
			return "", util.NewErrf("pm25[%s] type string, %s", r.Pm25, err)
		}
		text, err = getQualityText(iPm25)
		if err != nil {
			return "", err
		}
	} else {
		date = fmt.Sprintf("未来%d天", len(results))
		pm25, err = makeDaysAqiText(results)
		if err != nil {
			return "", err
		}
	}
	speechText := city + date + "，" + pm25
	if text != "" {
		speechText += "，" + text
	}
	return speechText, nil
}

func makeWindSpeechText(results ...*model.Result) (string, error) {
	if len(results) == 0 {
		return "", snv.ErrNoResult
	}
	var (
		err              error
		city, date, wind string
	)
	city = parseCityHint(results[0].City)
	//
	if len(results) == 1 {
		r := results[0]
		date = util.GetOralDate("2006-01-02", r.Date)
		wind = r.WindDir + "风" + r.WindLevel + "级"
	} else {
		date = fmt.Sprintf("未来%d天", len(results))
		wind, err = makeDaysWindText(results)
		if err != nil {
			return "", err
		}
	}
	return (city + date + "，" + wind), nil
}

func makeHumiditySpeechText(results ...*model.Result) (string, error) {
	if len(results) == 0 {
		return "", snv.ErrNoResult
	}
	city := parseCityHint(results[0].City)
	//
	if len(results) == 1 {
		r := results[0]
		temp := fmt.Sprintf("当前湿度是%s", r.Humidity)
		return (city + "，" + temp), nil
	} else {
		text := "目前不支持未来湿度查询"
		return text, nil
	}
}

func makeTempSpeechText(results ...*model.Result) (string, error) {
	if len(results) == 0 {
		return "", snv.ErrNoResult
	}
	city := parseCityHint(results[0].City)
	var date, temp string
	//
	if len(results) == 1 {
		r := results[0]
		date = util.GetOralDate("2006-01-02", r.Date)
		temp = fmt.Sprintf("气温%s度到%s度", r.MaxTemp, r.MinTemp)
	} else {
		date = fmt.Sprintf("未来%d天", len(results))
		var err error
		temp, err = makeDaysTempText(results)
		if err != nil {
			return "", err
		}
	}
	speechText := city + date + "，" + temp
	return speechText, nil
}

func makeWeatherSpeechText(results ...*model.Result) (string, error) {
	if len(results) == 0 {
		return "", snv.ErrNoResult
	}
	city := parseCityHint(results[0].City)
	var text, date, temp, wind, caution string
	//
	if len(results) == 1 {
		r := results[0]
		date = util.GetOralDate("2006-01-02", r.Date)
		text = r.Weather
		temp = fmt.Sprintf("气温%s度到%s度", r.MinTemp, r.MaxTemp)
		wind = r.WindDir + "风" + r.WindLevel + "级"
	} else {
		date = fmt.Sprintf("未来%d天", len(results))
		var err error
		temp, err = makeDaysTempText(results)
		if err != nil {
			return "", err
		}
		caution, err = makeAbnormalWeatherText(results)
		if err != nil {
			return "", err
		}
	}
	speechText := city + date + text + "，" + temp
	if wind != "" {
		speechText += "，" + wind
	}
	if caution != "" {
		speechText += "，" + caution
	}
	return speechText, nil
}

func parseCityHint(path string) string {
	var city string
	pa := strings.Split(path, ",")
	if len(pa) > 2 {
		if pa[0] == pa[1] {
			city = pa[0]
		} else {
			city = pa[1] + pa[0]
		}
	} else {
		city = path
	}
	return city
}

// output e.g.: 5月16日，5月18日有降雨
func makeAbnormalWeatherText(results []*model.Result) (string, error) {
	var days, conds []string
	for _, v := range results {
		if strings.Contains(v.Weather, "雨") ||
			strings.Contains(v.Weather, "雪") ||
			strings.Contains(v.Weather, "尘") ||
			strings.Contains(v.Weather, "沙") ||
			strings.Contains(v.Weather, "霾") ||
			strings.Contains(v.Weather, "大风") ||
			strings.Contains(v.Weather, "飓风") ||
			strings.Contains(v.Weather, "热带风暴") ||
			strings.Contains(v.Weather, "龙卷风") {
			s, err := util.FormatTime(v.Date, "2006-01-02", "1月2日")
			if err != nil {
				return "", err
			}
			days = append(days, s)
			conds = append(conds, v.Weather)
		}
	}
	if len(conds) == 0 {
		return "无异常天气", nil
	} else if len(conds) == 1 {
		return fmt.Sprintf("其中%s%s", days[0], conds[0]), nil
	} else {
		// squeeze days with the same condition
		var condDays []CondDays
		j := 0
		for i, v := range conds {
			if i == 0 {
				condDays = append(condDays, CondDays{v, []string{days[0]}})
				continue
			}
			if v != conds[j] {
				if i != j+1 {
					j++
					conds[j] = v
				}
				condDays = append(condDays, CondDays{v, []string{days[i]}})
			} else {
				condDays[j].days = append(condDays[j].days, days[i])
			}
		}
		var items []string
		for _, v := range condDays {
			days, err := util.SqueezeDaysText("1月2日", v.days)
			if err != nil {
				return "", err
			}
			items = append(items, fmt.Sprintf("%s%s", strings.Join(days, "，"), v.cond))
		}
		prefix := "另外"
		if len(days) < len(results) {
			prefix = "其中"
		}
		return prefix + strings.Join(items, "，"), nil
	}
}

// output e.g.: 最大风力5级，出现在5月22日，最低风力3级，出现在5月22日至5月24日
func makeDaysWindText(results []*model.Result) (string, error) {
	if len(results) == 0 {
		return "", errors.New("makeDaysWindText input results len is 0")
	}
	lowest, highest, loDays, hiDays, err := getLoHiDays(
		func(v *model.Result) (int, error) {
			return strconv.Atoi(v.WindLevel)
		},
		func(v *model.Result) (int, error) {
			return strconv.Atoi(v.WindLevel)
		},
		results...)
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf("最大风力%d级，出现在%s，最低风力%d级，出现在%s",
		highest, strings.Join(hiDays, "，"), lowest, strings.Join(loDays, "，"))
	return text, nil
}

// output e.g.: 最高气温31度，出现在5月14日，5月17日，最低气温23度，出现在5月16日
func makeDaysTempText(results []*model.Result) (string, error) {
	if len(results) == 0 {
		return "", errors.New("makeTempText input results len is 0")
	}
	lowest, highest, loDays, hiDays, err := getLoHiDays(
		func(v *model.Result) (int, error) {
			return strconv.Atoi(v.MinTemp)
		},
		func(v *model.Result) (int, error) {
			return strconv.Atoi(v.MaxTemp)
		},
		results...)
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf("最高气温%d度，出现在%s，最低气温%d度，出现在%s",
		highest, strings.Join(hiDays, "，"), lowest, strings.Join(loDays, "，"))
	return text, nil
}

// output e.g.: 最高气温31度，出现在5月14日，5月17日，最低气温23度，出现在5月16日
func makeDaysAqiText(results []*model.Result) (string, error) {
	if len(results) == 0 {
		return "", errors.New("makeDaysAqiText input params len is 0")
	}
	lowest, highest, loDays, hiDays, err := getLoHiDays(
		func(v *model.Result) (int, error) {
			return strconv.Atoi(v.Pm25)
		},
		func(v *model.Result) (int, error) {
			return strconv.Atoi(v.Pm25)
		},
		results...)
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf("最高Pm2.5为%d，出现在%s，最低Pm2.5为%d，出现在%s",
		highest, strings.Join(hiDays, "，"), lowest, strings.Join(loDays, "，"))
	return text, nil
}

func getLoHiDays(low, high func(v *model.Result) (int, error), results ...*model.Result) (
	lowest, highest int, loDays, hiDays []string, err error) {
	highest, lowest = -10000, 10000
	for _, v := range results {
		var (
			hi, lo int
			text   string
		)
		hi, err = high(v)
		if err != nil {
			return
		}
		if hi > highest {
			highest = hi
			hiDays = make([]string, 1)
			hiDays[0], err = util.FormatTime(v.Date, "2006-01-02", "1月2日")
			if err != nil {
				return
			}
		} else if hi == highest {
			text, err = util.FormatTime(v.Date, "2006-01-02", "1月2日")
			if err != nil {
				return
			}
			hiDays = append(hiDays, text)
			hiDays, err = util.SqueezeDaysText("1月2日", hiDays)
			if err != nil {
				return
			}
		}
		lo, err = low(v)
		if err != nil {
			return
		}
		if lo < lowest {
			lowest = lo
			loDays = make([]string, 1)
			loDays[0], err = util.FormatTime(v.Date, "2006-01-02", "1月2日")
			if err != nil {
				return
			}
		} else if lo == lowest {
			text, err = util.FormatTime(v.Date, "2006-01-02", "1月2日")
			if err != nil {
				return
			}
			loDays = append(loDays, text)
			loDays, err = util.SqueezeDaysText("1月2日", loDays)
			if err != nil {
				return
			}
		}
	}
	return
}

func getQualityText(pm25 int) (string, error) {
	var quality string
	switch {
	case pm25 < 50:
		quality = "优"
	case pm25 < 100:
		quality = "良"
	case pm25 < 150:
		quality = "轻度污染"
	case pm25 < 200:
		quality = "中度污染"
	case pm25 < 300:
		quality = "重度污染"
	case pm25 >= 300:
		quality = "严重污染"
	}
	v, err := util.GetCfgVal([]string{}, "weather", "aqiText", quality)
	if err != nil {
		return "", err
	}
	ss, ok := v.([]interface{})
	if !ok {
		return "", util.NewErrf("value[%+v] type assertion failed", v)
	}
	if len(ss) == 0 {
		return "", util.NewErrf("weather aqiText %s no hint", quality)
	}
	texts := util.Conv2StrSlice(ss)
	rand.Seed(int64(time.Now().Second()))
	return ("空气质量为" + quality + "，" + texts[rand.Intn(len(texts))]), nil
}

func getCityFromSysInfo(ctx *sp.Context) string {
	rawLoc, err := ctx.GetSysParameter("location").GetMapValue()
	if err != nil {
		return ""
	}
	if addr, ok := rawLoc["address"].(*model.Address); ok {
		province, city := addr.Province, addr.City
		if province != "" && city != "" {
			province = strings.TrimRight(province, "自治区省市")
			city = strings.TrimRight(city, "市县区自治州")
			if province == city {
				return city
			}
			return province + city
		}
	}
	lat, ok := rawLoc["latitude"].(float64)
	if !ok {
		return ""
	}
	log, ok := rawLoc["longitude"].(float64)
	if !ok {
		return ""
	}
	city, err := snv.GetCityName(lat, log)
	if err != nil {
		return ""
	}
	return city
}
