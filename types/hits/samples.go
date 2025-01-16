package hits

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/mattn/go-sqlite3"
	"github.com/statping-ng/statping-ng/types"
	"github.com/statping-ng/statping-ng/utils"
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"time"
)

// max sample hits
var SampleHits = 200000.

func Samples() error {
	log.Infoln("Inserting Sample Service Hits...")
	records_1 := createHitsAt(1, -90 * types.Day, 1 * time.Minute)
	utils.Log.Infoln(fmt.Sprintf("Adding %v Hit records to service 1", len(records_1)))
	if err := gormbulk.BulkInsert(db.GormDB(), records_1, db.ChunkSize()); err != nil {
		log.Error(err)
		return err
	}
	// Statping - Statping Github
	records_2 := createHitsAt(2, -90 * types.Day, 1 * time.Minute)
	utils.Log.Infoln(fmt.Sprintf("Adding %v Hit records to service 2", len(records_2)))
	if err := gormbulk.BulkInsert(db.GormDB(), records_2, db.ChunkSize()); err != nil {
		log.Error(err)
		return err
	}
	records_3 := createHitsAt(3, -60 * types.Day, 1 * time.Minute)
	utils.Log.Infoln(fmt.Sprintf("Adding %v Hit records to service 3", len(records_3)))
	if err := gormbulk.BulkInsert(db.GormDB(), records_3, db.ChunkSize()); err != nil {
		log.Error(err)
		return err
	}
	records_4 := createHitsAt(4, -90 * types.Day, 1 * time.Minute)
	utils.Log.Infoln(fmt.Sprintf("Adding %v Hit records to service 4", len(records_4)))
	if err := gormbulk.BulkInsert(db.GormDB(), records_4, db.ChunkSize()); err != nil {
		log.Error(err)
		return err
	}
	records_5 := createHitsAt(5, -15 * types.Day, 1 * time.Minute)
	utils.Log.Infoln(fmt.Sprintf("Adding %v Hit records to service 5", len(records_5)))
	if err := gormbulk.BulkInsert(db.GormDB(), records_5, db.ChunkSize()); err != nil {
		log.Error(err)
		return err
	}
	records_6 := createHitsAt(6, -75 * types.Day, 1 * time.Minute)
	utils.Log.Infoln(fmt.Sprintf("Adding %v Hit records to service 6", len(records_6)))
	if err := gormbulk.BulkInsert(db.GormDB(), records_6, db.ChunkSize()); err != nil {
		log.Error(err)
		return err
	}
	// Service 7 is permanatly offline
	return nil
}

func createHitsAt(serviceID int64, daysToCreate time.Duration, createEvery time.Duration) []interface{} {
	log.Infoln(fmt.Sprintf("Adding Sample records to service #%d...", serviceID))

	createdAt := utils.Now().Add(daysToCreate)
	p := utils.NewPerlin(2, 2, 5, utils.Now().UnixNano())

	var records []interface{}
	for hi := 0.; hi <= SampleHits; hi++ {
		latency := p.Noise1D(hi / 500)

		hit := &Hit{
			Service:   serviceID,
			Latency:   int64(latency * 10000000),
			PingTime:  int64(latency * 5000000),
			CreatedAt: createdAt,
		}

		records = append(records, hit)

		if createdAt.After(utils.Now()) {
			break
		}
		createdAt = createdAt.Add(createEvery)
	}

	return records
}
