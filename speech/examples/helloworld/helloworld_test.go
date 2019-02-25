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

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/dialog/model"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
	sp "roobo.com/rosai-skills-kit-sdk-for-go/speech/speechlet"
)

var (
	addr    net.Addr
	starter sync.Once
)

const (
	router   = "/helloworld"
	userId   = "rosai1.ask.account.001"
	appId    = "rosai1.ask.developer.001"
	skillId  = "rosai1.ask.skill.helloworld.v1.0"
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

	intentResponse = `{
  "version": "2.0",
  "status": {
    "code": 0
  },
  "context": {
    "parameters": {
      "date": {
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
            "source": "Hello world"
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
	dlg := model.NewDialog(model.NewIntent("HelloWorldIntent", false))
	rh := sp.RequestHandler{
		AppId:       "rosai1.ask.skill.helloworld.12345",
		Speechlet:   &HelloWorld{},
		DialogModel: model.NewDialogModel().WithDialog(dlg),
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
	reqId, ts := "12345", "2018-04-06T15:30:02+08:00"
	sysInfo := sp.NewCtxSystem().
		WithUser(sp.NewUser(userId, appId)).
		WithSkill(sp.NewSkill(skillId)).
		WithDevice(sp.NewDevice(deviceId))
	ctx := sp.NewContext().WithSystem(sysInfo)

	starter.Do(func() { setupMockServer(t) })
	// LaunchRequest
	runResponseTest(t, launchResponse, func() *sp.RequestEnvelope {
		req := sp.NewLaunchRequest(reqId, ts)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
	// IntentRequest
	runResponseTest(t, intentResponse, func() *sp.RequestEnvelope {
		intent := slu.NewIntent("HelloWorldIntent").
			WithSlot(slu.NewSlot("date").WithStringValue("2018-04-11")).
			WithSlot(slu.NewSlot("city"))
		req := sp.NewIntentRequest(reqId, ts, intent)
		return sp.NewRequestEnvelope().WithContext(ctx).WithRequest(req)
	})
}
