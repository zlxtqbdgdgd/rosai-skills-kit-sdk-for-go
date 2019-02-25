package speechlet

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"reflect"

	"strings"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/dialog/directives"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/dialog/model"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/ui"
)

type RequestHandler struct {
	AppId       string
	Speechlet   Speechlet
	DialogModel *model.DialogModel

	RequestVerifiers  []RequestVerifier
	ResponseVerifiers []ResponseVerifier

	SlotHandler         reflect.Value
	DialogModelCallback DialogModelCallback
}

type DialogModelCallback interface {
	GetDialogModel(ctx *Context) *model.DialogModel
}

func (rh *RequestHandler) HandleCall(reqBytes []byte) ([]byte, error) {
	start := time.Now()
	reqEn, err := makeRequestEnvelope(reqBytes)
	if err != nil {
		log.Printf("ERROR] reqBytes: %s, error: %s", string(reqBytes), err)
		return nil, err
	}
	defer func() {
		eclipse := float64(time.Since(start).Nanoseconds()) / 1e6
		log.Printf("Request[%s] total cost: %6.3f ms", reqEn.Request.GetRequestId(), eclipse)
	}()
	// verify request
	for _, v := range rh.RequestVerifiers {
		if !v.Verify(reqEn) {
			eString := fmt.Sprintf("Could not validate Request %s using verifier %T,"+
				" rejiecting request", reqEn.Request.GetRequestId(), v)
			log.Println(eString)
			return nil, errors.New(eString)
		}
	}
	// dispatch and handle request to get response
	resp, ctx, err := rh.dispatchCall(reqEn)
	// verify response
	//session := req.Session
	//for _, v := range rh.responseVerifiers {
	//	if v.Verify(resp, session) {
	//		eString := fmt.Sprintf("Could not validate Response %s using verifier %T,"+
	//			" rejiecting response", req.Request.GetRequestId(), v)
	//		log.Println(eString)
	//		return nil, errors.New(eString)
	//	}
	//}
	var status *Status
	if err == nil {
		status = NewGoodStatus()
	} else {
		if err == ErrServiceMismatched {
			status = NewMismatchStatus(err.Error())
		} else {
			status = NewInternalErrStatus(err.Error())
		}
		log.Printf("Warning] Request: %s, error: %s", reqEn.Request.GetRequestId(), err)
	}
	// make results
	var results []*Result
	if resp != nil {
		results = resp.Results
	}
	// make RequestEnvelope
	respEn := NewResponseEnvelope().WithStatus(status).WithContext(ctx).WithResults(results...)
	// serialize response
	respBytes, err := json.MarshalIndent(respEn, "", "  ")
	if err != nil {
		return nil, err
	}
	log.Printf("Request[%s] response << %s", reqEn.Request.GetRequestId(), string(respBytes))
	return respBytes, nil
}

func (rh *RequestHandler) dispatchCall(reqEn *RequestEnvelope) (
	resp *Response, ctx *Context, err error) {
	if reqEn == nil || reqEn.Context == nil {
		return nil, nil, errors.New("RequestEnvelope or it's Context is nil")
	}
	userId, appId := reqEn.Context.GetUserId(), reqEn.Context.GetAppId()
	deviceId, skillId := reqEn.Context.GetDeviceId(), reqEn.Context.GetSkillId()
	if appId == "" || deviceId == "" || skillId == "" {
		return nil, nil, errors.New(fmt.Sprintf("AppId[%s], DeviceId[%s], SkillId[%s]"+
			" not allowed empty", appId, deviceId, skillId))
	}
	session, err := FetchSessionFromHistory(userId, appId, deviceId, skillId)
	if err != nil {
		return nil, nil, errors.New("fetch session failed: " + err.Error())
	}
	/*if rh.DialogModel == nil &&reqEn.Request.GetType() != IntentsRequestType && rh.PrivateDialogModel == nil {
		return nil, nil, errors.New(fmt.Sprintf("Request[%s] DialogModel is nil",
			reqEn.Request.GetRequestId()))
	}*/
	var dm *model.DialogModel = rh.DialogModel
	if dm == nil {
		dm = rh.DialogModelCallback.GetDialogModel(reqEn.Context)
	}
	if dm == nil {
		return nil, nil, errors.New(fmt.Sprintf("Request[%s] DialogModel is nil",
			reqEn.Request.GetRequestId()))
	}
	// debug log
	ssBytes, err := json.MarshalIndent(session, "", "  ")
	log.Printf("Fetch request[%s] session: %s", reqEn.Request.GetRequestId(), string(ssBytes))
	// If this is a new session, invoke the speechlet's onSessionStarted life-cycle method.
	if session.New {
		err = rh.Speechlet.OnSessionStarted(reqEn)
		if err != nil {
			return nil, nil, err
		}
	}
	//
	switch reqEn.Request.GetType() {
	default:
		log.Printf("Warning] unkown request type: %s", reqEn.Request.GetType())
	case SessionEndedRequestType:
		err = rh.Speechlet.OnSessionEnded(reqEn)
	case LaunchRequestType:
		resp, err = rh.Speechlet.OnLaunch(reqEn)
	case IntentRequestType:
		resp, ctx, err = rh.handleIntentRequest(reqEn, session, dm)
	case IntentsRequestType:
		// this stage, only call OnIntent, slots info are handled in dst
		resp, ctx, err = rh.Speechlet.OnIntent(reqEn)
	}
	return resp, ctx, err
}

