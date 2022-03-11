package doctorday

type SlotStatus int

const (
	SlotNotScheduled SlotStatus = iota
	SlotAvailable
	SlotBooked
)
