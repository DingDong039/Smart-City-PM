package persistence

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"Smart-City-PM/data"
)

// SaveToCSV saves the air quality readings to a CSV file in a new folder.
func SaveToCSV(readings []data.AirQualityReading, folderName string) error {
	// Create folder with the provided name if it doesn't exist
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create folder: %v", err)
	}

	// Generate the file path with a timestamp
	timestamp := time.Now().Format("20060102_150405")
	filePath := filepath.Join(folderName, fmt.Sprintf("air_quality_%s.csv", timestamp))

	// Create the CSV file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Timestamp", "PM2.5", "CO2"}); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	// Write data
	for _, reading := range readings {
		record := []string{
			reading.Timestamp.Format(time.RFC3339),
			fmt.Sprintf("%.2f", reading.PM25),
			fmt.Sprintf("%.2f", reading.CO2),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record to CSV: %v", err)
		}
	}

	fmt.Printf("CSV file saved successfully at: %s\n", filePath)
	return nil
}
