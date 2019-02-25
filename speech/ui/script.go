package ui

type ScriptItemInterface interface {
	GetType() string
	GetSource() string
}

type ScriptType string

const (
	H5Type ScriptType = "H5"
)

type ScriptItem struct {
	Type   ScriptType `json:"type"`
	Source string     `json:"source"`
}

func (si *ScriptItem) GetType() ScriptType {
	return si.Type
}

func (si *ScriptItem) GetSource() string {
	return si.Source
}

func (si *ScriptItem) SetSource(text string) {
	si.Source = text
}

func NewScriptItem(typ ScriptType, s string) *ScriptItem {
	return &ScriptItem{Type: typ, Source: s}
}

func NewH5ScriptItem(text string) *ScriptItem {
	return &ScriptItem{Type: H5Type, Source: text}
}