func (rh *RequestHandler) handleIntentRequest(reqEn *RequestEnvelope,
	session *Session, dm *model.DialogModel) (resp *Response, ctx *Context, err error) {
	// pre handle request
	req, ok := reqEn.Request.(*IntentRequest)
	if !ok {
		return nil, nil, errors.New(fmt.Sprintf("assert request[%+v] to "+
			"IntentRequest failed, type: %T", reqEn.Request, reqEn.Request))
	}
	var ask string
	if ask, err = rh.preHandleIntentRequest(req, session, dm); err != nil {
		return nil, nil, err
	}
	if ask != "" {
		return NewAskResponse(ask), nil, nil
	}
	reqEn.Request = req
	// debug log
	bytes, _ := json.MarshalIndent(reqEn, "", "  ")
	log.Printf("OnIntent RequestEnvelope[%s]: %s", req.GetRequestId(), string(bytes))
	// OnIntent
	resp, ctx, err = rh.Speechlet.OnIntent(reqEn)
	if resp == nil || err != nil {
		return nil, nil, err
	}
	//if ctx == nil {
	//	ctx = NewContext().WithSystem(reqEn.Context.GetSystem())
	//}
	// debug log
	bytes, _ = json.MarshalIndent(resp, "", "  ")
	log.Printf("response[%s] from bots: %s", req.GetRequestId(), string(bytes))
	//
	if resp.HasDirectives() {
		log.Printf("INFO] Request[%s] response has directives", req.GetRequestId())
		resp, err = rh.handleDirectiveResponse(req, resp, dm)
		session.WithUpdatedIntent(req.Intent)
	} /* else {
		if resp.ShouldEnded() {
			session.ClearAllIntents()
		} else {
			session.MergeIntent(req.Intent)
		}
	}*/

	// try to resolve response
	resolveResponse(req.Intent, resp)

	if resp.ShouldEnded() {
		session.ClearAllIntents()
	} else {
		session.MergeIntent(req.Intent)
	}

	// push session to cache between multiply servers
	if err := PushSessionToCache(session); err != nil {
		log.Printf("Warning] PushSessionToCache[%s] error: %s", req.GetRequestId(), err)
	}
	ssBytes, _ := json.MarshalIndent(session, "", "  ")
	log.Printf("push request[%s] session to cache: %s", req.GetRequestId(), string(ssBytes))
	// share slots information to context
	ctx = rh.shareSlotsToContext(req.Intent, ctx, dm)
	ctx.ClearSystemInfo()
	//ssBytes, _ = json.MarshalIndent(ctx, "", "  ")
	//log.Printf("[%s] context to be shared: %s", req.GetRequestId(), string(ssBytes))
	return resp, ctx, err
}

func resolveResponse(intent *slu.Intent, resp *Response) {
	if resp == nil || len(resp.Results) == 0 {
		return
	}
	var paramMap map[string]string = make(map[string]string)
	for _, v := range intent.Slots {
		paramMap[v.Name] = v.GetStringOrgin()

	}
	for _, v := range resp.Results {
		_tryResolveParams(&v.Hint, paramMap)
		if v.OutputSpeech != nil {
			for _, speechItem := range v.OutputSpeech.Items {
				if speechItem.Type == ui.PlainTextType {
					_tryResolveParams(&speechItem.Source, paramMap)
				}
			}
		}
	}
}

