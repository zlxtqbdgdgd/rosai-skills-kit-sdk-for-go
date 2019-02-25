package speechlet

import (
	"encoding/json"
	"testing"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/interfaces/system"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
)

const (
	ssStartedReqFormat = `{
  "type": "SessionStartedRequest",
  "requestId": "12345",
  "timestamp": "2018-04-06T15:30:02+08:00"
}`
	ssEndedReqFormat = `{
  "type": "SessionEndedRequest",
  "requestId": "12345",
  "timestamp": "2018-04-06T15:30:02+08:00",
  "reason": "USER_INITIATED",
  "error": {
    "type": "INTERNAL_SERVICE_ERROR",
    "message": "test internal service error"
  }
}`
	launchReqFormat = `{
  "type": "LaunchRequest",
  "requestId": "12345",
  "timestamp": "2018-04-06T15:30:02+08:00"
}`

	intentReqFormat = `{
  "type": "IntentRequest",
  "requestId": "12345",
  "timestamp": "2018-04-06T15:30:02+08:00",
  "intent": {
    "name": "PlanMyTrip",
    "confirmationStatus": "CONFIRMED",
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
  },
  "dialogState": "COMPLETED"
}`
)

const (
	IntentPlanMyTrip = "PlanMyTrip"
	SlotTravelDate   = "travelDate"
	SlotToCity       = "toCity"
	SlotFromCity     = "fromCity"
	SlotActivity     = "activity"
)

var (
	reqId, ts = "12345", "2018-04-06T15:30:02+08:00"
)

func runFormatTest(t *testing.T, prefix, reqRef string, f func() Request) {
	req := f()
	bytes, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != reqRef {
		t.Fatalf("%s format want: %s, got: %s", prefix, reqRef, string(bytes))
	}
}

func TestRequestFormat(t *testing.T) {
	// SessionStartedRequest
	runFormatTest(t, "ssStartedReqFormat", ssStartedReqFormat, func() Request {
		//ts := time.Now().Local().Format(time.RFC3339)
		return NewSessionStartedRequest(reqId, ts)
	})
	// SessionEndedRequest
	runFormatTest(t, "ssEndedReqFormat", ssEndedReqFormat, func() Request {
		return NewSessionEndedRequest(reqId, ts, USER_INITIATED,
			&system.Error{
				Type:    system.INTERNAL_SERVICE_ERROR,
				Message: "test internal service error",
			})
	})
	// LaunchRequest
	runFormatTest(t, "ssLaunchReqFormat", launchReqFormat, func() Request {
		return NewLaunchRequest(reqId, ts)
	})
	// IntentRequest
	runFormatTest(t, "ssIntentReqFormat", intentReqFormat, func() Request {
		intent := slu.NewIntent(IntentPlanMyTrip).WithStatus(slu.CONFIRMED).
			WithSlot(slu.NewSlot(SlotFromCity).WithStringValue("Beijing").WithStatus(slu.CONFIRMED)).
			WithSlot(slu.NewSlot(SlotToCity).WithStringValue("Sanya").WithStatus(slu.CONFIRMED)).
			WithSlot(slu.NewSlot(SlotTravelDate).WithStringValue("2018-04-05"))
		return NewIntentRequest(reqId, ts, intent).WithDialogState(slu.COMPLETED)
	})
}

func TestIntentRequestStatus(t *testing.T) {
	intent := slu.NewIntent(IntentPlanMyTrip).
		WithSlot(slu.NewSlot(SlotTravelDate).WithStringValue("2018-04-11"))
	intReq := NewIntentRequest(reqId, ts, intent)
	mi := rh.DialogModel.GetIntent(intent.Name)
	intReq.GetIntent().GetSlot(SlotTravelDate)
	intReq.GetIntent().WithSlot(slu.NewSlot(SlotToCity).WithStringValue("Seattle"))
	intReq.GetIntent().WithSlot(slu.NewSlot(SlotFromCity).WithStringValue("Beijing"))
	if !intent.Completed(mi) {
		intReqBytes, _ := json.MarshalIndent(intReq, "", "  ")
		miBytes, _ := json.MarshalIndent(mi, "", "  ")
		t.Fatalf("intent request should not set COMPLETED\nintent request: %s\n"+
			"model intent: %s", string(intReqBytes), string(miBytes))
	}
}
