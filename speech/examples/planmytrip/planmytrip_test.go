package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
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
	router   = "/planmytrip"
	userId   = "rosai1.ask.account.001"
	appId    = "rosai1.ask.developer.001"
	skillId  = "rosai1.ask.skill.planmytrip.v1.0"
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
            "source": "Welcome to the Rosai Skills Kit, you can say hello"
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
      "fromCity": {
        "orgin": null,
        "normType": "String",
        "norm": "Beijing"
      },
      "travelDate": {
        "orgin": null,
        "normType": "String",
        "norm": "2018-04-11"
      }
    }
  },
  "results": [
    {
      "outputSpeech": {
        "items": [
          {
            "type": "PlainText",
            "source": "Where are you traveling to?"
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
      "fromCity": {
        "orgin": null,
        "normType": "String",
        "norm": "Beijing"
      },
      "toCity": {
        "orgin": null,
        "normType": "String",
        "norm": "Seattle"
      },
      "travelDate": {
        "orgin": null,
        "normType": "String",
        "norm": "2018-04-11"
      }
    }
  },
  "results": [
    {
      "outputSpeech": {
        "items": [
          {
            "type": "PlainText",
            "source": "your trip from Beijing to Seattle on 2018-04-11 have been planned, enjoy yourself!"
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
		AppId:       "rosai1.ask.skill.planmytrip.12345",
		Speechlet:   &PlanMyTrip{},
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
	respBytes, err := snet.Post(url, reqBytes)
	if string(respBytes) != resp {
		//debug.PrintStack()
		t.Fatalf("got: %s\nwant: %s", string(respBytes), resp)
	}
}

func TestRequestDispatch(t *testing.T) {
	// setup mock server
	starter.Do(func() { setupMockServer(t) })
	// request info
	reqId, ts := "12345", "2018-04-06T15:30:02+08:00"
	sysInfo := sp.NewCtxSystem().
		WithUser(sp.NewUser(userId, appId)).
		WithSkill(sp.NewSkill(skillId)).
		WithDevice(sp.NewDevice(deviceId))
	ctx := sp.NewContext().WithSystem(sysInfo)
	// LaunchRequest
	runResponseTest(t, launchResponse, func() *sp.RequestEnvelope {
		req := sp.NewLaunchRequest(reqId, ts)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
	// IntentRequest
	// 1. slot: travelDate
	runResponseTest(t, respElicitToCity, func() *sp.RequestEnvelope {
		intent := slu.NewIntent(IntentPlanMyTrip).WithSubName("TestSubIntent").
			WithSlot(slu.NewSlot(SlotTravelDate).WithStringValue("2018-04-11"))
		req := sp.NewIntentRequest(reqId, ts, intent)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
	log.Println("pass runResponseTest for slot travelDate")
	// 1. slot: toCity
	runResponseTest(t, respDialogCompleted, func() *sp.RequestEnvelope {
		intent := slu.NewIntent(IntentPlanMyTrip).
			WithSlot(slu.NewSlot(SlotToCity).WithStringValue("Seattle"))
		req := sp.NewIntentRequest(reqId, ts, intent)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
	log.Println("pass runResponseTest for slot toCity")
}
