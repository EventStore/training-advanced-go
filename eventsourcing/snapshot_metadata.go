package eventsourcing

type SnapshotMetadata struct {
	Version int `json:"version"`
}

func NewSnapshotMetadata(version int) SnapshotMetadata {
	return SnapshotMetadata{
		Version: version,
	}
}
