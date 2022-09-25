package services

import (
	"time"

	"github.com/statping-ng/statping-ng/types/checkins"
	"github.com/statping-ng/statping-ng/types/hits"
)

func (s *Service) HitsColumnID() (string, int64) {
	return "service", s.Id
}

func (s *Service) FirstHit() *hits.Hit {
	return hits.AllHits(s).First()
}

func (s *Service) LastHit() *hits.Hit {
	return hits.AllHits(s).Last()
}

func (s *Service) AllHits() hits.Hitters {
	return hits.AllHits(s)
}

func (s *Service) AllCheckinHits() checkins.CheckinHitters {
	return checkins.AllCheckinHits(s.Id)
}

func (s *Service) FirstCheckinHit() *checkins.CheckinHit {
	return checkins.AllCheckinHits(s.Id).First()
}

func (s *Service) HitsSince(t time.Time) hits.Hitters {
	return hits.Since(t, s)
}