func _tryResolveParams(unresolved *string, paramMap map[string]string) bool {
	if *unresolved == "" {
		return true
	}

	var resolved string
	var isParam bool = false
	var lastIndex int = 0
	unresolvedToken := strings.Split(*unresolved, "")
	var unresolvedLength = len(unresolvedToken)
	for i := 0; i < unresolvedLength; i++ {
		//log.Printf("======================>", unresolvedToken[i], i, unresolvedLength,
		//	(unresolvedToken[i] == "{")/*, unresolvedToken[i+1] == "$"*/)
		if unresolvedToken[i] == "{" && i+1 < unresolvedLength &&
			unresolvedToken[i+1] == "$" {
			// 支持嵌套 {%x{%y%}%}
			resolved += strings.Join(unresolvedToken[lastIndex:i], "")
			//log.Printf("————————————————————————————————————>>>>>>>>>>>>>>>>", resolved,
			//	lastIndex, i-lastIndex)
			lastIndex = i + 2
			isParam = true
			i++
		}
		if unresolvedToken[i] == "}" /* && i+1 < unresolvedLength &&unresolvedToken[i+1] == "}"*/ {
			if !isParam {
				log.Printf("Cannot resolve params in string, invalid format. %s", *unresolved)
				return false
			}
			var paramName = strings.Join(unresolvedToken[lastIndex:i], "")
			//log.Printf("==========================>>>>>>>>>>>>>>>>", paramName,
			//	lastIndex, i-lastIndex)
			if v, ok := paramMap[paramName]; !ok {
				log.Printf("Cannot resolve params in string, unknow param. %s %s",
					*unresolved, paramName)
				return false
			} else {
				resolved += v
			}
			lastIndex = i + 1
			isParam = false
			//i++
		}
	}

	if isParam {
		log.Printf("Cannot resolve params in string, invalid format. %s", *unresolved)
		return false
	}
	if lastIndex < unresolvedLength {
		resolved += strings.Join(unresolvedToken[lastIndex:], "")
	}

	*unresolved = resolved

	data, _ := json.Marshal(paramMap)
	log.Printf("tryResolveParams:[unresolved=%s][resolved=%s][%s]", *unresolved,
		resolved, string(data))

	return true
}

func (rh *RequestHandler) shareSlotsToContext(intent *slu.Intent, ctx *Context,
	dm *model.DialogModel) *Context {
	if ctx == nil {
		ctx = NewContext()
	}
	mi := dm.GetIntent(intent.Name)
	for _, v := range intent.Slots {
		//log.Printf("==========>", v.Name, mi.GetSlot(v.Name).ConcealRequired)

		modslot := mi.GetSlot(v.Name)
		if modslot != nil && modslot.ConcealRequired {
			continue
		}

		if v.GetValue() != nil {
			ctx.WithParameter(v.Name, v.Value)
		}
	}
	return ctx
}

func (rh *RequestHandler) preHandleIntentRequest(req *IntentRequest,
	session *Session, dm *model.DialogModel) (string, error) {
	// make request with full slots
	// 1. new a intent from dialog model
	intentName := req.IntentName()
	intent := slu.NewIntentFromModel(dm, intentName)
	if intent == nil {
		return "", errors.New("NewIntentFromModel failed, intent name mismatched")
	}
	intent.WithSubName(req.SubIntentName())
	// 2. fetch history slot values from session
	intent.Merge(session.GetUpdatedIntent(intentName))
	if intent.Started() {
		req.DialogState = slu.STARTED
	} else {
		req.DialogState = slu.IN_PROGRESS
	}
	// debug log
	//bytes, _ := json.MarshalIndent(intent, "", "  ")
	//log.Printf("intent merged session history: %s", string(bytes))
	// 3. update intent by using slot values in request
	intent.Merge(req.Intent)
	// debug log
	//bytes, _ = json.MarshalIndent(intent, "", "  ")
	//log.Printf("intent merged request: %s", string(bytes))
	// 4. clean slots without value
	mi := dm.GetIntent(intent.Name)
	intent.CleanSlots(mi)

	for _, reqslot := range req.Intent.Slots {
		modslot := mi.GetSlot(reqslot.Name)
		if modslot == nil || modslot.Handler == "" {
			continue
		}
		params := make([]reflect.Value, 2)
		params[0] = reflect.ValueOf(reqslot)
		params[1] = reflect.ValueOf(intent)
		value := rh.SlotHandler.MethodByName(modslot.Handler)
		if !value.IsValid() {
			return "", errors.New("Call SlotHandler failed, handler name mismatched")
		}
		values := value.Call(params)
		if len(values) == 0 {
			return "", errors.New("Call SlotHandler failed, invalid return")
		}
		ask := values[0].String()
		if ask != "" {
			return ask, nil
		}
	}

	// 5. make requestEnvelope with new intent
	req.Intent = intent

	// 6. set intent request status
	if intent.Completed(mi) {
		req.DialogState = slu.COMPLETED
	}
	return "", nil
}

func (rh *RequestHandler) handleDirectiveResponse(req *IntentRequest,
	resp *Response, dm *model.DialogModel) (*Response, error) {
	var err error
	for _, v := range resp.GetDirectives() {
		switch v.GetType() {
		default:
			return nil, errors.New(fmt.Sprintf("response Directive[%+v] type not found", v))
		case directives.DelegateType:
			updatedIntent := v.GetUpdatedIntent()
			req.Intent.Merge(updatedIntent)
			resp, err = rh.handleDelegateDirective(req.Intent, dm)
		}
	}
	return resp, err
}

