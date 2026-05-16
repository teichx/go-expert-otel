package models

import "math"

type Temperature struct {
	City string  `json:"city"`
	C    float64 `json:"temp_C"`
	F    float64 `json:"temp_F"`
	K    float64 `json:"temp_K"`
}

func NewTemperature(city string, celsius float64) Temperature {
	return Temperature{
		City: city,
		C:    math.Round(celsius*10) / 10,
		F:    math.Round((celsius*1.8+32)*10) / 10,
		K:    math.Round((celsius+273)*10) / 10,
	}
}
