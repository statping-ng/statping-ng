package checkins

import (
	"time"

	"github.com/statping-ng/statping-ng/database"
)

type CheckinHitters struct {
	db database.Database
}

func (h CheckinHitters) Db() database.Database {
	return h.db
}

func (h CheckinHitters) First() *CheckinHit {
	var checkinHit CheckinHit
	h.db.Order("checkin_hits.id ASC").Limit(1).Find(&checkinHit)
	return &checkinHit
}

func (h CheckinHitters) Count() int {
	var count int
	h.db.Count(&count)
	return count
}

func (h CheckinHitters) Since(t time.Time) CheckinHitters {
	timestamp := db.FormatTime(t)
	return CheckinHitters{h.db.Where("checkin_hits.created_at > ?", timestamp)}
}

func (c *Checkin) LastHit() *CheckinHit {
	var hit CheckinHit
	dbHits.Where("checkin = ?", c.Id).Last(&hit)
	return &hit
}

func (c *Checkin) Hits() []*CheckinHit {
	var hits []*CheckinHit
	dbHits.Where("checkin = ?", c.Id).Order("id DESC").Limit(32).Find(&hits)
	c.AllHits = hits
	return hits
}

func (c *CheckinHit) Create() error {
	q := dbHits.Create(c)
	return q.Error()
}

func (c *CheckinHit) Update() error {
	q := dbHits.Update(c)
	return q.Error()
}

func (c *CheckinHit) Delete() error {
	q := dbHits.Delete(c)
	return q.Error()
}

func AllCheckinHits(serviceId int64) CheckinHitters {
	return CheckinHitters{dbHits.Joins("JOIN checkins ON checkins.id = checkin").Where("checkins.service = ?", serviceId).Order("checkin_hits.id DESC")}
}
