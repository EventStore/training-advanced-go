package readmodel

import "time"

type ArchivableDaysRepository interface {
	Add(day ArchivableDay) error
	FindAll(dateThreshold time.Time) ([]ArchivableDay, error)
}
