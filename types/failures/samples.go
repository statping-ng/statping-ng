package failures

import (
	"fmt"
	"github.com/statping-ng/statping-ng/types"
	"github.com/statping-ng/statping-ng/utils"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"time"
	"math/rand"
)

var (
	log = utils.Log.WithField("type", "failure")
)

func Example() Failure {
	return Failure{
		Id:        48533,
		Issue:     "Response did not response a 200 status code",
		Method:    "",
		MethodId:  0,
		ErrorCode: 404,
		Service:   1,
		Checkin:   0,
		PingTime:  48309,
		Reason:    "status_code",
		CreatedAt: utils.Now(),
	}
}

func createFailuresForService(serviceID int64, start time.Time, end time.Time, chanceOfFailure float64) []interface{} {
	var records []interface{}
	currentTime := start
	for currentTime.Before(end) {
		// Randomly decide if an outage should occur
		if rand.Float64() < chanceOfFailure {
			// Determine random outage length
			// This is so we can display all the different
			// severities of outages
			outageLengths := []time.Duration{
				4 * types.Day,
				24 * types.Hour,
				3 * types.Hour,
				60 * types.Minute,
				5 * types.Minute,
			}
			outageLength := outageLengths[rand.Intn(len(outageLengths))]

			// Simulate failures for the duration of the outage
			outageEnd := currentTime.Add(outageLength)
			// Initialize the current day and count of records.
			currentDay := currentTime.Day()
			dayRecordCount := 0
			
			for currentTime.Before(outageEnd) {
				f1 := &Failure{
					Service:   serviceID,
					Issue:     "Server failure",
					Reason:    "lookup",
					CreatedAt: currentTime.UTC(),
				}
				records = append(records, f1)
				dayRecordCount++
				currentTime = currentTime.Add(1 * time.Minute) // simulate the next ping after 1 minute

				// Check if we have moved to the next day
				if currentDay != currentTime.Day() {
					// Log the number of failures created for the previous day
					utils.Log.Infoln("Created", dayRecordCount, "Failures for Service", serviceID, "for day", currentTime.Add(-1 * time.Minute).UTC())

					// Reset for the new day
					currentDay = currentTime.Day()
					dayRecordCount = 0
				}
			}

			// Don't forget to log the count for the last day of the outage period as well
			if dayRecordCount > 0 {
				utils.Log.Infoln("Created", dayRecordCount, "Failures for Service", serviceID, "for day", currentTime.UTC())
			}
		} else {
			currentTime = currentTime.Add(5 * time.Minute) // No outage, move to the next random check period
		}
	}
	return records
}

func Samples() error {
	utils.Log.Infoln("Inserting Sample Service Failures...")
	endDate := utils.Now()                        // Up to current time

	chanceOfFailure := 0.0003 // very small chance of starting an outage at any given minute

	// Only add failures to services 3 and 4
	records_3 := createFailuresForService(3, utils.Now().Add(-60 * types.Day), endDate, chanceOfFailure)
	utils.Log.Infoln(fmt.Sprintf("Adding %v Failure records to service 3", len(records_3)))
	if err := gormbulk.BulkInsert(db.GormDB(), records_3, db.ChunkSize()); err != nil {
		log.Error(err)
		return err
	}
	records_4 := createFailuresForService(4, utils.Now().Add(-90 * types.Day), endDate, chanceOfFailure)
	utils.Log.Infoln(fmt.Sprintf("Adding %v Failure records to service 4", len(records_4)))
	if err := gormbulk.BulkInsert(db.GormDB(), records_4, db.ChunkSize()); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
