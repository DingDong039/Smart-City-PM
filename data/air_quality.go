package data

import (
	"encoding/json"
	"fmt"
	"time"
)

// AirQualityReading represents a single sensor reading.
type AirQualityReading struct {
	Timestamp time.Time `json:"timestamp"`
	PM25      float64   `json:"pm25"`
	CO2       float64   `json:"co2"`
}

// ParseReadings takes a byte array containing JSON data and unmarshals it into a slice of AirQualityReading structs.
func ParseReadings(data []byte) ([]AirQualityReading, error) {
	var readings []AirQualityReading
	err := json.Unmarshal(data, &readings)
	if err != nil {
		return nil, fmt.Errorf("error parsing sensor data: %v", err)
	}
	return readings, nil
}

// CalculateAverage calculates the average value of each pollutant across all readings.
func CalculateAverage(readings []AirQualityReading) map[string]float64 {
	var totalPM25, totalCO2 float64
	count := len(readings)

	for _, reading := range readings {
		totalPM25 += reading.PM25
		totalCO2 += reading.CO2
	}

	averages := map[string]float64{
		"pm25": totalPM25 / float64(count),
		"co2":  totalCO2 / float64(count),
	}

	return averages
}

// FindHighestPollutantByHour groups readings by hour and finds the pollutant with the highest average for each hour.
func FindHighestPollutantByHour(readings []AirQualityReading) map[int]string {
	// Maps to store total pollutant values and count of readings for each hour
	hourlyTotals := make(map[int]map[string]float64) // e.g., hourlyTotals[hour]["pm25"]
	hourlyCounts := make(map[int]int)                // e.g., hourlyCounts[hour]

	// Initialize maps for each hour's data
	for _, reading := range readings {
		hour := reading.Timestamp.Hour()

		// Ensure the map for this hour is initialized
		if _, exists := hourlyTotals[hour]; !exists {
			hourlyTotals[hour] = map[string]float64{
				"pm25": 0,
				"co2":  0,
			}
		}

		// Accumulate total pollutant values
		hourlyTotals[hour]["pm25"] += reading.PM25
		hourlyTotals[hour]["co2"] += reading.CO2

		// Count the number of readings for each hour
		hourlyCounts[hour]++
	}

	// Map to store the highest pollutant by hour
	highestPollutantByHour := make(map[int]string)

	// Loop through each hour and calculate averages, then determine the highest pollutant
	for hour, totals := range hourlyTotals {
		// Calculate averages for the hour
		avgPM25 := totals["pm25"] / float64(hourlyCounts[hour])
		avgCO2 := totals["co2"] / float64(hourlyCounts[hour])

		// Determine the pollutant with the highest average for the hour
		if avgPM25 > avgCO2 {
			highestPollutantByHour[hour] = "pm25"
		} else {
			highestPollutantByHour[hour] = "co2"
		}
	}

	return highestPollutantByHour
}
