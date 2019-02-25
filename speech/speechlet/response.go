package speechlet

import (
	"encoding/json"
	"log"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/dialog/directives"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/ui"
)

type ResponseEnvelope struct {
	Version string    `json:"version"`
	Status  *Status   `json:"status"`
	Context *Context  `json:"context,omitempty"`
	Results []*Result `json:"results,omitempty"`
}

type ResponseEnvelopeRaw struct {
	Version string          `json:"version"`
	Status  *Status         `json:"status"`
	Context *Context        `json:"context,omitempty"`
	Results json.RawMessage `json:"results,omitempty"`
}

func (raw *ResponseEnvelopeRaw) GetResults() []*Result {
	var resultsRaw []*ResultRaw
	if err := json.Unmarshal([]byte(raw.Results), &resultsRaw); err != nil {
		log.Println(err)
		return nil
	}
	/*var typePre struct {
		Type string `json:"type"`
	}*/
	results := make([]*Result, 0)
	for _, v := range resultsRaw {
		item := new(Result)
		item.Hint = v.Hint
		item.FormatType = v.FormatType
		item.OutputSpeech = v.OutputSpeech
		item.Script = v.Script
		item.Data = v.Data
		/*for _, w := range v.Directives {
			var dir Directive
			if err := json.Unmarshal([]byte(w), &typePre); err != nil {
				log.Println(err)
				return nil
			}
			switch DirectiveType(typePre.Type) {
			case DisplayDirectiveType:
				var dis DisplayDirectiveRaw
				if err := json.Unmarshal([]byte(w), &dis); err != nil {
					log.Println(err)
					return nil
				}
				var disp DisplayDirective
				disp.Type = dis.Type
				disp.Hint = dis.Hint
				disp.Suggestions = dis.Suggestions
				disp.Card = dis.GetCard()
				dir = &disp
			case EventDirectiveType:
				var event EventDirective
				if err := json.Unmarshal([]byte(w), &event); err != nil {
					log.Println(err)
				}
				dir = &event
			default:
				log.Printf("unrecognize directive type: %s", typePre)
				return nil
			}
			item.Directives = append(item.Directives, dir)
		}*/
		results = append(results, item)
	}
	return results
}

func NewResponseEnvelope() *ResponseEnvelope {
	return &ResponseEnvelope{Version: speech.Version}
}

func (respEn *ResponseEnvelope) WithStatus(sta *Status) *ResponseEnvelope {
	respEn.Status = sta
	return respEn
}

func (respEn *ResponseEnvelope) WithContext(ctx *Context) *ResponseEnvelope {
	respEn.Context = ctx
	return respEn
}

func (respEn *ResponseEnvelope) WithResults(results ...*Result) *ResponseEnvelope {
	respEn.Results = results
	return respEn
}

func NewErrResponseEnvelope(detail string) *ResponseEnvelope {
	return NewResponseEnvelope().WithStatus(NewInternalErrStatus(detail))
}

type Response struct {
	Results          []*Result              `json:"results,omitempty"`
	Directives       []directives.Directive `json:"directives,omitempty"`
	ShouldEndSession bool                   `json:"shouldEndSession"`
}

func NewAskResponse(ask string) *Response {
	return makeSpeechResponse(ask, false)
}

func NewTellResponse(answer string) *Response {
	return makeSpeechResponse(answer, true)
}

func NewDelegateResponse(directives []directives.Directive) *Response {
	return &Response{
		Directives:       directives,
		ShouldEndSession: false,
	}
}

func NewResponse() *Response {
	return &Response{}
}

func (resp *Response) WithResults(results ...*Result) *Response {
	resp.Results = results
	return resp
}

func (resp *Response) GetResults() []*Result {
	return resp.Results
}

func (resp *Response) GetFirstResult() *Result {
	if len(resp.Results) == 0 {
		return nil
	}
	return resp.Results[0]
}

func (resp *Response) AppendResults(results ...*Result) *Response {
	resp.Results = append(resp.Results, results...)
	return resp
}

func (resp *Response) WithDerectives(directives []directives.Directive) *Response {
	resp.Directives = directives
	return resp
}

func (resp *Response) WithShouldEndSession(b bool) *Response {
	resp.ShouldEndSession = b
	return resp
}

func (resp *Response) ShouldEnded() bool {
	if resp == nil {
		return true
	}
	return resp.ShouldEndSession
}

func (resp *Response) SetEnded(b bool) {
	if resp == nil {
		resp.ShouldEndSession = b
	}
}

type FormatType string

const (
	CardFormat      FormatType = "CardFormat"
	PlainTextFormat FormatType = "SpeechFormat"
	SsmlFormat      FormatType = "SsmlFormat"

	TextFormat  FormatType = "text"
	AudioFormat FormatType = "audio"
	ListFormat  FormatType = "list"
)

type Timeout struct {
	TimeInMillseconds int    `json:"timeInMs"`
	Action            string `json:"action"`
}

type Result struct {
	FormatType   FormatType   `json:"formatType,omitempty"`
	Hint         string       `json:"hint,omitempty"`
	OutputSpeech *SpeechItems `json:"outputSpeech,omitempty"`
	Script       *ScriptItems `json:"script,omitempty"`
	//Directives   []Directive  `json:"directives,omitempty"`
	Data interface{} `json:"data,omitempty"`

	Timeout  *Timeout   `json:"timeout,omitempty"`
	Emotions []*Emotion `json:"emotions,omitempty"`
}

type ResultRaw struct {
	FormatType   FormatType   `json:"formatType,omitempty"`
	Hint         string       `json:"hint,omitempty"`
	OutputSpeech *SpeechItems `json:"outputSpeech,omitempty"`
	Script       *ScriptItems `json:"script,omitempty"`
	//Directives   []json.RawMessage `json:"directives,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type SpeechItems struct {
	Items []*ui.SpeechItem `json:"items,omitempty"`
}

