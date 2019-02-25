package model

import "encoding/json"

type DialogModel struct {
	Dialog  Dialog    `json:"dialog"`
	Prompts []*Prompt `json:"prompts"`
}

func NewDialogModel() *DialogModel {
	return &DialogModel{}
}

func (dm *DialogModel) WithDialog(d Dialog) *DialogModel {
	dm.Dialog = d
	return dm
}

func (dm *DialogModel) WithPrompts(prompts ...*Prompt) *DialogModel {
	dm.Prompts = append(dm.Prompts, prompts...)
	return dm
}

func (dm *DialogModel) GetRandomPrompt(id string) *Prompt {
	if dm == nil {
		return nil
	}
	for _, v := range dm.Prompts {
		if v != nil && v.ID == id {
			return v //v.GetRandomValue()
		}
	}
	return nil
}

func (dm *DialogModel) GetSlotElicit(intentName, slotName string) *Prompt {
	if dm == nil {
		return nil
	}
	slot := dm.GetSlot(intentName, slotName)
	if slot == nil {
		return nil
	}
	return dm.GetRandomPrompt(slot.Prompts.Elicitation)
}

func (dm *DialogModel) GetSlotConfirmation(intentName, slotName string) *Prompt {
	if dm == nil {
		return nil
	}
	slot := dm.GetSlot(intentName, slotName)
	if slot == nil {
		return nil
	}
	return dm.GetRandomPrompt(slot.Prompts.Confirmation)
}

func (dm *DialogModel) GetIntentConfirmation(intentName string) *Prompt {
	if dm == nil {
		return nil
	}
	intent := dm.GetIntent(intentName)
	if intent == nil {
		return nil
	}
	return dm.GetRandomPrompt(intent.Prompts.Confirmation)
}

func (dm *DialogModel) GetIntentResult(intentName string) *Prompt {
	if dm == nil {
		return nil
	}
	intent := dm.GetIntent(intentName)
	if intent == nil {
		return nil
	}
	return dm.GetRandomPrompt(intent.Prompts.Result)
}

func (dm *DialogModel) GetIntent(name string) *Intent {
	if dm == nil {
		return nil
	}
	for _, v := range dm.Dialog.Intents {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (dm *DialogModel) GetSlot(intentName, slotName string) *Slot {
	intent := dm.GetIntent(intentName)
	if intent == nil {
		return nil
	}
	return intent.GetSlot(slotName)
}

func (dm *DialogModel) Verify() bool {
	for _, intent := range dm.Dialog.Intents {
		if intent.NeedConfirm() && dm.GetIntentConfirmation(intent.Name) == nil {
			return false
		}
		if intent.NeedResult() && dm.GetIntentResult(intent.Name) == nil {
			return false
		}
		for _, slot := range intent.Slots {
			if slot.NeedElicit() && dm.GetSlotElicit(intent.Name, slot.Name) == nil {
				return false
			}
			if slot.NeedConfirm() && dm.GetSlotConfirmation(intent.Name, slot.Name) == nil {
				return false
			}
		}
	}
	return true
}

type Dialog struct {
	Intents []*Intent `json:"intents"`
}

func NewDialog(intents ...*Intent) Dialog {
	d := Dialog{}
	d.Intents = append(d.Intents, intents...)
	return d
}

type Intent struct {
	Name                 string    `json:"name"`
	ConfirmationRequired bool      `json:"confirmationRequired"`
	ResultRequired       bool      `json:"resultRequired"`
	Prompts              PromptIds `json:"prompts"`
	Slots                []*Slot   `json:"slots"`
}

func NewIntent(name string, confirmationRequired bool) *Intent {
	return &Intent{Name: name, ConfirmationRequired: confirmationRequired}
}

func (intent *Intent) WithPrompts(ids PromptIds) *Intent {
	intent.Prompts = ids
	return intent
}

func (intent *Intent) WithSlots(slot ...*Slot) *Intent {
	intent.Slots = append(intent.Slots, slot...)
	return intent
}

func (intent *Intent) GetSlot(name string) *Slot {
	for _, v := range intent.Slots {
		if v.Name == name {
			return v
		}
	}
	return nil
}

type Slot struct {
	Name                 string    `json:"name"`
	Type                 string    `json:"type"`
	ConfirmationRequired bool      `json:"confirmationRequired"`
	ElicitationRequired  bool      `json:"elicitationRequired"`
	Prompts              PromptIds `json:"prompts"`
	Handler              string    `json:"handler"`
	ConcealRequired      bool      `json:"concealRequired"`
}

func NewSlot(name, typ string, c, e bool) *Slot {
	return &Slot{
		Name:                 name,
		Type:                 typ,
		ConfirmationRequired: c,
		ElicitationRequired:  e,
	}
}

func (slot *Slot) WithPrompts(ids PromptIds) *Slot {
	slot.Prompts = ids
	return slot
}

type PromptIds struct {
	Confirmation string `json:"confirmation,omitempty"`
	Elicitation  string `json:"elicitation,omitempty"`
	Result       string `json:"result,omitempty"`
}

func NewPromptIds(c, e string) PromptIds {
	return PromptIds{Confirmation: c, Elicitation: e}
}

type Prompt struct {
	ID         string       `json:"id"`
	Variations []*Variation `json:"variations"`
}

/*func NewPrompt(id string, v []*Variation) *Prompt {
	return &Prompt{ID: id, Variations: v}
}*/

/*func (p *Prompt) GetRandomValue() string {
	if p == nil || len(p.Variations) == 0 {
		return ""
	}
	rand.Seed(int64(time.Now().Minute()))
	i := rand.Intn(len(p.Variations))
	if p.Variations[i] == nil {
		return ""
	}
	return p.Variations[i].Value
}*/

type Variation struct {
	Type  string            `json:"type"`
	Value []json.RawMessage `json:"value"`
}

/*func NewVariation(t, v string) *Variation {
	return &Variation{Type: t, Value: v}
}*/

func (intent *Intent) NeedConfirm() bool {
	if intent == nil {
		return false
	}
	return intent.ConfirmationRequired
}

func (intent *Intent) NeedResult() bool {
	if intent == nil {
		return false
	}
	return intent.ResultRequired
}

func (slot *Slot) NeedConfirm() bool {
	if slot == nil {
		return false
	}
	return slot.ConfirmationRequired
}

func (slot *Slot) NeedElicit() bool {
	if slot == nil {
		return false
	}
	return slot.ElicitationRequired
}
