package main

import (
	"errors"
	"fmt"
	"log"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/dialog/directives"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
	sp "roobo.com/rosai-skills-kit-sdk-for-go/speech/speechlet"
)

const (
	IntentPlanMyTrip = "PlanMyTrip"
	SlotTravelDate   = "travelDate"
	SlotToCity       = "toCity"
	SlotFromCity     = "fromCity"
	SlotActivity     = "activity"
)

type PlanMyTrip struct {
	fromCity, toCity, travelDate, activity string
}

func (pmt *PlanMyTrip) UpdateSlotsValues(intent *slu.Intent) {
	pmt.fromCity = intent.GetSlot(SlotFromCity).GetStringValue()
	pmt.toCity = intent.GetSlot(SlotToCity).GetStringValue()
	pmt.travelDate = intent.GetSlot(SlotTravelDate).GetStringValue()
	pmt.activity = intent.GetSlot(SlotActivity).GetStringValue()
}

func (pmt *PlanMyTrip) OnSessionStarted(re *sp.RequestEnvelope) error {
	log.Printf("INFO] OnSessionStarted requestId=%s, sessionId=%s", "rosai.123", "helloworld123")
	return nil
}

func (pmt *PlanMyTrip) OnLaunch(re *sp.RequestEnvelope) (*sp.Response, error) {
	log.Printf("INFO] OnLaunch requestId=%s, sessionId=%s", "rosai.123", "helloworld123")
	return getWelcomeResponse()
}

func getWelcomeResponse() (*sp.Response, error) {
	speechText := "Welcome to the Rosai Skills Kit, you can say hello"
	return sp.NewAskResponse(speechText), nil
}

func (pmt *PlanMyTrip) OnIntent(re *sp.RequestEnvelope) (
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
	case "PlanMyTrip":
		dialogState := request.DialogState
		if dialogState == slu.STARTED {
			intent.SetSlot(slu.NewSlot(SlotFromCity).WithStringValue("Beijing"))
			directive := directives.NewDelegateDirective(intent)
			directives := []directives.Directive{directive}
			return sp.NewDelegateResponse(directives), nil, nil
		} else if dialogState == slu.COMPLETED {
			pmt.UpdateSlotsValues(intent)
			return pmt.getTellResponse()
		} else {
			// This is executed when the dialog is in state e.g. IN_PROGESS.
			// If there is only one slot this shouldn't be called
			directive := directives.NewDelegateDirective(intent)
			directives := []directives.Directive{directive}
			return sp.NewDelegateResponse(directives), nil, nil
		}
	case "ROSAI.HelpIntent":
		return getHelpResponse()
	default:
		tip := fmt.Sprintf("Intent(%s) is unsupported. Please try something else.", intentName)
		return sp.NewAskResponse(tip), nil, nil
	}
}

func (pmt *PlanMyTrip) getTellResponse() (*sp.Response, *sp.Context, error) {
	var speechText string
	if pmt.activity == "" {
		speechText = fmt.Sprintf("your trip from %s to %s on %s have been planned, "+
			"enjoy yourself!", pmt.fromCity, pmt.toCity, pmt.travelDate)
	} else {
		speechText = fmt.Sprintf("your trip from %s to %s on %s for %s have been planned, "+
			"enjoy yourself!", pmt.fromCity, pmt.toCity, pmt.travelDate, pmt.activity)
	}
	return sp.NewTellResponse(speechText), nil, nil
}

func getHelpResponse() (*sp.Response, *sp.Context, error) {
	speechText := "You can say hello to me!"
	return sp.NewAskResponse(speechText), nil, nil
}

func (pmt *PlanMyTrip) OnSessionEnded(request *sp.RequestEnvelope) error {
	log.Printf("INFO] OnSessionEnded requestId=%s, sessionId=%d", "rosai.4", "helloworld.4")
	return nil
}
