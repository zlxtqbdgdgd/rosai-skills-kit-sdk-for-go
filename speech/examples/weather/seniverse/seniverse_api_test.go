package seniverse

import (
	"log"
	"strings"
	"testing"
	"time"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/examples/weather/db"
	"roobo.com/sailor/util"
)

func init() {
	if err := util.InitConf("../conf/app.json"); err != nil {
		log.Fatal(err)
	}
	err1 := ConfWeatherApi()
	err2 := ConfMysql()
	err3 := db.ConfRedis()
	if err1 != nil || err2 != nil || err3 != nil {
		util.NewErrf("init conf error: %v,%v,%v", err1, err2, err3)
	}
}

func runGetCityTest(t *testing.T, city string, want string) {
	id, err := GetCityID(city)
	if err != nil {
		t.Fatalf("city: %s, error: %s", city, err)
	}
	if id != want {
		t.Fatalf("got city: %s id: %s, want: %s", city, id, want)
	}
}

func TestGetCityID(t *testing.T) {
	runGetCityTest(t, "北京", "WX4FBXXFKE4F")
	runGetCityTest(t, "长沙", "WT029G15ETRJ")
	runGetCityTest(t, "北京朝阳", "WX4G17JWZEK7")
	runGetCityTest(t, "岳阳湘阴", "WT070TH9NEKX")
	runGetCityTest(t, "辽宁朝阳", "WXMSKZT4B3TT")
	runGetCityTest(t, "北京通州", "WX4GN7JYY527")
	runGetCityTest(t, "江苏通州", "WTWN7PBRX3RJ")
	runGetCityTest(t, "天津", "WWGQDCW6TBW1")
	//runGetCityTest(t, "朝阳", "WX4G17JWZEK7")
}

func runGetCityNameTest(t *testing.T, lat, log float64, want string) {
	id, err := GetCityName(lat, log)
	if err != nil {
		t.Fatalf("lat: %f, log: %f, error: %s", lat, log, err)
	}
	if id != want {
		t.Fatalf("got lat: %f, log: %f, name: %s, want: %s", lat, log, id, want)
	}
}

func TestGetCityName(t *testing.T) {
	runGetCityNameTest(t, 39.919815, 116.43324, "北京")
	runGetCityNameTest(t, 29.645717, 106.469728, "重庆")
	runGetCityNameTest(t, 39.921282, 116.45991, "北京朝阳")
	runGetCityNameTest(t, 28.669832, 112.909645, "岳阳湘阴")
	runGetCityNameTest(t, 29.148504, 113.141516, "湖南岳阳")
}

func TestLogRespBody(t *testing.T) {
	raw := []byte{'a', 'b', 'c'}
	LogRespBody("prefix:", raw, 200)
	var raw2 []byte
	LogRespBody("prefix:", raw2, 200)
	raw = nil
	LogRespBody("prefix:", raw, 200)
	lRaw := make([]byte, 300)
	for i, _ := range lRaw {
		lRaw[i] = '*'
	}
	lRaw[199], lRaw[200] = 'a', 'd'
	LogRespBody("prefix:", lRaw, 200)
}

func runForcastOneDayTest(t *testing.T, city string, tm time.Time) {
	sw := &SeniverseWeather{}
	result, err := sw.ForcastOneDayWeather(city, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("result for city: %s, time: %v: %+v", city, tm, result)
	if time.Now().Format("2006-01-02") != result.Date ||
		!strings.Contains(result.City, city) || len(result.Weather) == 0 {
		t.Fatalf("weather date error: got: city: %s, date: %s, weather: %s,"+
			" want: city: %s, date: %s", result.City, result.Date, result.Weather,
			city, time.Now().Format("2006-01-02"))
	}
}

func TestForcastOneDayApi(t *testing.T) {
	runForcastOneDayTest(t, "长沙", time.Now().AddDate(0, 0, 1))
	runForcastOneDayTest(t, "南京", time.Now().AddDate(0, 0, 1))
	runForcastOneDayTest(t, "北京", time.Now())
}

func runSearchNow(t *testing.T, city string) {
	sw := &SeniverseWeather{}
	result, err := sw.SearchNow(city)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("weather result for city %s: %+v", city, result)
	if !strings.Contains(result.City, city) || len(result.Weather) == 0 {
		t.Fatalf("weather date error: got: city: %s, weather: %s, want: city: %s",
			result.City, result.Weather, city)
	}
}

func TestSearchNowApi(t *testing.T) {
	runSearchNow(t, "北京")
	runSearchNow(t, "长沙")
}

func runForcastOneDayAqiTest(t *testing.T, city string, tm time.Time) {
	sw := &SeniverseWeather{}
	result, err := sw.ForcastOneDayWeather(city, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("aqi result for city: %s, time: %v: %+v", city, tm, result)
	if time.Now().Format("2006-01-02") != result.Date || len(result.Pm25) == 0 ||
		!strings.Contains(result.City, city) {
		t.Fatalf("weather date error: got: city: %s, date: %s, pm25: %s,"+
			" want: city: %s, date: %s", result.City, result.Date, result.Pm25,
			city, time.Now().Format("2006-01-02"))
	}
}

func TestForcastOneDayAqiApi(t *testing.T) {
	runForcastOneDayAqiTest(t, "长沙", time.Now().AddDate(0, 0, 1))
	runForcastOneDayAqiTest(t, "南京", time.Now().AddDate(0, 0, 1))
	runForcastOneDayAqiTest(t, "北京", time.Now())
}
