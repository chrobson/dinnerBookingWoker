package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	url              = "https://book.dinnerbooking.com/pl/pl-PL/time_days/calendar_days/3203/1/2/2023-10-01/2023-12-31.json"
	retryInterval    = 60 * time.Second
	timeLayout       = "2006-01-02"
	closedRestaurant = 1
)

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
			time.Sleep(retryInterval)
			continue
		}
		defer response.Body.Close()

		// Parse the JSON response into a Response struct
		var jsonResponse Response
		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&jsonResponse); err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
			time.Sleep(retryInterval)
			continue
		}

		// Iterate through the availability data and display information
		for _, availability := range jsonResponse.Days {
			if availability.TimeDay.Closed == closedRestaurant {
				continue // Skip closed restaurant days
			}

			dateStr := availability.TimeDay.Day
			availabilityStatus := availability.Booking.AvailabilityStatus

			switch availabilityStatus {
			case 0, 2:
				continue // Skip days with no available tables or unknown status
			case 1:
				fmt.Printf("On %s: Tables are available.\n", dateStr)
			default:
				fmt.Printf("On %s: Availability status unknown.\n", dateStr)
			}
		}

		// Sleep before checking again
		time.Sleep(retryInterval)
	}
}

func main() {
	checkAvailability()
}
