package slu

import (
	"encoding/json"
	"errors"
	"fmt"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/dialog/model"
	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu/entityresolution"
)

var ErrNilSlot = errors.New("slot is nil")
var ErrNilIntent = errors.New("intent is nil")

// An Intent is the output of spoken language understanding (SLU) that represents
// what the user wants to do based on a predefined schema definition.

type DialogState string

const (
	// When a skill has a managed dialog configured, this field indicates the current dialog state
	// for the intent request.

	// Indicates this is the first turn in a multi-turn dialog. Skills can use this state to
	// trigger behavior that only needs to be executed in the first turn. For example, a skill
	// may wish to provide missing slot values that it can determine based on the current user,
	// or other such session information, but doesn't wish to perform that action for every turn
	// in a dialog.
	STARTED DialogState = "STARTED"

	// Indicates that a multi-turn dialog is in process and it is not the first turn. Skills may
	// assume that all of the required slot values and confirmations have not yet been provided,
	// and react accordingly (for instance by immediately returning a response containing a
	// DelegateDirective).
	IN_PROGRESS DialogState = "IN_PROGRESS"

	// Indicates that all required slot values and confirmations have been provided, the dialog
	// is considered complete, and the skill can proceed to fulfilling the intent. Nevertheless,
	// the skill may manually continue the dialog if it determines at runtime that it requires
	// more input in order to fulfill the intent, in which case it may return an appropriate
	// .DialogDirective and update slot values/confirmations as required.
	COMPLETED DialogState = "COMPLETED"
)

// ConfirmationStatus Indication of whether an intent or slot has been explicitly
// confirmed or denied by the user, or neither.
// Intents can be confirmed or denied using dialog.directives.ConfirmIntentDirective,
// or by indicating in the skill's configured dialog that the intent requires confirmation.
// Slots can be confirmed or denied using dialog.directives.ConfirmSlotDirective,
// or by indicating in the skill's configured dialog that the intent requires confirmation.
type ConfirmationStatus string

const (
	NONE      ConfirmationStatus = "NONE"
	CONFIRMED ConfirmationStatus = "CONFIRMED"
	DENIED    ConfirmationStatus = "DENIED"
)

type ValueType string

const (
	NullType   ValueType = "Null"
	IntType    ValueType = "Int"
	BoolType   ValueType = "Bool"
	FloatType  ValueType = "Float"
	StringType ValueType = "String"

	ArrayType    ValueType = "Array"
	StrArrayType ValueType = "StrArray"
	MapType      ValueType = "Map"
	MapArrayType ValueType = "MapArray"
)

var (
	SupportedValueTypes = []ValueType{
		IntType,
		BoolType,
		FloatType,
		StringType,
		ArrayType,
		StrArrayType,
		MapType,
		MapArrayType,
	}
)

var (
	NilErr   = errors.New("value is nil")
	TypeErr  = errors.New("type mismatch")
	ValueErr = errors.New("Value internal error")
)

type Value struct {
	Origin    interface{} `json:"orgin"`
	ValueType ValueType   `json:"normType"`
	Value     interface{} `json:"norm"`
	Logic     string      `json:"logic,omitempty"`
	Priority  interface{} `json:"priority,omitempty"`
	Tag       interface{} `json:"tag,omitempty"`
}

func NewValue(typ ValueType, v interface{}) *Value {
	return &Value{ValueType: typ, Value: v}
}

func NewIntValue(i int) *Value {
	return NewValue(IntType, i)
}

func NewBoolValue(b bool) *Value {
	return NewValue(BoolType, b)
}

func NewFloatValue(f float64) *Value {
	return NewValue(FloatType, f)
}

func NewStringValue(s string) *Value {
	return NewValue(StringType, s)
}

func NewArrayValue(a []interface{}) *Value {
	return NewValue(ArrayType, a)
}

func NewStrArrayValue(a []string) *Value {
	return NewValue(StrArrayType, a)
}

func NewMapValue(m map[string]interface{}) *Value {
	return NewValue(MapType, m)
}

func NewStrMapValue(m string) *Value {
	return NewValue(MapType, m)
}

func NewMapArrayValue(ma []map[string]interface{}) *Value {
	return NewValue(MapArrayType, ma)
}

func NewStrMapArrayValue(ma string) *Value {
	return NewValue(MapArrayType, ma)
}

func (v *Value) WithOrigin(o interface{}) *Value {
	v.Origin = o
	return v
}

func (v *Value) GetOrigin() interface{} {
	return v.Origin
}

func (v *Value) SetOrigin(o interface{}) {
	v.Origin = o
}

func (v *Value) WithLogic(l string) *Value {
	v.Logic = l
	return v
}

func (v *Value) GetLogic() string {
	return v.Logic
}

func (v *Value) SetLogic(l string) {
	v.Logic = l
}

