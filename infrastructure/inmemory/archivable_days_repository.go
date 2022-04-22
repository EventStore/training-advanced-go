package inmemory

import (
	"time"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
)

type ArchivableDaysRepository struct {
	readmodel.ArchivableDaysRepository

	archivableDays []readmodel.ArchivableDay
}

func NewArchivableDaysRepository() *ArchivableDaysRepository {
	return &ArchivableDaysRepository{}
}

func (r *ArchivableDaysRepository) Add(d readmodel.ArchivableDay) error {
	r.archivableDays = append(r.archivableDays, d)
	return nil
}

func (r *ArchivableDaysRepository) FindAll(dateThreshold time.Time) ([]readmodel.ArchivableDay, error) {
	days := make([]readmodel.ArchivableDay, 0)
	for _, day := range r.archivableDays {
		if day.Date.Before(dateThreshold) || day.Date.Equal(dateThreshold) {
			days = append(days, day)
		}
	}
	return days, nil
}
