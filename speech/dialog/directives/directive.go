package directives

import "roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"

type Directive interface {
	GetType() (typ Type)
	GetUpdatedIntent() *slu.Intent
}

func NewConfirmSlotDirective(slotName string, intent *slu.Intent) Directive {
	return &ConfirmSlotDirective{
		dialogDirective: dialogDirective{
			Type:          "Dialog.ConfirmSlot",
			UpdatedIntent: intent,
		},
		SlotToConfirm: slotName,
	}
}

func NewConfirmIntentDirective(intent *slu.Intent) Directive {
	return &ConfirmIntentDirective{
		dialogDirective: dialogDirective{
			Type:          "Dialog.ConfirmIntent",
			UpdatedIntent: intent,
		},
	}
}

func NewDelegateDirective(intent *slu.Intent) Directive {
	return &DelegateDirective{
		dialogDirective: dialogDirective{
			Type:          "Dialog.Delegate",
			UpdatedIntent: intent,
		},
	}
}

func NewElicitSlotDirective(slotName string, intent *slu.Intent) Directive {
	return &ElicitSlotDirective{
		dialogDirective: dialogDirective{
			Type:          "Dialog.ElicitSlot",
			UpdatedIntent: intent,
		},
		SlotToElicit: slotName,
	}
}

type Type string

const (
	DelegateType      Type = "Dialog.Delegate"
	ElicitSlotType    Type = "Dialog.ElicitSlot"
	ConfirmSlotType   Type = "Dialog.ConfirmSlot"
	ConfirmIntentType Type = "Dialog.ConfirmIntent"
)

// A Directive which a skill may return to indicate that the skill is asking the user to
// confirm or deny a slot value. The skill must also provide output speech for the request.
// If the user confirms the slot value, subsequent requests to the skill for the same dialog
// session will have a confirmationStatus of ConfirmationStatus#CONFIRMED for that slot;
// if the user denies the value, the confirmationStatus will be ConfirmationStatus#DENIED.
type dialogDirective struct {
	// type can be one in "Dialog.ConfirmSlot", , ,
	Type Type `json:"type"`

	// Provide an updated intent object to use in subsequent turns of the dialog. All slot values
	// and confirmation provided by the updated intent will replace the existing values and
	// confirmations; if no slot values or confirmations are provided, then they will all be
	// unset. May be left unset or set to null to indicate that that there are no updates to the
	// existing slot values or confirmations.
	UpdatedIntent *slu.Intent `json:"updatedIntent,omitempty"`
}

func (d *dialogDirective) GetType() Type {
	return d.Type
}

func (d *dialogDirective) GetUpdatedIntent() *slu.Intent {
	return d.UpdatedIntent
}

// A Directive which a skill may return to indicate that the skill is asking the user to
// confirm or deny a slot value. The skill must also provide output speech for the request.
// If the user confirms the slot value, subsequent requests to the skill for the same dialog
// session will have a confirmationStatus of ConfirmationStatus#CONFIRMED} for that slot;
// if the user denies the value, the confirmationStatus will be ConfirmationStatus#DENIED.
type ConfirmSlotDirective struct {
	dialogDirective
	SlotToConfirm string `json:"slotToConfirm"`
}

// A Directive which a skill may return to indicate that the skill is asking the user to
// confirm or deny the overall intent. The skill must also provide output speech for the request.
// If the user confirms the intent, subsequent requests to the skill for the same dialog
// session will have a confirmationStatus of ConfirmationStatus#CONFIRMED for the intent;
// if the user denies the value, the confirmationStatus will be ConfirmationStatus#DENIED}.
//
// When a user confirms the intent, it is expected that the skill will then try to fulfill the
// intent. When a user denies the intent, or the user confirms the intent but the skill is unable
// to fulfill it, the skill may either end the session or give the user the opportunity
// to (re-)confirm or change one or more slot values via a sequence of and/or ElicitSlotDirective.
// In this case, the skill should set UpdatedIntent to clear any values or confirmation statuses
// as necessary.
type ConfirmIntentDirective struct {
	dialogDirective
}

// A Directive which a skill may return to indicate that dialog management should be delegated
// to rosai, based on the required slots and confirmations configured in the Skill Builder.
type DelegateDirective struct {
	dialogDirective
}

// A Directive which a skill may return to indicate that the skill is requesting the user to
// provide a value for a slot. The skill must also provide output speech for the request.
type ElicitSlotDirective struct {
	dialogDirective
	SlotToElicit string `json:"slotToElicit"`
}
