package readmodel

import "time"

type ArchivableDay struct {
	Id   string
	Date time.Time
}

func NewArchivableDay(id string, d time.Time) ArchivableDay {
	return ArchivableDay{
		Id: id,
		Date: d,
	}
}
