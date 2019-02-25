package speechlet

import (
	"encoding/json"
	"log"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/interfaces"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
)

type CtxParamInternalKey string

const (
	CtxParamKeyHitIntent       CtxParamInternalKey = "_subIntentHit"
	CtxParamKeyHitSkill        CtxParamInternalKey = "_skillHit"
	CtxParamKeyLastInputSlots  CtxParamInternalKey = "_lastInputSlots"
	CtxParamKeyLastOutputSlots CtxParamInternalKey = "_lastOutputSlots"

	NullCtxWord = "__null_context_word__"
)

type CtxParams map[string]*slu.Value

func NewCtxParams() CtxParams {
	return make(CtxParams)
}

func (pa CtxParams) WithParameter(k string, v *slu.Value) CtxParams {
	return pa.SetParameter(k, v)
}

func (pa CtxParams) GetParameter(k string) *slu.Value {
	if len(pa) == 0 {
		return nil
	}
	if v, ok := pa[k]; ok {
		return v
	} else {
		return nil
	}
}

func (pa CtxParams) DelParameter(k string) {
	delete(pa, k)
}

func (pa CtxParams) SetParameter(k string, v *slu.Value) CtxParams {
	if pa == nil {
		pa = make(CtxParams)
	}
	pa[k] = v
	return pa
}

func (pa CtxParams) SetIntValue(k string, v int) CtxParams {
	return pa.SetParameter(k, slu.NewIntValue(v))
}

func (pa CtxParams) GetIntValue(k string) int {
	v, e := pa.GetParameter(k).GetIntValue()
	if e != nil {
		return 0
	}
	return v
}

func (pa CtxParams) SetBoolValue(k string, v bool) CtxParams {
	return pa.SetParameter(k, slu.NewBoolValue(v))
}

func (pa CtxParams) GetBoolValue(k string) bool {
	v, e := pa.GetParameter(k).GetBoolValue()
	if e != nil {
		return false
	}
	return v
}

func (pa CtxParams) SetFloatValue(k string, v float64) CtxParams {
	return pa.SetParameter(k, slu.NewFloatValue(v))
}

func (pa CtxParams) GetFloatValue(k string) float64 {
	f, e := pa.GetParameter(k).GetFloatValue()
	if e != nil {
		return 0.0
	}
	return f
}

func (pa CtxParams) SetStrArrayValue(k string, v []string) CtxParams {
	return pa.SetParameter(k, slu.NewStrArrayValue(v))
}

func (pa CtxParams) GetStrArrayValue(k string) []string {
	s, e := pa.GetParameter(k).GetStrArrayValue()
	if e != nil {
		return nil
	}
	return s
}

func (pa CtxParams) SetStringValue(k, v string) CtxParams {
	return pa.SetParameter(k, slu.NewStringValue(v))
}

func (pa CtxParams) GetStringValue(k string) string {
	s, e := pa.GetParameter(k).GetStringValue()
	if e != nil {
		return ""
	}
	return s
}

func (pa CtxParams) SetMapValue(k string, m map[string]interface{}) CtxParams {
	return pa.SetParameter(k, slu.NewMapValue(m))
}

func (pa CtxParams) GetMapValue(k string) map[string]interface{} {
	m, e := pa.GetParameter(k).GetMapValue()
	if e != nil {
		return nil
	}
	return m
}

func (pa CtxParams) SetMapArrayValue(k string, ma []map[string]interface{}) CtxParams {
	return pa.SetParameter(k, slu.NewMapArrayValue(ma))
}

func (pa CtxParams) GetMapArrayValue(k string) []map[string]interface{} {
	ma, e := pa.GetParameter(k).GetMapArrayValue()
	if e != nil {
		return nil
	}
	return ma
}

type Context struct {
	System *System `json:"system,omitempty"`

	Context      string `json:"context,omitempty"`
	LifespanInMs int64  `json:"lifespanInMs,omitempty"`
	CtxParams    `json:"parameters,omitempty"`
}

func NewContext() *Context {
	return &Context{CtxParams: make(CtxParams)}
}

