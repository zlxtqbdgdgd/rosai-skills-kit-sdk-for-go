package speechlet

import (
	"roobo.com/rosai-skills-kit-sdk-for-go/speech"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/interfaces/system"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
)

type RequestType string

const (
	SessionStartedRequestType RequestType = "SessionStartedRequest"
	SessionEndedRequestType   RequestType = "SessionEndedRequest"
	LaunchRequestType         RequestType = "LaunchRequest"
	IntentRequestType         RequestType = "IntentRequest"
	IntentsRequestType        RequestType = "IntentsRequest"
)

type RequestEnvelope struct {
	Version string   `json:"version"`
	Context *Context `json:"context"`
	Request Request  `json:"request"`
}

func NewRequestEnvelope() *RequestEnvelope {
	return &RequestEnvelope{
		Version: speech.Version,
	}
}

func (re *RequestEnvelope) WithContext(ctx *Context) *RequestEnvelope {
	re.Context = ctx
	return re
}

func (re *RequestEnvelope) WithRequest(req Request) *RequestEnvelope {
	re.Request = req
	return re
}

type Request interface {
	GetType() RequestType
	GetRequestId() string
	GetTimestamp() string

	setBaseInfo(typ RequestType, id, ts string)
	SetType(typ RequestType)
	SetRequestId(reqId string)
	SetTimestamp(ts string)
}

func NewSessionStartedRequest(reqId, ts string) Request {
	ssr := new(SessionStartedRequest)
	ssr.speechletRequest = new(speechletRequest)
	ssr.setBaseInfo(SessionStartedRequestType, reqId, ts)
	return ssr
}

func NewSessionEndedRequest(reqId, ts string, reason Reason,
	err *system.Error) Request {
	ser := new(SessionEndedRequest)
	ser.speechletRequest = new(speechletRequest)
	ser.setBaseInfo(SessionEndedRequestType, reqId, ts)
	ser.Reason = reason
	if err != nil {
		ser.Err = err
	}
	return ser
}

func NewLaunchRequest(reqId, ts string) Request {
	lr := new(LaunchRequest)
	lr.speechletRequest = new(speechletRequest)
	lr.setBaseInfo(LaunchRequestType, reqId, ts)
	return lr
}

type speechletRequest struct {
	// e.g. SessionStartedRequest, LaunchRequest,IntentRequest,SessionEndedRequest
	Type RequestType `json:"type"`
	// eg. rosai1.pudding-api.request.82292535-81da-45bd-be0a-a6dbb443c3
	RequestId string `json:"requestId"`
	Timestamp string `json:"timestamp"` // eg. 2018-03-30T05:46:53Z
}

func (sr *speechletRequest) GetType() RequestType {
	return sr.Type
}

func (sr *speechletRequest) GetRequestId() string {
	return sr.RequestId
}

func (sr *speechletRequest) GetTimestamp() string {
	return sr.Timestamp
}

func (sr *speechletRequest) SetType(typ RequestType) {
	sr.Type = typ
}

func (sr *speechletRequest) SetRequestId(reqId string) {
	sr.RequestId = reqId
}

func (sr *speechletRequest) SetTimestamp(timestamp string) {
	sr.Timestamp = timestamp
}

func (sr *speechletRequest) setBaseInfo(typ RequestType, reqId, ts string) {
	sr.Type = typ
	sr.RequestId = reqId
	sr.Timestamp = ts
}

type CoreRequest struct {
	*speechletRequest
}

type SystemRequest struct {
	*speechletRequest
}

type SessionStartedRequest struct {
	CoreRequest
}

type LaunchRequest struct {
	CoreRequest
}

type IntentRequest struct {
	CoreRequest
	Intent      *slu.Intent     `json:"intent"`
	DialogState slu.DialogState `json:"dialogState,omitempty"`
}

func NewIntentRequest(reqId, ts string, intent *slu.Intent) *IntentRequest {
	ir := new(IntentRequest)
	ir.speechletRequest = new(speechletRequest)
	ir.setBaseInfo(IntentRequestType, reqId, ts)
	ir.Intent = intent
	return ir
}

func (ir *IntentRequest) WithDialogState(sta slu.DialogState) *IntentRequest {
	ir.DialogState = sta
	return ir
}

func (ir *IntentRequest) SetDialogState(sta slu.DialogState) {
	ir.DialogState = sta
}

func (intReq *IntentRequest) MergeIntent(obj *slu.Intent) {
	if obj == nil {
		return
	}
	if intReq.Intent == nil {
		intReq.Intent = slu.NewIntent(obj.Name).WithSubName(obj.SubName)
	}
	intReq.Intent.Merge(obj)
}

func (intReq *IntentRequest) IntentName() string {
	if intReq == nil || intReq.Intent == nil {
		return ""
	}
	return intReq.Intent.Name
}

func (intReq *IntentRequest) GetIntent() *slu.Intent {
	if intReq == nil {
		return nil
	}
	return intReq.Intent
}

func (intReq *IntentRequest) SubIntentName() string {
	if intReq == nil || intReq.Intent == nil {
		return ""
	}
	return intReq.Intent.SubName
}

type SessionEndedRequest struct {
	CoreRequest
	Reason Reason        `json:"reason"`
	Err    *system.Error `json:"error,omitempty"`
}

type IntentsRequest struct {
	CoreRequest
	Intents     []*slu.Intent   `json:"intents"`
	DialogState slu.DialogState `json:"dialogState,omitempty"`
}

func NewIntentsRequest(reqId, ts string, intents []*slu.Intent) *IntentsRequest {
	ir := new(IntentsRequest)
	ir.speechletRequest = new(speechletRequest)
	ir.setBaseInfo(IntentsRequestType, reqId, ts)
	ir.Intents = intents
	return ir
}
