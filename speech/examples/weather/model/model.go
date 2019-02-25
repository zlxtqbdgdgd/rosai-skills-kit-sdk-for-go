package model

type Result struct {
	Index          int    `json:"index"`
	Pm25           string `json:"pm25"`
	City           string `json:"city"`
	Focus          string `json:"focus"`
	Weather        string `json:"weather"`
	Temperature    string `json:"temperature"`
	MinTemp        string `json:"minTemp"`
	MaxTemp        string `json:"maxTemp"`
	Date           string `json:"date"`
	Humidity       string `json:"humidity"`
	WindDir        string `json:"windDir"`
	WindLevel      string `json:"windLevel"`
	WindDay        string `json:"windDay"`
	WindDayLevel   string `json:"windDayLevel"`
	WindNight      string `json:"windNight"`
	WindNightLevel string `json:"windNightLevel"`
	Alter          string `json:"alter"`
	Led            string `json:"led,omitempty"`
	TypeTag        string `json:"typeTag,omitempty"`
}

type Results []*Result

func NewResults() Results {
	return []*Result{}
}

func (results Results) Append(r ...*Result) Results {
	results = append(results, r...)
	return results
}

func NewResult(index int) *Result {
	return &Result{Index: index}
}

func (r *Result) WithPm25(pm25 string) *Result {
	r.Pm25 = pm25
	return r
}

func (r *Result) WithCity(city string) *Result {
	r.City = city
	return r
}

func (r *Result) WithFocus(focus string) *Result {
	r.Focus = focus
	return r
}

func (r *Result) WithWeather(weather string) *Result {
	r.Weather = weather
	return r
}

func (r *Result) WithTemperature(temperature string) *Result {
	r.Temperature = temperature
	return r
}

func (r *Result) WithMinTemp(minTemp string) *Result {
	r.MinTemp = minTemp
	return r
}

func (r *Result) WithMaxTemp(maxTemp string) *Result {
	r.MaxTemp = maxTemp
	return r
}

func (r *Result) WithDate(date string) *Result {
	r.Date = date
	return r
}

func (r *Result) WithHumidity(humidity string) *Result {
	r.Humidity = humidity
	return r
}

func (r *Result) WithWindDir(windDir string) *Result {
	r.WindDir = windDir
	return r
}

func (r *Result) WithWindLevel(windLevel string) *Result {
	r.WindLevel = windLevel
	return r
}

func (r *Result) WithWindDay(windDay string) *Result {
	r.WindDay = windDay
	return r
}

func (r *Result) WithWindDayLevel(windDayLevel string) *Result {
	r.WindDayLevel = windDayLevel
	return r
}

func (r *Result) WithWindNight(windNight string) *Result {
	r.WindNight = windNight
	return r
}

func (r *Result) WithWindNightLevel(windNightLevel string) *Result {
	r.WindNightLevel = windNightLevel
	return r
}

func (r *Result) WithAlter(alter string) *Result {
	r.Alter = alter
	return r
}

func (r *Result) WithLed(led string) *Result {
	r.Led = led
	return r
}

func (r *Result) WithTypeTag(typeTag string) *Result {
	r.TypeTag = typeTag
	return r
}

func (r *Result) SetFocus(f string) {
	r.Focus = f
}

type Address struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Detail   string `json:"detail"`
}

type Location struct {
	Latitude  float64  `json:"latitude,omitempty"`
	Longitude float64  `json:"longitude,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
