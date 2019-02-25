package main

import (
	"errors"
	"fmt"
	"log"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
	sp "roobo.com/rosai-skills-kit-sdk-for-go/speech/speechlet"
)

const (
	IntentTidePooler = "SearchTidePooler"
	SlotCity         = "city"
	SlotDate         = "date"
)

type TidePooler struct {
}

func (hw *TidePooler) OnSessionStarted(re *sp.RequestEnvelope) error {
	log.Printf("INFO] OnSessionStarted requestId=%s, sessionId=%s", "rosai.123", "helloworld123")
	return nil
}

func (hw *TidePooler) OnSessionEnded(request *sp.RequestEnvelope) error {
	log.Printf("INFO] OnSessionEnded requestId=%s, sessionId=%d", "rosai.4", "helloworld.4")
	return nil
}

func (hw *TidePooler) OnLaunch(re *sp.RequestEnvelope) (*sp.Response, error) {
	log.Printf("INFO] OnLaunch requestId=%s, sessionId=%s", "rosai.123", "helloworld123")
	return getWelcomeResponse(), nil
}

func (hw *TidePooler) OnIntent(re *sp.RequestEnvelope) (resp *sp.Response,
	ctx *sp.Context, err error) {
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
	inCtx := re.Context
	log.Printf("intent name: %s, context: %+v", intentName, ctx)
	switch intentName {
	case IntentTidePooler:
		return handleSearchTidePoolerIntent(intent, inCtx)
	case "ROSAI.HelpIntent":
		resp, ctx, err = getHelpResponse()
		return
	default:
		tip := fmt.Sprintf("Intent(%s) is unsupported.  Please try something else.", intentName)
		resp = sp.NewAskResponse(tip)
	}
	return resp, nil, err
}

func handleSearchTidePoolerIntent(intent *slu.Intent, inCtx *sp.Context) (
	*sp.Response, *sp.Context, error) {
	city := intent.GetSlot(SlotCity).GetStringValue()
	date := intent.GetSlot(SlotDate).GetStringValue()
	switch {
	default:
		resp, _, err := getHelpResponse()
		return resp, nil, err
	case city != "" && date != "":
		return getFinalTideResponse(city, date)
	case city != "":
		return handleCityDialogRequest(intent, inCtx)
	case date != "":
		return handleDateDialogRequest(intent, inCtx)
	}
}

func handleCityDialogRequest(intent *slu.Intent, inCtx *sp.Context) (
	*sp.Response, *sp.Context, error) {
	date := inCtx.GetStringValue(SlotDate)
	if date == "" {
		speechText := "你要查询哪一天的潮汐信息"
		resp := sp.NewAskResponse(speechText)
		return resp, nil, nil
	}
	city := intent.GetSlot(SlotCity).GetStringValue()
	return getFinalTideResponse(city, date)
}

func handleDateDialogRequest(intent *slu.Intent, inCtx *sp.Context) (
	*sp.Response, *sp.Context, error) {
	city := inCtx.GetStringValue(SlotCity)
	if city == "" {
		speechText := "你要查询哪个城市的潮汐信息"
		resp := sp.NewAskResponse(speechText)
		return resp, nil, nil
	}
	date := intent.GetSlot(SlotDate).GetStringValue()
	return getFinalTideResponse(city, date)
}

func getFinalTideResponse(city, date string, extra ...string) (*sp.Response,
	*sp.Context, error) {
	speechText := city + date + " 6点涨潮，18点退潮"
	resp := sp.NewTellResponse(speechText)
	return resp, nil, nil
}

func getHelpResponse() (resp *sp.Response, ctx *sp.Context, err error) {
	speechText := "你可以跟我说\"我要查今天北京的潮汐信息\""
	resp = sp.NewAskResponse(speechText)
	ctx = sp.NewContext()
	return
}

func getWelcomeResponse() *sp.Response {
	speechText := "Welcome to the Rosai Skills Kit, you can say hello"
	return sp.NewAskResponse(speechText)
}