func (v *Value) WithPriority(p interface{}) *Value {
	v.Priority = p
	return v
}

func (v *Value) GetPriority() interface{} {
	return v.Priority
}

func (v *Value) SetPriority(p interface{}) {
	v.Priority = p
}

func (v *Value) WithTag(t interface{}) *Value {
	v.Tag = t
	return v
}

func (v *Value) GetTag() interface{} {
	return v.Tag
}

func (v *Value) SetTag(t interface{}) {
	v.Tag = t
}

func (v *Value) SetValueType(t ValueType) {
	v.ValueType = t
}

func (v *Value) SetValue(val interface{}) {
	v.Value = val
}

func (v *Value) GetType() ValueType {
	if v == nil {
		return NullType
	}
	return v.ValueType
}

func (v *Value) HasValue() bool {
	if v == nil {
		return false
	}
	for _, t := range SupportedValueTypes {
		if t == v.ValueType {
			return (v.Value != nil)
		}
	}
	return false
}

func (v *Value) GetValue() interface{} {
	return v.Value
}

func (v *Value) GetIntValue() (int, error) {
	switch v.GetType() {
	default:
		return -1, TypeErr
	case IntType:
		if vv, ok := v.Value.(int); ok {
			return vv, nil
		} else if vv, ok := v.Value.(float64); ok {
			return int(vv), nil
		} else {
			return -1, ValueErr
		}
	}
}

func (v *Value) GetBoolValue() (bool, error) {
	switch v.GetType() {
	default:
		return false, TypeErr
	case BoolType:
		if vv, ok := v.Value.(bool); ok {
			return vv, nil
		} else {
			return false, ValueErr
		}
	}
}

func (v *Value) GetFloatValue() (float64, error) {
	switch v.GetType() {
	default:
		return 0.0, TypeErr
	case FloatType:
		if vv, ok := v.Value.(float64); ok {
			return vv, nil
		} else {
			return 0.0, ValueErr
		}
	}
}

func (v *Value) GetStringValue() (string, error) {
	switch v.GetType() {
	default:
		return "", TypeErr
	case StringType:
		if vv, ok := v.Value.(string); ok {
			return vv, nil
		} else {
			return "", ValueErr
		}
	}
}

func (v *Value) GetStringOrgin() (string, error) {
	switch v.GetType() {
	default:
		return "", TypeErr
	case StringType:
		if vv, ok := v.Origin.(string); ok {
			return vv, nil
		} else {
			return "", ValueErr
		}
	}
}

func (v *Value) GetArrayValue() ([]interface{}, error) {
	switch v.GetType() {
	default:
		return nil, TypeErr
	case ArrayType:
		if vv, ok := v.Value.([]interface{}); ok {
			return vv, nil
		} else {
			return nil, ValueErr
		}
	}
}

func (v *Value) GetStrArrayValue() ([]string, error) {
	switch v.GetType() {
	default:
		return nil, TypeErr
	case StrArrayType:
		if vv, ok := v.Value.([]string); ok {
			return vv, nil
		} else if vv, ok := v.Value.([]interface{}); ok {
			return Conv2StrSlice(vv), nil
		} else {
			return nil, ValueErr
		}
	}
}

func (v *Value) GetMapValue() (map[string]interface{}, error) {
	if v == nil {
		return nil, NilErr
	}
	if v.GetType() == MapType {
		if vv, ok := v.Value.(map[string]interface{}); ok {
			return vv, nil
		} else if s, ok := v.Value.(string); ok {
			// NOTE: to adjust to qu string map value
			var m map[string]interface{}
			if err := json.Unmarshal([]byte(s), &m); err != nil {
				return nil, ValueErr
			}
			return m, nil
		}
	}
	return nil, ValueErr
}

func (v *Value) GetMapArrayValue() ([]map[string]interface{}, error) {
	if v.GetType() == MapArrayType {
		if vv, ok := v.Value.([]map[string]interface{}); ok {
			return vv, nil
		} else if s, ok := v.Value.(string); ok {
			// NOTE: to adjust to qu string map value
			var ma []map[string]interface{}
			if err := json.Unmarshal([]byte(s), &ma); err != nil {
				return nil, ValueErr
			}
			return ma, nil
		}
	}
	return nil, ValueErr
}

type Slot struct {
	Name               string                        `json:"name"`
	Value              *Value                        `json:"value,omitempty"`
	ConfirmationStatus ConfirmationStatus            `json:"confirmationStatus"`
	Resolutions        *entityresolution.Resolutions `json:"resolutions,omitempty"`
}

func (slot *Slot) CanConfirm() bool {
	if slot == nil {
		return false
	}
	if slot.Value.HasValue() && slot.ConfirmationStatus == NONE {
		return true
	}
	return false
}

