package db

import (
	"log"
	"testing"
	"time"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/examples/weather/model"
	"roobo.com/sailor/util"
)

func init() {
	if err := util.InitConf("../conf/app.json"); err != nil {
		log.Fatal(err)
	}
	if err := ConfRedis(); err != nil {
		log.Fatalf(err.Error())
	}
}

var (
	result = model.NewResult(1).WithWeather("多云转晴").
		WithCity("北京").WithDate(time.Now().Format("2006-01-02")).
		WithTemperature("34").WithMaxTemp("38").WithMinTemp("22").
		WithWindDir("东南").WithWindLevel("3").
		WithWindDay("东南").WithWindDayLevel("3").
		WithWindNight("东南").WithWindNightLevel("3").
		WithPm25("25").WithHumidity("40")
	result2 = model.NewResult(2).WithWeather("阵雨转晴").
		WithCity("北京").WithDate(time.Now().AddDate(0, 0, 1).Format("2006-01-02")).
		WithTemperature("30").WithMaxTemp("36").WithMinTemp("12").
		WithWindDir("东").WithWindLevel("6").
		WithWindDay("南").WithWindDayLevel("2").
		WithWindNight("东南").WithWindNightLevel("4").
		WithPm25("10").WithHumidity("70")
	results = model.NewResults().Append(result, result2)
)

func TestStoreAndRestoreResults(t *testing.T) {
	today := time.Now()
	// Test one day result
	if err := StoreOneDayResult2Redis("北京", today, result); err != nil {
		t.Fatal(err)
	}
	resultRe, err := RestoreOneDayResult4Redis("北京", today)
	if err != nil {
		t.Fatal(err)
	}
	if *resultRe != *result {
		t.Fatalf("want: %+v, got: %+v", result, resultRe)
	}
	// Test days results
	if err := StoreDaysResults2Redis("北京", today, today.AddDate(0, 0, 1), results); err != nil {
		t.Fatal(err)
	}
	resultsRe, err := RestoreDaysResults4Redis("北京", today, today.AddDate(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range results {
		if *v != *resultsRe[i] {
			t.Fatalf("want: %+v, got: %+v", result, resultRe)
		}
	}
}