func (ctx *Context) WithSystem(sys *System) *Context {
	ctx.System = sys
	return ctx
}

func (ctx *Context) ClearSystemInfo() {
	ctx.System = nil
}

func (ctx *Context) GetSystem() *System {
	return ctx.System
}

type System struct {
	Skill     *Skill  `json:"skill,omitempty"`
	User      *User   `json:"user,omitempty"`
	Device    *Device `json:"device,omitempty"`
	CtxParams `json:"parameters,omitempty"`
}

func (sys *System) GetParameters() CtxParams {
	return sys.CtxParams
}

func (sys *System) SetParameters(params CtxParams) {
	sys.CtxParams = params
}

func (sys *System) WithParameters(params CtxParams) *System {
	sys.SetParameters(params)
	return sys
}

func (ctx *Context) GetSysParameters() CtxParams {
	return ctx.System.CtxParams
}

func (ctx *Context) GetSysParameter(k string) *slu.Value {
	return ctx.System.GetParameter(k)
}

func (ctx *Context) SetSysParameter(k string, v *slu.Value) {
	ctx.System.CtxParams = ctx.System.SetParameter(k, v)
}

func (ctx *Context) DelSysParameter(k string) {
	delete(ctx.System.CtxParams, k)
}

func (ctx *Context) GetSysIntValue(k string) int {
	return ctx.System.GetIntValue(k)
}

func (ctx *Context) GetSysBoolValue(k string) bool {
	return ctx.System.GetBoolValue(k)
}

func (ctx *Context) GetSysFloatValue(k string) float64 {
	return ctx.System.GetFloatValue(k)
}

func (ctx *Context) GetSysStringValue(k string) string {
	return ctx.System.GetStringValue(k)
}

func (ctx *Context) GetSysMapValue(k string) map[string]interface{} {
	return ctx.System.GetMapValue(k)
}

func (ctx *Context) GetSysMapArrayValue(k string) []map[string]interface{} {
	return ctx.System.GetMapArrayValue(k)
}

func NewCtxSystem() *System {
	return &System{}
}

func (sys *System) WithUser(user *User) *System {
	sys.User = user
	return sys
}

func (sys *System) WithDevice(device *Device) *System {
	sys.Device = device
	return sys
}

func (sys *System) WithSkill(skill *Skill) *System {
	sys.Skill = skill
	return sys
}

type Device struct {
	DeviceId            string                          `json:"deviceId"`
	SupportedInterfaces map[string]interfaces.Interface `json:"supportedInterfaces,omitempty"`
}

func NewDevice(id string) *Device {
	return &Device{DeviceId: id}
}

func (dev *Device) WithSupportedInterfaces(m map[string]interfaces.Interface) *Device {
	dev.SupportedInterfaces = m
	return dev
}

func (ctx *Context) GetUserId() string {
	if ctx == nil || ctx.System == nil || ctx.System.User == nil {
		return ""
	}
	return ctx.System.User.UserId
}

func (ctx *Context) GetUserAccessToken() string {
	if ctx == nil || ctx.System == nil || ctx.System.User == nil {
		return ""
	}
	return ctx.System.User.AccessToken
}

func (ctx *Context) GetAppId() string {
	if ctx == nil || ctx.System == nil || ctx.System.User == nil {
		return ""
	}
	return ctx.System.User.AppId
}

func (ctx *Context) GetSkillId() string {
	if ctx == nil || ctx.System == nil || ctx.System.Skill == nil {
		return ""
	}
	return ctx.System.Skill.SkillId
}

func (ctx *Context) GetDeviceId() string {
	if ctx == nil || ctx.System == nil || ctx.System.Device == nil {
		return ""
	}
	return ctx.System.Device.DeviceId
}

func (ctx *Context) WithContext(c string) *Context {
	ctx.Context = c
	return ctx
}

func (ctx *Context) GetContext() string {
	return ctx.Context
}

func (ctx *Context) SetContext(c string) {
	ctx.Context = c
}

func (ctx *Context) WithLifespanInMs(ms int) *Context {
	ctx.LifespanInMs = int64(ms)
	return ctx
}

