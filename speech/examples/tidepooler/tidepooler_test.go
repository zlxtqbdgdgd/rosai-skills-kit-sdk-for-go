package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"testing"

	snet "roobo.com/sailor/net"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
	sp "roobo.com/rosai-skills-kit-sdk-for-go/speech/speechlet"
)

var (
	addr    net.Addr
	starter sync.Once
)

const (
	router   = "/tidepooler"
	userId   = "rosai1.ask.account.001"
	appId    = "rosai1.ask.developer.001"
	skillId  = "rosai1.ask.skill.tidepooler.v1.0"
	deviceId = "rosai1.device.001"

	launchResponse = `{
  "version": "2.0",
  "status": {
    "code": 0
  },
  "results": [
    {
      "formatType": "SpeechFormat",
      "hint": "Welcome to the Rosai Skills Kit, you can say hello"
    }
  ]
}`
	finalResponse = `{
  "version": "2.0",
  "status": {
    "code": 0
  },
  "context": {
    "parameters": {
      "city": {
        "orgin": null,
        "normType": "String",
        "norm": "上海"
      },
      "date": {
        "orgin": null,
        "normType": "String",
        "norm": "2018-04-12"
      }
    }
  },
  "results": [
    {
      "outputSpeech": {
        "items": [
          {
            "type": "PlainText",
            "source": "上海2018-04-12 6点涨潮，18点退潮"
          }
        ]
      }
    }
  ]
}`
	withDateAskCityResponse = `{
  "version": "2.0",
  "status": {
    "code": 0
  },
  "context": {
    "parameters": {
      "date": {
        "orgin": null,
        "normType": "String",
        "norm": "2018-04-12"
      }
    }
  },
  "results": [
    {
      "outputSpeech": {
        "items": [
          {
            "type": "PlainText",
            "source": "你要查询哪个城市的潮汐信息"
          }
        ]
      }
    }
  ]
}`
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
		AppId:       "rosai1.ask.skill.tidepooler.12345",
		Speechlet:   &TidePooler{},
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
	reqBytes, err := json.Marshal(reqEn)
	if err != nil {
		t.Fatal(err)
	}
	url := "http://" + addr.String() + router
	respBytes, err := snet.Post(url, reqBytes)
	if string(respBytes) != resp {
		debug.PrintStack()
		t.Fatalf("got: %s, want: %s", string(respBytes), resp)
	}
}

func TestRequestDispatch(t *testing.T) {
	starter.Do(func() { setupMockServer(t) })
	reqId, ts := "tidepooler.12345", "2018-04-06T15:30:02+08:00"
	sysInfo := sp.NewCtxSystem().
		WithUser(sp.NewUser(userId, appId)).
		WithSkill(sp.NewSkill(skillId)).
		WithDevice(sp.NewDevice(deviceId))
	ctx := sp.NewContext().WithSystem(sysInfo)
	// LaunchRequest
	//runResponseTest(t, launchResponse, func() sp.RequestEnvelope {
	//	req := sp.NewLaunchRequest(reqId, locale, ts)
	//	return sp.RequestEnvelope{
	//		Version: "0.1",
	//		//Session: &sp.Session{
	//		//	IsNew:     true,
	//		//	SessionId: "rosai.api.session.123-456-789",
	//		//	Application: sp.Application{
	//		//		ApplicationId: "rosai.ask.skill.100-100-100",
	//		//	},
	//		//},
	//		Request: req,
	//	}
	//})
	//// OneShot IntentRequest
	//runResponseTest(t, finalResponse, func() *sp.RequestEnvelope {
	//	intent := slu.NewIntent(IntentTidePooler).
	//		WithSlot(slu.NewSlot(SlotCity).WithValue("上海")).
	//		WithSlot(slu.NewSlot(SlotDate).WithValue("2018-04-12"))
	//	req := sp.NewIntentRequest(reqId, locale, ts, intent, slu.STARTED)
	//	return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	//})
	// IntentRequest with Date
	runResponseTest(t, withDateAskCityResponse, func() *sp.RequestEnvelope {
		intent := slu.NewIntent(IntentTidePooler).
			WithSlot(slu.NewSlot(SlotDate).WithStringValue("2018-04-12"))
		req := sp.NewIntentRequest(reqId, ts, intent)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
	// IntentRequest with Date and context with city
	runResponseTest(t, finalResponse, func() *sp.RequestEnvelope {
		intent := slu.NewIntent(IntentTidePooler).
			WithSlot(slu.NewSlot(SlotCity).WithStringValue("上海"))
		req := sp.NewIntentRequest(reqId, ts, intent)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
}
