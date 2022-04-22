package eventsourcing

type ColdStorage interface {
	SaveAll(events []interface{})
}