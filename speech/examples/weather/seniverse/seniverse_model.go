package seniverse

type LocResp struct {
	Results []*Location `json:"results"`
}

type Location struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CountryCode    string `json:"country"`
	Path           string `json:"path"`
	TimeZone       string `json:"timezone"`
	TimeZoneOffset string `json:"timezone_offset"`
}

type LogLat struct {
	log float64 `json:"longitude"`
	lat float64 `json:"latitude"`
}

type ApiStatusResp struct {
	Status string `json:"status"`
	Code   string `json:"status_code"`
}

type NowCondResp struct {
	Results []*NowCondResult `json:"results"`
}

type NowCondResult struct {
	Location   *Location `json:"location"`
	Now        *NowCond  `json:"now"`
	LastUpdate string    `json:"last_update"`
}

type NowCond struct {
	Text          string `json:"text"`
	Code          string `json:"code"`
	Temperature   string `json:"temperature"`
	FeelsLike     string `json:"feels_like"`
	Pressure      string `json:"pressure"`
	Humidity      string `json:"humidity"`
	Visibility    string `json:"visibility"`
	WindDirection string `json:"wind_direction"`
	WindDegree    string `json:"wind_direction_degree"`
	WindSpeed     string `json:"wind_speed"`
	WindScale     string `json:"wind_scale"`
	Clouds        string `json:"clouds"`
	DewPoint      string `json:"dew_point"`
}

type DailyCondResp struct {
	Results []*DailyCondResults `json:"results"`
}

type DailyCondResults struct {
	Location   *Location    `json:"location"`
	Daily      []*DailyCond `json:"daily"`
	LastUpdate string       `json:"last_update"`
}

type DailyCond struct {
	Date          string `json:"date"`
	TextDay       string `json:"text_day"`
	CodeDay       string `json:"code_day"`
	TextNight     string `json:"text_night"`
	CodeNight     string `json:"code_night"`
	High          string `json:"high"`
	Low           string `json:"low"`
	Precip        string `json:"precip"`
	WindDirection string `json:"wind_direction"`
	WindDegree    string `json:"wind_direction_degree"`
	WindSpeed     string `json:"wind_speed"`
	WindScale     string `json:"wind_scale"`
}

type NowAirResp struct {
	Results []*DailyCondResults `json:"results"`
}

type NowAirResults struct {
	Location   *Location `json:"location"`
	Air        CityAir   `json:"air"`
	LastUpdate string    `json:"last_update"`
}

type CityAir struct {
	CityCond     NowAirCond       `json:"city"`
	StationsCond []StationAirCond `json:"stations"`
}

type NowAirCond struct {
	AirAqi
	PrimaryPollutant string `json:"primary_pollutant"`
	LastUpdate       string `json:"last_update"`
}

type DailyAirResp struct {
	Results []*DailyAirResults `json:"results"`
}

type DailyAirResults struct {
	Location   *Location       `json:"location"`
	Daily      []*DailyAirCond `json:"daily"`
	LastUpdate string          `json:"last_update"`
}

type DailyAirCond struct {
	AirAqi
	Date string `json:"date"`
}

type AirAqi struct {
	Aqi     string `json:"aqi"`
	Pm25    string `json:"pm25"`
	Pm10    string `json:"pm10"`
	So2     string `json:"so2"`
	No2     string `json:"no2"`
	Co      string `json:"co"`
	O3      string `json:"o3"`
	Quality string `json:"quality"`
}

type StationAirCond struct {
	AirAqi
	LogLat
	PrimaryPollutant string `json:"primary_pollutant"`
	LastUpdate       string `json:"last_update"`
	StationName      string `json:"station"`
}
