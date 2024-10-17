package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type SensorReading struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	PM25      float64   `json:"pm25"`
	CO2       float64   `json:"co2"`
}

var readings []SensorReading

// calculateAverage calculates the average value of each pollutant across all provided readings.
func calculateAverage(readings []SensorReading) map[string]float64 {
	pollutants := map[string]float64{
		"pm25": 0,
		"co2":  0,
	}
	count := len(readings)

	if count == 0 {
		return pollutants
	}

	for _, reading := range readings {
		pollutants["pm25"] += reading.PM25
		pollutants["co2"] += reading.CO2
	}

	// Calculate average values
	pollutants["pm25"] /= float64(count)
	pollutants["co2"] /= float64(count)

	return pollutants
}

func findHighestPollutantByHour(readings []SensorReading) map[int]string {
	// Maps to store pollutant values and the corresponding timestamp for each hour
	hourlyReadings := make(map[int]map[string][]float64)
	hourlyTimestamps := make(map[int]time.Time)

	// Group readings by hour
	for _, reading := range readings {
		hour := reading.Timestamp.Hour()

		// Initialize the map for each hour if not already initialized
		if _, exists := hourlyReadings[hour]; !exists {
			hourlyReadings[hour] = map[string][]float64{
				"pm25": {},
				"co2":  {},
			}
		}

		// Append the values to the corresponding hour
		hourlyReadings[hour]["pm25"] = append(hourlyReadings[hour]["pm25"], reading.PM25)
		hourlyReadings[hour]["co2"] = append(hourlyReadings[hour]["co2"], reading.CO2)

		// Store the first timestamp in the hour for later formatting (dd/mm/yyyy and HH:MM)
		// Assuming the timestamp for the first reading of the hour is representative
		if _, exists := hourlyTimestamps[hour]; !exists {
			hourlyTimestamps[hour] = reading.Timestamp
		}
	}

	// Find the highest average pollutant for each hour
	highestPollutants := make(map[int]string)

	for hour, pollutants := range hourlyReadings {
		// Handle cases where there might be no readings for pm25 or co2
		pm25Values := pollutants["pm25"]
		co2Values := pollutants["co2"]

		if len(pm25Values) == 0 || len(co2Values) == 0 {
			// If no readings are available for pm25 or co2, skip the hour
			fmt.Printf("No data for hour %d\n", hour)
			continue
		}

		// Calculate averages
		pm25Sum := sum(pm25Values)
		co2Sum := sum(co2Values)
		pm25Avg := pm25Sum / float64(len(pm25Values))
		co2Avg := co2Sum / float64(len(co2Values))

		// Format the averages to 2 decimal places
		pm25AvgFormatted := fmt.Sprintf("%.2f", pm25Avg)
		co2AvgFormatted := fmt.Sprintf("%.2f", co2Avg)

		// Get the formatted date (dd/mm/yyyy) for the first reading in the hour
		dateFormatted := hourlyTimestamps[hour].Format("02/01/2006")

		// Get the formatted time (HH:MM) for the first reading in the hour
		timeFormatted := hourlyTimestamps[hour].Format("15:04")

		// Print the formatted averages (optional)
		fmt.Printf("Date: %s | Time: %s | PM2.5: %s | CO2: %s\n", dateFormatted, timeFormatted, pm25AvgFormatted, co2AvgFormatted)

		// Determine the pollutant with the highest average
		if pm25Avg > co2Avg {
			highestPollutants[hour] = "pm25"
		} else {
			highestPollutants[hour] = "co2"
		}
	}

	return highestPollutants
}

// Helper function to sum values in a slice.
func sum(values []float64) float64 {
	var total float64
	for _, value := range values {
		total += value
	}
	return total
}

func main() {
	file, err := os.Open("Mock/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the records and print them
	for _, record := range records[1:] { // Skip the header row
		pm25, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		co2, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		timestamp, err := time.Parse(time.RFC3339, record[1])
		if err != nil {
			log.Fatal(err)
		}

		// Append the reading to the slice
		reading := SensorReading{
			SensorID:  record[0],
			Timestamp: timestamp,
			PM25:      pm25,
			CO2:       co2,
		}
		readings = append(readings, reading)
	}

	// Calculate the average pollutant values
	avgPollutants := calculateAverage(readings)
	fmt.Println("Average Pollutants:", avgPollutants)

	// Find the highest pollutant by hour
	highestPollutants := findHighestPollutantByHour(readings)
	fmt.Println("Highest Pollutant by Hour:", highestPollutants)

}
