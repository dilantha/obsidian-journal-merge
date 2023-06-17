package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Get the 'Journal' folder path from command line argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide the 'Journal' folder path as a command line argument.")
		return
	}
	journalFolder := os.Args[1]
	dailyFolder := filepath.Join(journalFolder, "Daily")
	monthlyFolder := filepath.Join(journalFolder, "Monthly")

	// Create the 'Monthly' folder if it doesn't exist
	if _, err := os.Stat(monthlyFolder); os.IsNotExist(err) {
		os.Mkdir(monthlyFolder, 0755)
	}

	// Read the list of daily journal files
	files, err := ioutil.ReadDir(dailyFolder)
	if err != nil {
		fmt.Println("Error reading daily folder:", err)
		return
	}

	// Map to store the content for each month
	monthlyNotes := make(map[string]string)

	// Iterate through the daily files
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		// Extract the year and month from the file name
		yearMonth := strings.Split(file.Name(), "-")[0:2]
		year, month := yearMonth[0], yearMonth[1]

		// Create the monthly folder if it doesn't exist
		if _, err := os.Stat(monthlyFolder); os.IsNotExist(err) {
			os.Mkdir(monthlyFolder, 0755)
		}

		// Read the content of the daily file
		filePath := filepath.Join(dailyFolder, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
			continue
		}

		// Append the content to the monthly note
		monthlyNote := monthlyNotes[fmt.Sprintf("%s-%s", year, month)]
		monthlyNote += fmt.Sprintf("## %s\n\n%s\n\n", file.Name(), string(content))
		monthlyNotes[fmt.Sprintf("%s-%s", year, month)] = monthlyNote
	}

	// Create the monthly note files
	for key, value := range monthlyNotes {
		monthNotePath := filepath.Join(monthlyFolder, key+".md")
		if err := ioutil.WriteFile(monthNotePath, []byte(value), 0644); err != nil {
			fmt.Printf("Error creating monthly note %s: %v\n", key, err)
		}
	}

	fmt.Println("Journal notes merged into monthly folders successfully.")
}
