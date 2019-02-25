package main

import (
	"errors"
	"fmt"
	"log"

	sp "roobo.com/rosai-skills-kit-sdk-for-go/speech/speechlet"
)

type HelloWorld struct {
}

func (hw *HelloWorld) OnSessionStarted(re *sp.RequestEnvelope) error {
	log.Printf("INFO] OnSessionStarted requestId=%s, sessionId=%s", "rosai.123", "helloworld123")
	return nil
}

func (hw *HelloWorld) OnLaunch(re *sp.RequestEnvelope) (*sp.Response, error) {
	log.Printf("INFO] OnLaunch requestId=%s, sessionId=%s", "rosai.123", "helloworld123")
	return getWelcomeResponse(), nil
}

func getWelcomeResponse() *sp.Response {
	speechText := "Welcome to the Rosai Skills Kit, you can say hello"
	return sp.NewAskResponse(speechText)
}

func (hw *HelloWorld) OnIntent(re *sp.RequestEnvelope) (
	*sp.Response, *sp.Context, error) {
	request, ok := re.Request.(*sp.IntentRequest)
	if !ok {
		log.Printf("INFO] re.Request type: %T, value: %+v", re.Request, re.Request)
		return nil, nil, errors.New("OnIntent assert requestEnvelope for IntentRequest failed")
	}
	log.Printf("INFO] onIntent requestId=%s", request.RequestId)
	intent := request.Intent
	intentName := ""
	if intent != nil {
		intentName = intent.Name
	}
	log.Println("intent name:", intentName)
	switch intentName {
	case "HelloWorldIntent":
		return getHelloResponse()
	case "ROSAI.HelpIntent":
		return getHelpResponse()
	default:
		tip := fmt.Sprintf("Intent(%s) is unsupported.  Please try something else.", intentName)
		resp := sp.NewAskResponse(tip)
		return resp, nil, nil
	}
}

func getHelloResponse() (*sp.Response, *sp.Context, error) {
	speechText := "Hello world"
	resp := sp.NewTellResponse(speechText)
	return resp, nil, nil
}

func getHelpResponse() (*sp.Response, *sp.Context, error) {
	speechText := "You can say hello to me!"
	resp := sp.NewAskResponse(speechText)
	return resp, nil, nil
}

func (hw *HelloWorld) OnSessionEnded(request *sp.RequestEnvelope) error {
	log.Printf("INFO] OnSessionEnded requestId=%s, sessionId=%d", "rosai.4", "helloworld.4")
	return nil
}