func (ctx *Context) GetLifespanInMs() int {
	return int(ctx.LifespanInMs)
}

func (ctx *Context) SetLifespanInMs(ms int) {
	ctx.LifespanInMs = int64(ms)
}

func (ctx *Context) WithParameters(params CtxParams) *Context {
	ctx.CtxParams = params
	return ctx
}

func (ctx *Context) GetParameters() CtxParams {
	return ctx.CtxParams
}

func (ctx *Context) SetParameters(params CtxParams) {
	ctx.CtxParams = params
}

func (ctx *Context) WithParameter(k string, v *slu.Value) *Context {
	ctx.CtxParams = ctx.CtxParams.SetParameter(k, v)
	return ctx
}

func (ctx *Context) SetParameter(k string, v *slu.Value) {
	ctx.CtxParams = ctx.CtxParams.SetParameter(k, v)
}

func (ctx *Context) SetIntValue(k string, v int) {
	ctx.CtxParams = ctx.CtxParams.SetIntValue(k, v)
}

func (ctx *Context) WithIntValue(k string, v int) *Context {
	ctx.CtxParams = ctx.CtxParams.SetIntValue(k, v)
	return ctx
}

func (ctx *Context) SetBoolValue(k string, v bool) {
	ctx.CtxParams = ctx.CtxParams.SetBoolValue(k, v)
}

func (ctx *Context) WithBoolValue(k string, v bool) *Context {
	ctx.CtxParams = ctx.CtxParams.SetBoolValue(k, v)
	return ctx
}

func (ctx *Context) SetFloatValue(k string, v float64) {
	ctx.CtxParams = ctx.CtxParams.SetFloatValue(k, v)
}

func (ctx *Context) WithFloatValue(k string, v float64) *Context {
	ctx.CtxParams = ctx.CtxParams.SetFloatValue(k, v)
	return ctx
}

func (ctx *Context) SetStringValue(k, v string) {
	ctx.CtxParams = ctx.CtxParams.SetStringValue(k, v)
}

func (ctx *Context) WithStringValue(k, v string) *Context {
	ctx.CtxParams = ctx.CtxParams.SetStringValue(k, v)
	return ctx
}

func (ctx *Context) SetStrArrayValue(k string, v []string) {
	ctx.CtxParams = ctx.CtxParams.SetStrArrayValue(k, v)
}

func (ctx *Context) WithStrArrayValue(k string, v []string) *Context {
	ctx.CtxParams = ctx.CtxParams.SetStrArrayValue(k, v)
	return ctx
}

func (ctx *Context) SetMapValue(k string, v map[string]interface{}) {
	ctx.CtxParams = ctx.CtxParams.SetMapValue(k, v)
}

func (ctx *Context) WithMapValue(k string, v map[string]interface{}) *Context {
	ctx.CtxParams = ctx.CtxParams.SetMapValue(k, v)
	return ctx
}

func (ctx *Context) SetMapArrayValue(k string, ma []map[string]interface{}) {
	ctx.CtxParams = ctx.CtxParams.SetMapArrayValue(k, ma)
}

func (ctx *Context) WithMapArrayValue(k string, ma []map[string]interface{}) *Context {
	ctx.CtxParams = ctx.CtxParams.SetMapArrayValue(k, ma)
	return ctx
}

func (ctx *Context) SetInternalIntValue(k CtxParamInternalKey, v int) {
	ctx.SetIntValue(string(k), v)
}

func (ctx *Context) GetInternalIntValue(k CtxParamInternalKey) int {
	return ctx.GetIntValue(string(k))
}

func (ctx *Context) SetInternalBoolValue(k CtxParamInternalKey, v bool) {
	ctx.SetBoolValue(string(k), v)
}

func (ctx *Context) GetInternalBoolValue(k CtxParamInternalKey) bool {
	return ctx.GetBoolValue(string(k))
}

func (ctx *Context) SetInternalFloatValue(k CtxParamInternalKey, v float64) {
	ctx.SetFloatValue(string(k), v)
}

func (ctx *Context) GetInternalFloatValue(k CtxParamInternalKey) float64 {
	return ctx.GetFloatValue(string(k))
}

