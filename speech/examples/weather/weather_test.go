package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	snet "roobo.com/sailor/net"
	"roobo.com/sailor/util"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/examples/weather/model"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
	sp "roobo.com/rosai-skills-kit-sdk-for-go/speech/speechlet"
)

const (
	router   = "/weather"
	userId   = "rosai1.ask.account.001"
	appId    = "rosai1.ask.skill.developer.001"
	skillId  = "rosai1.ask.skill.weather.v1.0"
	deviceId = "rosai1.device.001"

	launchResponse = `{
  "version": "2.0",
  "status": {
    "code": 0
  },
  "results": [
    {
      "outputSpeech": {
        "items": [
          {
            "type": "PlainText",
            "source": "欢迎使用roobo天气服务，您可以对我说“北京今天天气”"
          }
        ]
      }
    }
  ]
}`

	respElicitToCity = `{
  "version": "2.0",
  "status": {
    "code": 0
  },
  "context": {
    "parameters": {
      "%s": {
        "orgin": null,
        "normType": "String",
        "norm": "%s"
      },
      "focus": {
        "orgin": null,
        "normType": "String",
        "norm": "天气"
      }
    }
  },
  "results": [
    {
      "outputSpeech": {
        "items": [
          {
            "type": "PlainText",
            "source": "你要查询哪个城市的天气"
          }
        ]
      }
    }
  ]
}`

	respDialogCompleted = `{
  "version": "2.0",
  "status": {
    "code": 0
  },
  "context": {
    "parameters": {
      "city": {
        "orgin": null,
        "normType": "String",
        "norm": "北京"
      },
      "%s": {
        "orgin": null,
        "normType": "String",
        "norm": "%s"
      },
      "focus": {
        "orgin": null,
        "normType": "String",
        "norm": "天气"
      }
    }
  },
  "results": [
    {
      "formatType": "SpeechFormat",
      "hint": "北京今天，气温12度到27度，东南风3至2级，多云",
      "data": {
        "index": 1,
        "pm25": "76",
        "city": "北京",
        "focus": "weather",
        "weather": "多云",
        "temperature": "22",
        "minTemp": "12",
        "maxTemp": "27",
        "date": "2018-05-09",
        "humidity": "31",
        "windDir": "西南风",
        "windLevel": "2",
        "windDay": "东南风",
        "windDayLevel": "3",
        "windNight": "东南风",
        "windNightLevel": "2",
        "alter": ""
      }
    }
  ]
}`
)

var (
	addr    net.Addr
	starter sync.Once

	// request info
	reqId, ts = "12345", "2018-04-06T15:30:02+08:00"
	sysInfo   = sp.NewCtxSystem().
			WithUser(sp.NewUser(userId, appId)).
			WithSkill(sp.NewSkill(skillId)).
			WithDevice(sp.NewDevice(deviceId))
	ctx = sp.NewContext().WithSystem(sysInfo)
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func setupMockServer(t *testing.T) {
	dm, err := getDialogModel()
	if err != nil {
		t.Fatal(err)
	}
	rh := sp.RequestHandler{
		AppId:       "rosai1.ask.skill.planmytrip.12345",
		Speechlet:   &Weather{},
		DialogModel: dm,
	}
	http.Handle(router, &rh)
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen - %s", err.Error())
	}
	go func() {
		err = http.Serve(ln, nil)
		if err != nil {
			t.Fatalf("failed to start HTTP server - %s", err.Error())
		}
	}()
	addr = ln.Addr()
}

