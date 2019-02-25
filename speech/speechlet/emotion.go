package speechlet

type EmotionCode string

const (
	Peaceful EmotionCode = "A001"

	Happy             EmotionCode = "B001"
	Exciting          EmotionCode = "B002"
	ExtremelyExciting EmotionCode = "B003"

	Dissatisfy     EmotionCode = "C001"
	Displeasure    EmotionCode = "C002"
	Angry          EmotionCode = "C003"
	Rage           EmotionCode = "C004"
	ExtremelyAngry EmotionCode = "C005"

	Aggrieved      EmotionCode = "D001"
	Depressed      EmotionCode = "D002"
	Disappointment EmotionCode = "D003"
	Guilt          EmotionCode = "D004"
	Sadness        EmotionCode = "D005"
	Pain           EmotionCode = "D006"

	Surprised EmotionCode = "E001"
	Shocked   EmotionCode = "E002"

	Fear EmotionCode = "F001"

	Envy  EmotionCode = "G001"
	Hate  EmotionCode = "G002"
	Abhor EmotionCode = "G003"
	Blame EmotionCode = "G004"

	Irritability EmotionCode = "H001"

	Vigilant EmotionCode = "I001"
	Anxious  EmotionCode = "I002"
	Flurried EmotionCode = "I003"

	Reassuring EmotionCode = "J001"
	Satisfied  EmotionCode = "J002"
	Relax      EmotionCode = "J003"

	Respect EmotionCode = "K001"
	Adore   EmotionCode = "K002"
	Praise  EmotionCode = "K003"
	Believe EmotionCode = "K004"
	Love    EmotionCode = "K005"

	Shy EmotionCode = "L001"

	Miss  EmotionCode = "M001"
	Await EmotionCode = "M002"
	Yearn EmotionCode = "M003"

	Bless EmotionCode = "N001"

	Curious EmotionCode = "O001"
	Doubt   EmotionCode = "O002"

	Boring EmotionCode = "P001"
	Dull   EmotionCode = "P002"
	Sleepy EmotionCode = "P003"
	Tired  EmotionCode = "P004"
)

type LevelCode int

const (
	LPeace      LevelCode = 0 // A
	LHappy      LevelCode = 1 // B
	LAngry      LevelCode = 2
	LSad        LevelCode = 3
	LSuprise    LevelCode = 4
	LFear       LevelCode = 5
	LEnvy       LevelCode = 6
	LAnnoy      LevelCode = 7
	LAnxious    LevelCode = 8
	LReassuring LevelCode = 9
	LRespect    LevelCode = 10
	LShy        LevelCode = 11
	LMiss       LevelCode = 12
	LBless      LevelCode = 13
	LDoubt      LevelCode = 14
	LBoring     LevelCode = 15
)

func getLevelCodeByEmotionCode(code EmotionCode) LevelCode {
	return LevelCode(code[0] - 'A')
}

type Emotion struct {
	Type  string      `json:"type"`
	Level LevelCode   `json:"level"`
	Code  EmotionCode `json:"code"`
}

func NewEmotionWithEmotionCode(code EmotionCode) *Emotion {
	return &Emotion{
		Level: getLevelCodeByEmotionCode(code),
		Code:  code,
	}
}

func NewDefaultEmotion() *Emotion {
	return &Emotion{
		Type:  "answer",
		Level: 0,
		Code:  "A001",
	}
}

func (e *Emotion) WithType(typ string) *Emotion {
	e.Type = typ
	return e
}

var EmotionCodeLst = []EmotionCode{
	Peaceful,
	Happy,
	Exciting,
	ExtremelyExciting,
	Dissatisfy,
	Displeasure,
	Angry,
	Rage,
	ExtremelyAngry,
	Aggrieved,
	Depressed,
	Disappointment,
	Guilt,
	Sadness,
	Pain,
	Surprised,
	Shocked,
	Fear,
	Envy,
	Hate,
	Abhor,
	Blame,
	Irritability,
	Vigilant,
	Anxious,
	Flurried,
	Reassuring,
	Satisfied,
	Relax,
	Respect,
	Adore,
	Praise,
	Believe,
	Love,
	Shy,
	Miss,
	Await,
	Yearn,
	Bless,
	Curious,
	Doubt,
	Boring,
	Dull,
	Sleepy,
	Tired}

func Conv2EmotionCode(code int) EmotionCode {
	if code > len(EmotionCodeLst)-1 {
		return Peaceful
	}

	return EmotionCodeLst[code]
}
