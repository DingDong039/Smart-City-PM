package main

import (
	"Smart-City-PM/data"
	"Smart-City-PM/persistence"
	"fmt"
	"log"
	"os"
)

func main() {
	// Example: Read data from JSON file
	jsonData, err := os.ReadFile("./Mock/caution.json")
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// Parse the sensor data
	readings, err := data.ParseReadings(jsonData)
	if err != nil {
		log.Fatalf("Failed to parse readings: %v", err)
	}

	// Calculate averages
	averages := data.CalculateAverage(readings)
	fmt.Printf("Averages: %v\n", averages)

	// Find highest pollutant by hour
	highestByHour := data.FindHighestPollutantByHour(readings)
	fmt.Printf("Highest pollutant by hour: %v\n", highestByHour)

	// Save the data to CSV
	err = persistence.SaveToCSV(readings, "LogFile")
	if err != nil {
		log.Fatalf("Failed to save CSV: %v", err)
	}
}
