package SensorReading

import (
	"time"
)

type AirQualityReading struct {
	SensorID  string    `json:"sensor_id"` // ID of the sensor that reported the reading
	Timestamp time.Time `json:"timestamp"` // Time at which the reading was recorded
	PM25      float64   `json:"pm25"`      // PM2.5 pollutant level
	CO2       float64   `json:"co2"`       // CO2 pollutant level
}
