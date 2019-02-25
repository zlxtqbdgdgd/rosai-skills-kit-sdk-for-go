package ui

type SpeechItemInterface interface {
	GetType() string
	GetSource() string
}

type SpeechType string

const (
	PlainTextType SpeechType = "PlainText"
	SSMLType      SpeechType = "SSML"
	AudioType     SpeechType = "Audio"
)

type SpeechItem struct {
	Type   SpeechType `json:"type"`
	Source string     `json:"source"`
}

func (si *SpeechItem) GetType() SpeechType {
	return si.Type
}

func (si *SpeechItem) GetSource() string {
	return si.Source
}

func (si *SpeechItem) SetSource(text string) {
	si.Source = text
}

func NewSpeechItem(typ SpeechType, s string) *SpeechItem {
	return &SpeechItem{Type: typ, Source: s}
}

func NewPlainTextSpeechItem(text string) *SpeechItem {
	return &SpeechItem{Type: PlainTextType, Source: text}
}

func NewSsmlSpeechItem(ssml string) *SpeechItem {
	return &SpeechItem{Type: SSMLType, Source: ssml}
}

func NewAudioSpeechItem(url string) *SpeechItem {
	return &SpeechItem{Type: AudioType, Source: url}
}
