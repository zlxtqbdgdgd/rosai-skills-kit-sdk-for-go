package speechlet

import (
	"bytes"
	"encoding/json"
	"testing"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/dialog/directives"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/ui"
)

const (
	respFormat = `{
  "directives": [
    {
      "type": "Dialog.Delegate",
      "updatedIntent": {
        "name": "PlanMyTrip",
        "confirmationStatus": "NONE",
        "slots": {
          "fromCity": {
            "name": "fromCity",
            "value": {
              "orgin": null,
              "normType": "String",
              "norm": "Beijing"
            },
            "confirmationStatus": "CONFIRMED"
          },
          "toCity": {
            "name": "toCity",
            "value": {
              "orgin": null,
              "normType": "String",
              "norm": "Sanya"
            },
            "confirmationStatus": "CONFIRMED"
          },
          "travelDate": {
            "name": "travelDate",
            "value": {
              "orgin": null,
              "normType": "String",
              "norm": "2018-04-05"
            },
            "confirmationStatus": "NONE"
          }
        }
      }
    }
  ],
  "shouldEndSession": true
}`

	respWithDisplayDirective = `{
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
      "date": {
        "orgin": null,
        "normType": "String",
        "norm": "2018-06-21"
      }
    }
  },
  "results": [
    {
      "outputSpeech": {
        "items": [
          {
            "type": "PlainText",
            "source": "北京今天多云，气温23度到35度，东南风2级"
          },
          {
            "type": "SSML",
            "source": "\u003cspeak\u003e北京今天\u003cemphasis level=\"strong\"\u003e多云\u003c/emphasis\u003e，气温23度到35度，东南风2级\u003c/speak\u003e"
          },
          {
            "type": "Audio",
            "source": "https://ai.roobo.com/weather/wind_2.mp3"
          },
          {
            "type": "PlainText",
            "source": "您还可以跟我说 北京空气质量?"
          }
        ]
      },
      "directives": [
        {
          "type": "Display.Customized",
          "hint": "多云，气温23度到35度，东南风2级",
          "card": {
            "type": "Standard",
            "title": "北京天气",
            "content": "多云，气温23度到35度，东南风2级",
            "image": {
              "url": "www.roobo.com/aicloud/skills/weather/cloudy.jpg"
            }
          },
          "suggestions": [
            "明天",
            "后天"
          ]
        },
        {
          "type": "ROSAI.EVENT",
          "event": {
            "name": "ROSAI.ContinueEvent",
            "period": 5000
          }
        }
      ],
      "data": {
        "city": "北京",
        "date": "2018-06-21",
        "quality": "轻度污染",
        "temperature": "33",
        "weather": "多云"
      }
    }
  ]
}`
)

func TestResponseFormat(t *testing.T) {
	intent := slu.NewIntent("PlanMyTrip").
		WithSlot(slu.NewSlot("fromCity").WithStringValue("Beijing").WithStatus(slu.CONFIRMED)).
		WithSlot(slu.NewSlot("toCity").WithStringValue("Sanya").WithStatus(slu.CONFIRMED)).
		WithSlot(slu.NewSlot("travelDate").WithStringValue("2018-04-05"))
	directives := []directives.Directive{directives.NewDelegateDirective(intent)}
	sr := Response{
		Directives:       directives,
		ShouldEndSession: true,
	}
	bytes, err := json.MarshalIndent(sr, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != respFormat {
		t.Fatalf("response format want: %s(len:%d), got: %s(len:%d)",
			respFormat, len(respFormat), string(bytes), len(bytes))
	}
}

func TestResponseEnvWithOutputDisplay(t *testing.T) {
	status := NewGoodStatus()
	ctx := NewContext().WithStringValue("city", "北京").WithStringValue("date", "2018-06-21")

	hint := "北京今天多云，气温23度到35度，东南风2级"
	result := NewResult().
		WithOutputSpeech(NewSpeechItems(ui.NewPlainTextSpeechItem(hint),
			ui.NewSsmlSpeechItem("<speak>北京今天"+
				"<emphasis level=\"strong\">多云</emphasis>，气温23度到35度，东南风2级</speak>"),
			ui.NewAudioSpeechItem("https://ai.roobo.com/weather/wind_2.mp3"),
			ui.NewPlainTextSpeechItem("您还可以跟我说 北京空气质量?"))).
		WithDisplayDirective(NewDisplayDirective().
			WithHint("多云，气温23度到35度，东南风2级").
			WithCard(ui.NewStandardCard("北京天气", "多云，气温23度到35度，东南风2级").
				WithImage("www.roobo.com/aicloud/skills/weather/cloudy.jpg")).
			WithSuggestions("明天", "后天")).
		WithEventDirective(NewEventDirective("ROSAI.ContinueEvent", 5000)).
		WithData(map[string]interface{}{
			"city":        "北京",
			"date":        "2018-06-21",
			"weather":     "多云",
			"quality":     "轻度污染",
			"temperature": "33",
		})

	respEn := NewResponseEnvelope().WithStatus(status).WithContext(ctx).WithResults(result)

	bytesl, _ := json.MarshalIndent(respEn, "", "  ")
	if string(bytesl) != respWithDisplayDirective {
		t.Fatalf("want: %s\n, got: %s\n(%d/%d)", respWithDisplayDirective, string(bytesl),
			len(respWithDisplayDirective), len(string(bytesl)))
	}
	// test ResponseEnvelopeRaw
	var raw ResponseEnvelopeRaw
	if err := json.Unmarshal([]byte(respWithDisplayDirective), &raw); err != nil {
		t.Fatal(err)
	}
	results := raw.GetResults()
	if results == nil {
		t.Fatal("get results is nil")
	}
	bytesl, _ = json.MarshalIndent(results, "", "  ")
	bytesC, _ := json.MarshalIndent(respEn.Results, "", "  ")
	if !bytes.Equal(bytesl, bytesC) {
		t.Fatal("ResponseEnvelopeRaw unmarshal failed, result is not equal to origin")
	}
}