func runResponseTest(t *testing.T, resp string, f func() *sp.RequestEnvelope) {
	reqEn := f()
	reqBytes, err := json.MarshalIndent(reqEn, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	url := "http://" + addr.String() + router
	respBytes, err := snet.Post(url, reqBytes, 3000)
	if err != nil {
		t.Fatal(err)
	}
	s := string(respBytes)
	if s != resp {
		if strings.Contains(s, `"data"`) {
			return
		}
		//debug.PrintStack()
		t.Fatalf("got: %s\nwant: %s\nlen(%d/%d)", string(respBytes), resp,
			len(s), len(resp))
	}
}

func TestForcastOneDayWeather(t *testing.T) {
	// setup mock server
	starter.Do(func() { setupMockServer(t) })
	// LaunchRequest
	runResponseTest(t, launchResponse, func() *sp.RequestEnvelope {
		req := sp.NewLaunchRequest(reqId, ts)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
	// IntentRequest
	// 1. slot: date
	date := time.Now().Format("2006-01-02")
	respStr := fmt.Sprintf(respElicitToCity, SlotDate, date)
	runResponseTest(t, respStr, func() *sp.RequestEnvelope {
		intent := slu.NewIntent(IntentSearchOneDay).
			WithSlot(slu.NewSlot(SlotDate).WithStringValue(date)).
			WithSlot(slu.NewSlot(SlotFocus).WithStringValue(WeatherFocus))
		req := sp.NewIntentRequest(reqId, ts, intent)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
	log.Println("pass runResponseTest for slot date")
	// 1. slot: city
	respStr = fmt.Sprintf(respDialogCompleted, SlotDate, date)
	runResponseTest(t, respStr, func() *sp.RequestEnvelope {
		intent := slu.NewIntent(IntentSearchOneDay).
			WithSlot(slu.NewSlot(SlotCity).WithStringValue("北京"))
		req := sp.NewIntentRequest(reqId, ts, intent)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
	log.Println("pass runResponseTest for slot city")
}

func TestForcastDaysWeather(t *testing.T) {
	// setup mock server
	starter.Do(func() { setupMockServer(t) })
	// IntentRequest
	// 1. slot: date
	start := time.Now().Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 3).Format("2006-01-02")
	duration := start + "/" + end
	respStr := fmt.Sprintf(respElicitToCity, SlotDuration, duration)
	runResponseTest(t, respStr, func() *sp.RequestEnvelope {
		intent := slu.NewIntent(IntentSearchDays).
			WithSlot(slu.NewSlot(SlotDuration).WithStringValue(duration)).
			WithSlot(slu.NewSlot(SlotFocus).WithStringValue(WeatherFocus))
		req := sp.NewIntentRequest(reqId, ts, intent)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
	log.Println("pass runResponseTest for slot date")
	// 1. slot: city
	respStr = fmt.Sprintf(respDialogCompleted, SlotDuration, duration)
	runResponseTest(t, respStr, func() *sp.RequestEnvelope {
		intent := slu.NewIntent(IntentSearchDays).
			WithSlot(slu.NewSlot(SlotCity).WithStringValue("北京"))
		req := sp.NewIntentRequest(reqId, ts, intent)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
	log.Println("pass runResponseTest for slot city")
}

func TestMakeDaysTempText(t *testing.T) {
	//
	results := model.NewResults().
		Append(model.NewResult(1).WithMaxTemp("10").WithMinTemp("6").WithDate("2018-05-14")).
		Append(model.NewResult(2).WithMaxTemp("12").WithMinTemp("3").WithDate("2018-05-15")).
		Append(model.NewResult(3).WithMaxTemp("15").WithMinTemp("1").WithDate("2018-05-16")).
		Append(model.NewResult(4).WithMaxTemp("12").WithMinTemp("-2").WithDate("2018-05-17")).
		Append(model.NewResult(5).WithMaxTemp("15").WithMinTemp("6").WithDate("2018-05-18")).
		Append(model.NewResult(6).WithMaxTemp("14").WithMinTemp("-3").WithDate("2018-05-19")).
		Append(model.NewResult(7).WithMaxTemp("14").WithMinTemp("-3").WithDate("2018-05-20"))
	text, err := makeDaysTempText(results)
	if err != nil {
		t.Fatal(err)
	}
	if text != "最高气温15度，出现在5月16日，5月18日，最低气温-3度，出现在5月19日至5月20日" {
		t.Fatalf("got: %s, want: %s", text, "最高气温15度，出现在5月16日，5月18日，最低气温-3度，出现在5月19日至5月20日")
	}
}

func TestMakeAbnormalWeatherText(t *testing.T) {
	results := model.NewResults().
		Append(model.NewResult(1).WithWeather("晴").WithDate("2018-05-14")).
		Append(model.NewResult(2).WithWeather("阵雨").WithDate("2018-05-15")).
		Append(model.NewResult(3).WithWeather("阵雨").WithDate("2018-05-16")).
		Append(model.NewResult(4).WithWeather("阵雨").WithDate("2018-05-17")).
		Append(model.NewResult(5).WithWeather("晴").WithDate("2018-05-18")).
		Append(model.NewResult(6).WithWeather("沙尘").WithDate("2018-05-19")).
		Append(model.NewResult(7).WithWeather("暴雪").WithDate("2018-05-20"))
	text, err := makeAbnormalWeatherText(results)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(text)
	if text != "其中5月15日至5月17日阵雨，5月19日沙尘，5月20日暴雪" {
		t.Fatalf("got: %s, want: %s", text, "其中5月15日至5月17日阵雨，5月19日沙尘，5月20日暴雪")
	}
}

func TestMakeAqiSpeechText(t *testing.T) {
	result := model.NewResult(1).WithCity("朝阳,北京,中国").
		WithDate(time.Now().Format("2006-01-02")).WithPm25("168")
	text, err := makeAqiSpeechText(result)
	if err != nil {
		t.Fatal(err)
	}
	v, _ := util.GetCfgVal([]string{}, "weather", "aqiText", "中度污染")
	ss := util.Conv2StrSlice(v.([]interface{}))
	var (
		texts []string
		s     = "北京朝阳今天，PM2.5是168，空气质量为中度污染，"
	)
	for _, v := range ss {
		texts = append(texts, s+v)
	}
	if !util.ContainsStr(texts, text) {
		t.Fatalf("got: %s not in strings: %v", text, texts)
	}
}

func runGetCityFromSysInfoTest(t *testing.T, want string, ctx *sp.Context) {
	city := getCityFromSysInfo(ctx)
	if city != want {
		t.Fatalf("want: %s, got: %s", want, city)
	}
}

func NewTestGetCityCtx(lat, log float64, prov, city string) *sp.Context {
	params := sp.NewCtxParams().WithParameter("location",
		slu.NewValue(slu.MapType,
			map[string]interface{}{
				"latitude":  lat,
				"longitude": log,
				"address": &model.Address{
					Country:  "中国",
					Province: prov,
					City:     city,
					Detail:   "",
				},
			}))
	sys := sp.NewCtxSystem().WithParameters(params)
	ctx := sp.NewContext().WithSystem(sys)
	return ctx
}

func TestGetCityFromSysInfo(t *testing.T) {
	runGetCityFromSysInfoTest(t, "北京",
		NewTestGetCityCtx(39.919815, 116.43324, "北京市", "北京市"))
	runGetCityFromSysInfoTest(t, "内蒙古阿拉善盟",
		NewTestGetCityCtx(39.919815, 116.43324, "内蒙古自治区", "阿拉善盟"))
	runGetCityFromSysInfoTest(t, "湖北武汉",
		NewTestGetCityCtx(39.919815, 116.43324, "湖北省", "武汉市"))
	//runGetCityFromSysInfoTest(t, "北京",
	//	NewTestGetCityCtx(39.919815, 116.43324, "", ""))
}
