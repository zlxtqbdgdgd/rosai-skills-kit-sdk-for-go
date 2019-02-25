package ui

type CardType string

const (
	TextCardType     CardType = "Text"
	StandardCardType CardType = "Standard"
	ImagesCardType   CardType = "Images"
	ListCardType     CardType = "List"
	TimerCardType    CardType = "Timer"
)

type CardInterface interface {
	GetType() CardType
}

type TextCard struct {
	Type    CardType `json:"type"`
	Title   string   `json:"title"`
	Content string   `json:"content,omitempty"`
}

func NewTextCard(title, content string) *TextCard {
	return &TextCard{Type: TextCardType, Title: title, Content: content}
}

func (tc *TextCard) GetType() CardType {
	return tc.Type
}

type StandardCard struct {
	TextCard
	Image *Image `json:"image"`
}

func (sc *StandardCard) GetType() CardType {
	return sc.Type
}

func NewStandardCard(title, content string) *StandardCard {
	sc := &StandardCard{TextCard: *NewTextCard(title, content)}
	sc.Type = StandardCardType
	return sc
}

func (sc *StandardCard) WithContent(content string) *StandardCard {
	sc.Content = content
	return sc
}

func (sc *StandardCard) GetTitle() string {
	return sc.Title
}

func (sc *StandardCard) GetContent() string {
	return sc.Content
}

func (sc *StandardCard) WithImage(url string) *StandardCard {
	sc.Image = NewImage(url)
	return sc
}

func (sc *StandardCard) WithBulletImage(url string, bs *BulletScreen) *StandardCard {
	sc.Image = NewImage(url).WithBullet(bs)
	return sc
}

type ImagesCard struct {
	Type CardType `json:"type"`
	List []*Image `json:"list"`
}

func NewImagesCard(list ...*Image) *ImagesCard {
	return &ImagesCard{Type: ImagesCardType, List: list}
}

func (ic *ImagesCard) GetType() CardType {
	return ic.Type
}

func (ic *ImagesCard) AppendImages(list ...*Image) *ImagesCard {
	ic.List = append(ic.List, list...)
	return ic
}

type ListCard struct {
	Type CardType        `json:"type"`
	List []*StandardCard `json:"list"`
}

func NewListCard(list ...*StandardCard) *ListCard {
	return &ListCard{Type: ListCardType, List: list}
}

func (lc *ListCard) GetType() CardType {
	return lc.Type
}

func (lc *ListCard) AppendCards(list ...*StandardCard) *ListCard {
	lc.List = append(lc.List, list...)
	return lc
}

type BulletPositionType string

const (
	Top         BulletPositionType = "top"
	Bottom      BulletPositionType = "bottom"
	Left        BulletPositionType = "left"
	Right       BulletPositionType = "right"
	Center      BulletPositionType = "center"
	Full        BulletPositionType = "full"
	TopLeft     BulletPositionType = "top-left"
	TopRight    BulletPositionType = "top-right"
	BottomLeft  BulletPositionType = "bottom-left"
	BottomRight BulletPositionType = "bottom-right"
)

type BulletScreen struct {
	URL    string `json:"imageUrl,omitempty"`
	Text   string `json:"text,omitempty"`
	Period int    `json:"period,omitempty"`

	TextPosition  BulletPositionType `json:"textPosition,omitempty"`
	ImagePosition BulletPositionType `json:"imagePosition,omitempty"`
}

func NewBulletScreen() *BulletScreen {
	return &BulletScreen{}
}

func (bs *BulletScreen) WithImage(url string, pos BulletPositionType) *BulletScreen {
	bs.URL = url
	bs.ImagePosition = pos
	return bs
}

func (bs *BulletScreen) WithText(text string, pos BulletPositionType) *BulletScreen {
	bs.Text = text
	bs.TextPosition = pos
	return bs
}

func (bs *BulletScreen) WithPeriod(period int) *BulletScreen {
	bs.Period = period
	return bs
}

type Image struct {
	URL          string        `json:"url"`
	BulletScreen *BulletScreen `json:"bulletScreen,omitempty"`
}

func NewImage(url string) *Image {
	return &Image{URL: url}
}

func (i *Image) WithBullet(bs *BulletScreen) *Image {
	i.BulletScreen = bs
	return i
}