func NewSpeechItems(items ...*ui.SpeechItem) *SpeechItems {
	return &SpeechItems{Items: items}
}

func (si *SpeechItems) AppendItems(items ...*ui.SpeechItem) *SpeechItems {
	si.Items = append(si.Items, items...)
	return si
}

func (si *SpeechItems) AppendPlainTextItems(texts ...string) *SpeechItems {
	for _, v := range texts {
		si.Items = append(si.Items, ui.NewPlainTextSpeechItem(v))
	}
	return si
}

func (si *SpeechItems) AppendSsmlItems(ssmls ...string) *SpeechItems {
	for _, v := range ssmls {
		si.Items = append(si.Items, ui.NewSsmlSpeechItem(v))
	}
	return si
}

func (si *SpeechItems) AppendAudioItems(urls ...string) *SpeechItems {
	for _, v := range urls {
		si.Items = append(si.Items, ui.NewAudioSpeechItem(v))
	}
	return si
}

type ScriptItems struct {
	Items []*ui.ScriptItem `json:"items,omitempty"`
}

func NewScriptItems(items ...*ui.ScriptItem) *ScriptItems {
	return &ScriptItems{Items: items}
}

func (si *ScriptItems) AppendItems(items ...*ui.ScriptItem) *ScriptItems {
	si.Items = append(si.Items, items...)
	return si
}

func (si *ScriptItems) AppendH5Items(texts ...string) *ScriptItems {
	for _, v := range texts {
		si.Items = append(si.Items, ui.NewH5ScriptItem(v))
	}
	return si
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) WithHint(hint string) *Result {
	r.Hint = hint
	return r
}

func (r *Result) GetHint() string {
	return r.Hint
}

func (r *Result) SetHint(hint string) {
	r.Hint = hint
}

func (r *Result) WithFormatType(ft FormatType) *Result {
	r.FormatType = ft
	return r
}

func (r *Result) WithData(data interface{}) *Result {
	r.Data = data
	return r
}

func (r *Result) WithOutputSpeech(items *SpeechItems) *Result {
	r.OutputSpeech = items
	return r
}

func (r *Result) WithOutputPlainTextSpeech(texts ...string) *Result {
	if r.OutputSpeech == nil {
		r.OutputSpeech = NewSpeechItems()
	}
	r.OutputSpeech.AppendPlainTextItems(texts...)
	return r
}

func (r *Result) GetFirstOutputPlainTextSpeech() (string, bool) {
	speech := r.GetOutputSpeech()
	for _, v := range speech.Items {
		if v.GetType() == ui.PlainTextType {
			return v.GetSource(), true
		}
	}
	return "", false
}

func (r *Result) SetFirstOutputPlainTextSpeech(text string) bool {
	speech := r.GetOutputSpeech()
	for _, v := range speech.Items {
		if v.GetType() == ui.PlainTextType {
			v.SetSource(text)
			return true
		}
	}
	return false
}

func (r *Result) WithOutputSsmlSpeech(ssmls ...string) *Result {
	if r.OutputSpeech == nil {
		r.OutputSpeech = NewSpeechItems()
	}
	r.OutputSpeech.AppendSsmlItems(ssmls...)
	return r
}

func (r *Result) WithOutputAudioSpeech(urls ...string) *Result {
	if r.OutputSpeech == nil {
		r.OutputSpeech = NewSpeechItems()
	}
	r.OutputSpeech.AppendAudioItems(urls...)
	return r
}

/*func (r *Result) GetDirectives() []Directive {
	return r.Directives
}

func (r *Result) WithDisplayDirective(dir *DisplayDirective) *Result {
	for i, v := range r.Directives {
		if v.GetType() == DisplayDirectiveType {
			r.Directives[i] = dir
			return r
		}
	}
	r.Directives = append(r.Directives, dir)
	return r
}

func (r *Result) AppendDisplayDirective(dir *DisplayDirective) *Result {
	r.Directives = append(r.Directives, dir)
	return r
}

func (r *Result) WithEventDirective(dir *EventDirective) *Result {
	for i, v := range r.Directives {
		if v.GetType() == EventDirectiveType {
			r.Directives[i] = dir
			return r
		}
	}
	r.Directives = append(r.Directives, dir)
	return r
}

func (r *Result) GetDisplayDirective() *DisplayDirective {
	for _, v := range r.Directives {
		if v.GetType() == DisplayDirectiveType {
			if w, ok := v.(*DisplayDirective); ok {
				return w
			} else {
				return nil
			}
		}
	}
	return nil
}

func (r *Result) GetEventDirective() *EventDirective {
	for _, v := range r.Directives {
		if v.GetType() == EventDirectiveType {
			if w, ok := v.(*EventDirective); ok {
				return w
			} else {
				return nil
			}
		}
	}
	return nil
}*/

// deprecated
func (r *Result) GetFormatType() FormatType {
	return r.FormatType
}

func (r *Result) GetOutputSpeech() *SpeechItems {
	return r.OutputSpeech
}

func (r *Response) HasDirectives() bool {
	if r == nil {
		return false
	}
	if len(r.Directives) == 0 {
		return false
	}
	return true
}

func (r *Response) GetDirectives() []directives.Directive {
	if !r.HasDirectives() {
		return nil
	}
	return r.Directives
}

func makeSpeechResponse(text string, shouldEnd bool) *Response {
	result := NewResult().WithOutputSpeech(NewSpeechItems(
		ui.NewPlainTextSpeechItem(text)))
	return NewResponse().WithResults(result).WithShouldEndSession(shouldEnd)
}