func (slot *Slot) CanElicit() bool {
	if slot == nil {
		return false
	}
	if slot.Value.HasValue() {
		return false
	}
	return true
}

func (slot *Slot) HasValue() bool {
	if slot == nil {
		return false
	}
	return slot.Value.HasValue()
}

func NewSlot(name string) *Slot {
	return &Slot{Name: name, ConfirmationStatus: NONE}
}

func (slot *Slot) WithValue(value *Value) *Slot {
	slot.Value = value
	return slot
}

func (slot *Slot) WithIntValue(i int) *Slot {
	slot.Value = NewIntValue(i)
	return slot
}

func (slot *Slot) WithBoolValue(b bool) *Slot {
	slot.Value = NewBoolValue(b)
	return slot
}

func (slot *Slot) WithFloatValue(f float64) *Slot {
	slot.Value = NewFloatValue(f)
	return slot
}

func (slot *Slot) WithStringValue(s string) *Slot {
	slot.Value = NewStringValue(s)
	return slot
}

func (slot *Slot) WithArrayValue(a []interface{}) *Slot {
	slot.Value = NewArrayValue(a)
	return slot
}

func (slot *Slot) WithStrArrayValue(a []string) *Slot {
	slot.Value = NewStrArrayValue(a)
	return slot
}

func (slot *Slot) WithMapValue(m map[string]interface{}) *Slot {
	slot.Value = NewMapValue(m)
	return slot
}

func (slot *Slot) WithMapArrayValue(ma []map[string]interface{}) *Slot {
	slot.Value = NewMapArrayValue(ma)
	return slot
}

func (slot *Slot) WithStatus(sta ConfirmationStatus) *Slot {
	if slot == nil {
		return nil
	}
	slot.ConfirmationStatus = sta
	return slot
}

func (slot *Slot) GetValue() *Value {
	if slot == nil {
		return nil
	}
	return slot.Value
}

func (slot *Slot) GetIntValue() int {
	v, err := slot.GetValue().GetIntValue()
	if err != nil {
		//log.Printf("slot %+v get int value error: %s", slot, err)
		return 0
	}
	return v
}

func (slot *Slot) GetBoolValue() bool {
	v, err := slot.GetValue().GetBoolValue()
	if err != nil {
		//log.Printf("slot %+v get bool value error: %s", slot, err)
		return false
	}
	return v
}

func (slot *Slot) GetFloatValue() float64 {
	v, err := slot.GetValue().GetFloatValue()
	if err != nil {
		//log.Printf("slot %+v get float value error: %s", slot, err)
		return 0.0
	}
	return v
}

func (slot *Slot) GetStringValue() string {
	v, err := slot.GetValue().GetStringValue()
	if err != nil {
		//log.Printf("slot %+v get string value error: %s", slot, err)
		return ""
	}
	return v
}

func (slot *Slot) GetStringOrgin() string {
	v, err := slot.GetValue().GetStringOrgin()
	if err != nil {
		//log.Printf("slot %+v get string value error: %s", slot, err)
		return ""
	}
	return v
}

func (slot *Slot) GetArrayValue() []interface{} {
	v, err := slot.GetValue().GetArrayValue()
	if err != nil {
		//log.Printf("slot %+v get array value error: %s", slot, err)
		return nil
	}
	return v
}

func (slot *Slot) GetStrArrayValue() []string {
	v, err := slot.GetValue().GetStrArrayValue()
	if err != nil {
		//log.Printf("slot %+v get string array value error: %s", slot, err)
		return nil
	}
	return v
}

func (slot *Slot) GetMapValue() map[string]interface{} {
	v, err := slot.GetValue().GetMapValue()
	if err != nil {
		//log.Printf("slot %+v get map value error: %s", slot, err)
		return nil
	}
	return v
}

func (slot *Slot) GetMapArrayValue() []map[string]interface{} {
	v, err := slot.GetValue().GetMapArrayValue()
	if err != nil {
		//log.Printf("slot %+v get map array value error: %s", slot, err)
		return nil
	}
	return v
}

func (slot *Slot) SetValue(value *Value) error {
	if slot == nil {
		return ErrNilSlot
	}
	slot.Value = value
	return nil
}

func (slot *Slot) GetStatus() ConfirmationStatus {
	if slot == nil {
		return NONE
	}
	return slot.ConfirmationStatus
}

func (slot *Slot) SetStatus(sta ConfirmationStatus) error {
	if slot == nil {
		return ErrNilSlot
	}
	slot.ConfirmationStatus = sta
	return nil
}

// An Intent is the output of spoken language understanding (SLU) that represents
// what the user wants to do based on a predefined schema definition.
type Intent struct {
	Name               string             `json:"name"`
	SubName            string             `json:"subName,omitempty"`
	ConfirmationStatus ConfirmationStatus `json:"confirmationStatus,omitempty"`
	Slots              map[string]*Slot   `json:"slots"`
}