func (ctx *Context) SetInternalStringValue(k CtxParamInternalKey, v string) {
	ctx.SetStringValue(string(k), v)
}

func (ctx *Context) GetInternalStringValue(k CtxParamInternalKey) string {
	return ctx.GetStringValue(string(k))
}

func (ctx *Context) SetInternalStrArrayValue(k CtxParamInternalKey, v []string) {
	ctx.SetStrArrayValue(string(k), v)
}

func (ctx *Context) GetInternalStrArrayValue(k CtxParamInternalKey) []string {
	return ctx.GetStrArrayValue(string(k))
}

func (ctx *Context) SetInternalMapValue(k CtxParamInternalKey, v map[string]interface{}) {
	ctx.SetMapValue(string(k), v)
}

func (ctx *Context) GetInternalMapValue(k CtxParamInternalKey) map[string]interface{} {
	return ctx.GetMapValue(string(k))
}

func (ctx *Context) DelInternalParameter(k CtxParamInternalKey) {
	delete(ctx.CtxParams, string(k))
}

func (ctx *Context) SetSubIntentHit(name string) {
	ctx.SetStringValue(string(CtxParamKeyHitIntent), name)
}

func (ctx *Context) GetLastSubIntentHit() string {
	s := ctx.GetStringValue(string(CtxParamKeyHitIntent))
	// NOTE: to be compatible with old definition temporaryly, remove it after cc update
	if s == "" {
		return ctx.GetStringValue(".subIntentHit")
	}
	return s
}

func (ctx *Context) SetSkillHit(skill string) {
	ctx.SetStringValue(string(CtxParamKeyHitSkill), skill)
}

func (ctx *Context) GetLastSkillHit() string {
	s := ctx.GetStringValue(string(CtxParamKeyHitSkill))
	// NOTE: to be compatible with old definition temporaryly, remove it after cc update
	if s == "" {
		return ctx.GetStringValue(".skillHit")
	}
	return s
}

func (ctx *Context) SetLastInputSlots(ss ...*slu.Slot) {
	slots := ctx.GetMapValue(string(CtxParamKeyLastInputSlots))
	if slots == nil {
		slots = make(map[string]interface{})
	}
	for _, v := range ss {
		slots[v.Name] = v
	}
	ctx.SetMapValue(string(CtxParamKeyLastInputSlots), slots)
}

func (ctx *Context) GetLastInputSlots() map[string]*slu.Slot {
	m := ctx.GetMapValue(string(CtxParamKeyLastInputSlots))
	// NOTE: to be compatible with old definition temporaryly, remove it after cc update
	if m == nil {
		m = ctx.GetMapValue(".lastInputSlots")
	}
	ms := make(map[string]*slu.Slot)
	for k, v := range m {
		s := new(slu.Slot)
		raw, _ := json.Marshal(v)
		if err := json.Unmarshal(raw, s); err != nil {
			log.Printf("error: %s when GetLastOutputSlots", err)
			continue
		}
		ms[k] = s
	}
	return ms
}

func (ctx *Context) SetLastOutputSlots(ss ...*slu.Slot) {
	slots := ctx.GetMapValue(string(CtxParamKeyLastOutputSlots))
	if slots == nil {
		slots = make(map[string]interface{})
	}
	for _, v := range ss {
		slots[v.Name] = v
	}
	ctx.SetMapValue(string(CtxParamKeyLastOutputSlots), slots)
}

func (ctx *Context) GetLastOutputSlots() map[string]*slu.Slot {
	m := ctx.GetMapValue(string(CtxParamKeyLastOutputSlots))
	// NOTE: to be compatible with old definition temporaryly, remove it after cc update
	if m == nil {
		m = ctx.GetMapValue(".lastOutputSlots")
	}
	ms := make(map[string]*slu.Slot)
	for k, v := range m {
		s := new(slu.Slot)
		raw, _ := json.Marshal(v)
		if err := json.Unmarshal(raw, s); err != nil {
			log.Printf("error: %s when GetLastOutputSlots", err)
			continue
		}
		ms[k] = s
	}
	return ms
}