func (rh *RequestHandler) handleDelegateDirective(intent *slu.Intent,
	dm *model.DialogModel) (*Response, error) {
	var result *Result = nil
	mi := dm.GetIntent(intent.Name)
	for _, v := range mi.Slots {
		if v.NeedElicit() && intent.CanElicit(v.Name) {
			result = makeResultFromPrompt(dm.GetSlotElicit(intent.Name, v.Name))
			break
		}
		if v.NeedConfirm() && intent.CanConfirm(v.Name) {
			result = makeResultFromPrompt(dm.GetSlotConfirmation(intent.Name, v.Name))
			break
		}
	}
	if result == nil {
		if mi.NeedResult() {
			result = makeResultFromPrompt(dm.GetIntentResult(intent.Name))
			resp := NewResponse().WithResults(result).WithShouldEndSession(true)
			return resp, nil
		} else {
			//return nil, errors.New("delegate speechText is empty")
			// intended left blank.
		}

	}
	resp := NewResponse().WithResults(result).WithShouldEndSession(false)
	return resp, nil
}

func parseRequestType(raw interface{}) (string, error) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return "", errors.New(fmt.Sprintf("request assert[%+v] map failed, type: %T",
			m, m))
	}
	v, ok := m["type"]
	if !ok {
		return "", errors.New(fmt.Sprintf("request[%+v] no type field", m))
	}
	t, ok := v.(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("request assert[%+v] type string failed", v))
	}
	return t, nil
}

func makeRequestEnvelope(reqBytes []byte) (*RequestEnvelope, error) {
	// deserialize request
	var initRE struct {
		Version string      `json:"version"`
		Session *Session    `json:"session,omitempty"`
		Context *Context    `json:"context"`
		Request interface{} `json:"request"`
	}
	if err := json.Unmarshal(reqBytes, &initRE); err != nil {
		return nil, err
	}
	typ, err := parseRequestType(initRE.Request)
	if err != nil {
		return nil, err
	}
	var req Request
	bytes, _ := json.Marshal(initRE.Request)
	switch RequestType(typ) {
	default:
		return nil, errors.New("request type not found")
	case SessionStartedRequestType:
		var r SessionStartedRequest
		if err = json.Unmarshal(bytes, &r); err == nil {
			req = &r
		}
	case SessionEndedRequestType:
		var r SessionEndedRequest
		if err = json.Unmarshal(bytes, &r); err == nil {
			req = &r
		}
	case LaunchRequestType:
		var r LaunchRequest
		if err = json.Unmarshal(bytes, &r); err == nil {
			req = &r
		}
	case IntentRequestType:
		var r IntentRequest
		if err = json.Unmarshal(bytes, &r); err == nil {
			req = &r
		}
	case IntentsRequestType:
		var r IntentsRequest
		if err = json.Unmarshal(bytes, &r); err == nil {
			req = &r
		}
	}
	if err != nil {
		return nil, err
	}
	reqEn := RequestEnvelope{
		Version: initRE.Version,
		//Session: initRE.Session,
		Context: initRE.Context,
		Request: req,
	}
	return &reqEn, nil
}

func makeResultFromPrompt(prompt *model.Prompt) *Result {
	if prompt == nil {
		return nil
	}

	result := NewResult()
	var firstText bool = false
	for _, v := range prompt.Variations {
		rand.Seed(int64(time.Now().Second()))
		idx := rand.Intn(len(v.Value))
		selectedValue := v.Value[idx]

		if ui.SpeechType(v.Type) == ui.PlainTextType {
			var plainText string
			if err := json.Unmarshal([]byte(selectedValue), &plainText); err != nil {
				log.Printf("DisplayDirectiveRaw %+v GetCard error: %s", selectedValue, err)
				continue
			}

			if !firstText {
				result.WithHint(plainText)
				firstText = true
			}

			result.WithOutputPlainTextSpeech(plainText)
		}
		if ui.SpeechType(v.Type) == ui.AudioType {
			var audio string
			if err := json.Unmarshal([]byte(selectedValue), &audio); err != nil {
				log.Printf("DisplayDirectiveRaw %+v GetCard error: %s", selectedValue, err)
				continue
			}

			result.WithOutputAudioSpeech(audio)
		}
		/*if DirectiveType(v.Type) == DisplayDirectiveType {
			var raw DisplayDirectiveRaw
			if err := json.Unmarshal([]byte(selectedValue), &raw); err != nil {
				log.Printf("DisplayDirectiveRaw %+v GetCard error: %s", selectedValue, err)
				continue
			}

			result.WithDisplayDirective(NewDisplayDirective().
				WithCard(raw.GetCard()))
		}*/
	}
	return result
}
