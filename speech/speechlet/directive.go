package speechlet

/*
type DirectiveType string

const (
	DisplayDirectiveType DirectiveType = "Display.Customized"
	EventDirectiveType   DirectiveType = "ROSAI.EVENT"
)

type Directive interface {
	GetType() DirectiveType
}

type DisplayDirective struct {
	Type        DirectiveType    `json:"type"`
	Hint        string           `json:"hint,omitempty"`
	Card        ui.CardInterface `json:"card,omitempty"`
	Suggestions []string         `json:"suggestions,omitempty"`
}

type DisplayDirectiveRaw struct {
	Type        DirectiveType   `json:"type"`
	Hint        string          `json:"hint,omitempty"`
	Card        json.RawMessage `json:"card,omitempty"`
	Suggestions []string        `json:"suggestions,omitempty"`
}

func (ddr *DisplayDirectiveRaw) GetCard() ui.CardInterface {
	var typePre struct {
		Type string `json:"type"`
	}
	var err error
	if err = json.Unmarshal([]byte(ddr.Card), &typePre); err != nil {
		log.Printf("DisplayDirectiveRaw %+v GetCard error: %s", ddr, err)
		return nil
	}
	switch ui.CardType(typePre.Type) {
	case ui.ImagesCardType:
		var ic ui.ImagesCard
		if err = json.Unmarshal([]byte(ddr.Card), &ic); err == nil {
			return &ic
		}
	case ui.StandardCardType:
		var sc ui.StandardCard
		if err = json.Unmarshal([]byte(ddr.Card), &sc); err == nil {
			return &sc
		}
	case ui.TextCardType:
		var tc ui.TextCard
		if err = json.Unmarshal([]byte(ddr.Card), &tc); err == nil {
			return &tc
		}
	case ui.ListCardType:
		var lc ui.ListCard
		if err = json.Unmarshal([]byte(ddr.Card), &lc); err == nil {
			return &lc
		}
	default:
		log.Printf("DisplayDirectiveRaw %+v unrecgnized card type: %s", ddr, typePre.Type)
	}
	if err != nil {
		log.Printf("DisplayDirectiveRaw %+v GetCard error: %s", ddr, err)
	}
	return nil
}

func NewDisplayDirective() *DisplayDirective {
	return &DisplayDirective{Type: DisplayDirectiveType}
}

func (dd *DisplayDirective) GetType() DirectiveType {
	return dd.Type
}

func (dd *DisplayDirective) WithHint(hint string) *DisplayDirective {
	dd.Hint = hint
	return dd
}

func (dd *DisplayDirective) WithCard(card ui.CardInterface) *DisplayDirective {
	dd.Card = card
	return dd
}

func (dd *DisplayDirective) WithSuggestions(ss ...string) *DisplayDirective {
	dd.Suggestions = append(dd.Suggestions, ss...)
	return dd
}

func (dd *DisplayDirective) GetHint() string {
	return dd.Hint
}

func (dd *DisplayDirective) GetCard() ui.CardInterface {
	return dd.Card
}

func (dd *DisplayDirective) GetSuggestions() []string {
	return dd.Suggestions
}

type EventDirective struct {
	Type  DirectiveType `json:"type"`
	Event *Event        `json:"event"`
}

type Event struct {
	Name   string `json:"name"`
	Period int    `json:"period"`
}

func NewEvent(name string, period int) *Event {
	return &Event{Name: name, Period: period}
}

func NewEventDirective(name string, period int) *EventDirective {
	return &EventDirective{
		Type:  EventDirectiveType,
		Event: NewEvent(name, period),
	}
}

func (ed *EventDirective) GetType() DirectiveType {
	return ed.Type
}

func (ed *EventDirective) GetEventName() string {
	if ed.Event == nil {
		return ""
	}
	return ed.Event.Name
}

func (ed *EventDirective) GetEventPeriod() int {
	if ed.Event == nil {
		return 0
	}
	return ed.Event.Period
}
*/
