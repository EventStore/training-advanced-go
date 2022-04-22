package eventsourcing

type AggregateRootSnapshot interface {
	GetSnapshot() interface{}
	GetSnapshotVersion() int
	LoadSnapshot(snapshot interface{}, version int)
}

type AggregateRootSnapshotBase struct {
	AggregateRootBase
	AggregateRootSnapshot

	snapshotVersion int
	loadSnapshot    func(snapshot interface{})
	getSnapshot     func() interface{}
}

func NewAggregateRootSnapshot() AggregateRootSnapshotBase {
	return AggregateRootSnapshotBase{
		AggregateRootBase: NewAggregateRoot(),
	}
}

func (a *AggregateRootSnapshotBase) RegisterSnapshot(load func(snapshot interface{}), get func() interface{}) {
	a.loadSnapshot = load
	a.getSnapshot = get
}

func (a *AggregateRootSnapshotBase) LoadSnapshot(snapshot interface{}, version int) {
	a.loadSnapshot(snapshot)
	a.version = version
	a.snapshotVersion = version
}

func (a *AggregateRootSnapshotBase) GetSnapshot() interface{} {
	return a.getSnapshot()
}

func (a *AggregateRootSnapshotBase) GetSnapshotVersion() int {
	return a.snapshotVersion
}