func (intent *Intent) Started() bool {
	if intent == nil {
		return false
	}
	for _, v := range intent.Slots {
		if v.HasValue() {
			return false
		}
	}
	return true
}

func (intent *Intent) Completed(mi *model.Intent) bool {
	if intent == nil || mi == nil {
		return false
	}
	if !intent.SlotsConfirmed(mi) {
		return false
	}
	if mi.NeedConfirm() {
		return intent.ConfirmationStatus == NONE
	} else {
		return true
	}
	/*if mi.NeedConfirm() && intent.ConfirmationStatus != NONE {
		return false
	}
	if mi.NeedResult() {
		return false
	}
	return true*/
}

func (intent *Intent) SlotsConfirmed(mi *model.Intent) bool {
	//bytes, _ := json.MarshalIndent(mi, "", "  ")
	//log.Printf("model.Intent: %s", string(bytes))
	for _, v := range mi.Slots {
		if v.NeedElicit() && intent.CanElicit(v.Name) {
			//log.Printf("%s elicit need: %t, can: %t", v.Name, v.NeedElicit(), intent.CanElicit(v.Name))
			return false
		}
		if v.NeedConfirm() && intent.CanConfirm(v.Name) {
			//log.Printf("%s confrim need: %t, can: %t", v.Name, v.NeedConfirm(), intent.CanConfirm(v.Name))
			return false
		}
	}
	return true
}

func (intent *Intent) CanElicit(slotName string) bool {
	if intent == nil {
		return false
	}
	if slot, ok := intent.Slots[slotName]; ok {
		return slot.CanElicit()
	}
	return false
}

func (intent *Intent) CanConfirm(slotName string) bool {
	if intent == nil {
		return false
	}
	if slot, ok := intent.Slots[slotName]; ok {
		return slot.CanConfirm()
	}
	return false
}

func (intent *Intent) Merge(obj *Intent) bool {
	if obj == nil || intent == nil {
		return false
	}
	if intent.Name != obj.Name {
		return false
	}
	if len(obj.Slots) == 0 {
		return false
	}
	if len(intent.Slots) == 0 {
		intent.Slots = make(map[string]*Slot)
	}
	for k, v := range obj.Slots {
		if v.Value.HasValue() {
			intent.Slots[k] = v
		}
	}
	if intent.ConfirmationStatus == NONE {
		intent.ConfirmationStatus = obj.ConfirmationStatus
	}
	return true
}

func (intent *Intent) CleanSlots(mi *model.Intent) {
	for k, v := range intent.Slots {
		if !v.HasValue() && !mi.GetSlot(v.Name).NeedElicit() {
			delete(intent.Slots, k)
		}
	}
}

func NewIntent(name string) *Intent {
	return &Intent{
		Name:               name,
		Slots:              make(map[string]*Slot),
		ConfirmationStatus: NONE,
	}
}

func NewIntentFromModel(dm *model.DialogModel, name string) *Intent {
	if dm == nil || name == "" {
		return nil
	}
	for _, v := range dm.Dialog.Intents {
		if v.Name == name {
			intent := NewIntent(name)
			for _, vv := range v.Slots {
				slot := NewSlot(vv.Name)
				intent.WithSlot(slot)
			}
			return intent
		}
	}
	return nil
}

func (intent *Intent) WithSubName(name string) *Intent {
	intent.SubName = name
	return intent
}

func (intent *Intent) WithSlot(slot *Slot) *Intent {
	if intent == nil {
		return nil
	}
	if slot == nil || slot.Name == "" {
		return intent
	}
	intent.Slots[slot.Name] = slot
	return intent
}

func (intent *Intent) WithStatus(sta ConfirmationStatus) *Intent {
	if intent == nil {
		return nil
	}
	intent.ConfirmationStatus = sta
	return intent
}

func (intent *Intent) SetSlot(slot *Slot) {
	intent.WithSlot(slot)
}

func (intent *Intent) GetSlot(name string) *Slot {
	if intent == nil || intent.Slots == nil {
		return nil
	}
	if slot, ok := intent.Slots[name]; ok {
		return slot
	} else {
		return nil
	}
}

func (intent *Intent) GetStatus(name string) ConfirmationStatus {
	if intent == nil {
		return NONE
	}
	return intent.ConfirmationStatus
}

func (intent *Intent) SetStatus(sta ConfirmationStatus) error {
	if intent == nil {
		return ErrNilIntent
	}
	intent.ConfirmationStatus = sta
	return nil
}

func Conv2StrSlice(si []interface{}) []string {
	var ss []string
	for _, v := range si {
		ss = append(ss, fmt.Sprint(v))
	}
	return ss
}
