package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Define the URL for fetching availability data
var url = "https://book.dinnerbooking.com/pl/pl-PL/time_days/calendar_days/3203/1/2/2023-10-01/2023-12-31.json"

// Availability represents the availability data structure
type Availability struct {
	TimeDay struct {
		AreaID          int    `json:"area_id"`
		Day             string `json:"day"`
		Closed          int    `json:"closed"`
		HasBookingMenus bool   `json:"has_booking_menus"`
	} `json:"TimeDay"`
	TimeNote []interface{} `json:"TimeNote"`
	Booking  struct {
		Pax                int `json:"pax"`
		AvailabilityStatus int `json:"availability_status"`
	} `json:"Booking"`
}

// Response represents the JSON response structure
type Response struct {
	Days []Availability `json:"days"`
}

// Function to check availability for a specific date
func checkAvailability() {
	for {
		// Make the HTTP GET request
		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("An error occurred: %v\n", err)
			time.Sleep(60 * time.Second) // Sleep for 60 seconds before trying again
			continue
		}
		defer response.Body.Close()

		// Parse the JSON response into a Response struct
		var jsonResponse Response
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&jsonResponse)
		if err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
			time.Sleep(60 * time.Second) // Sleep for 60 seconds before trying again
			continue
		}

		// Iterate through the availability data and display information
		for _, availability := range jsonResponse.Days {
			dateStr := availability.TimeDay.Day
			closed := availability.TimeDay.Closed
			availabilityStatus := availability.Booking.AvailabilityStatus

			// Skip days where the restaurant is closed
			if closed == 1 {
				continue // Skip to the next day
			}

			if availabilityStatus == 0 || availabilityStatus == 2 {
				continue // Skip to the next day
			} else if availabilityStatus == 1 {
				fmt.Printf("On %s: Tables are available.\n", dateStr)
			} else {
				fmt.Printf("On %s: Availability status unknown.\n", dateStr)
			}
		}

		// Sleep for 60 seconds (1 minute) before checking again
		time.Sleep(60 * time.Second)
	}
}

func main() {
	checkAvailability()
}
