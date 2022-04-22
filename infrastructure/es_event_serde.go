package infrastructure

import (
	"encoding/json"
	"reflect"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/EventStore/training-introduction-go/eventsourcing"
	"github.com/gofrs/uuid"
)

type EsEventSerde struct {
	typeMapper *eventsourcing.TypeMapper
}

func NewEsEventSerde(tm *eventsourcing.TypeMapper) *EsEventSerde {
	return &EsEventSerde{
		typeMapper: tm,
	}
}

func (s *EsEventSerde) Serialize(event interface{}, m EventMetadata) (esdb.EventData, error) {
	typeToData, err := s.typeMapper.GetTypeToData(GetValueType(event))
	if err != nil {
		return esdb.EventData{}, err
	}

	id, err := uuid.NewV4()
	if err != nil {
		return esdb.EventData{}, err
	}

	name, jsonData := typeToData(event)
	dataBytes, err := json.Marshal(jsonData)
	if err != nil {
		return esdb.EventData{}, err
	}

	metadataBytes, err := json.Marshal(m)
	if err != nil {
		return esdb.EventData{}, err
	}

	return esdb.EventData{
		EventID:     id,
		ContentType: esdb.JsonContentType,
		EventType:   name,
		Data:        dataBytes,
		Metadata:    metadataBytes,
	}, nil
}

func (s *EsEventSerde) Deserialize(r *esdb.ResolvedEvent) (interface{}, *EventMetadata, error) {
	dataToType, err := s.typeMapper.GetDataToType(r.Event.EventType)
	if err != nil {
		return nil, nil, err
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(r.Event.Data, &m)
	if err != nil {
		return nil, nil, err
	}

	metadata := EventMetadata{}
	err = json.Unmarshal(r.Event.UserMetadata, &metadata)
	if err != nil {
		return nil, nil, err
	}

	return dataToType(m), &metadata, nil
}

func (s *EsEventSerde) SerializeSnapshot(snapshot interface{}, m eventsourcing.SnapshotMetadata) (esdb.EventData, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return esdb.EventData{}, err
	}

	name, err := s.typeMapper.GetTypeName(snapshot)
	if err != nil {
		return esdb.EventData{}, err
	}

	snapshotBytes, err := json.Marshal(snapshot)
	if err != nil {
		return esdb.EventData{}, err
	}

	metadataBytes, err := json.Marshal(m)
	if err != nil {
		return esdb.EventData{}, err
	}

	return esdb.EventData{
		EventID:     id,
		ContentType: esdb.JsonContentType,
		EventType:   name,
		Data:        snapshotBytes,
		Metadata:    metadataBytes,
	}, nil
}

func (s *EsEventSerde) DeserializeSnapshot(r *esdb.ResolvedEvent) (interface{}, *eventsourcing.SnapshotMetadata, error) {
	snapshot, err := s.deserializeToType(r.Event)
	if err != nil {
		return nil, nil, err
	}

	metadata := eventsourcing.SnapshotMetadata{}
	err = json.Unmarshal(r.Event.UserMetadata, &metadata)
	if err != nil {
		return nil, nil, err
	}

	return snapshot, &metadata, nil
}

func (s *EsEventSerde) SerializeCommand(command interface{}, m CommandMetadata) (esdb.EventData, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return esdb.EventData{}, err
	}

	name, err := s.typeMapper.GetTypeName(command)
	if err != nil {
		return esdb.EventData{}, err
	}

	commandBytes, err := json.Marshal(command)
	if err != nil {
		return esdb.EventData{}, err
	}

	metadataBytes, err := json.Marshal(m)
	if err != nil {
		return esdb.EventData{}, err
	}

	return esdb.EventData{
		EventID:     id,
		ContentType: esdb.JsonContentType,
		EventType:   name,
		Data:        commandBytes,
		Metadata:    metadataBytes,
	}, nil
}

func (s *EsEventSerde) DeserializeCommand(r *esdb.ResolvedEvent) (interface{}, CommandMetadata, error) {
	metadata := CommandMetadata{}
	cmd, err := s.deserializeToType(r.Event)
	if err != nil {
		return nil, metadata, err
	}

	err = json.Unmarshal(r.Event.UserMetadata, &metadata)
	if err != nil {
		return nil, metadata, err
	}

	return cmd, metadata, nil
}

func (s *EsEventSerde) deserializeToType(e *esdb.RecordedEvent) (interface{}, error) {
	t, err := s.typeMapper.GetType(e.EventType)
	if err != nil {
		return nil, err
	}

	cmd := reflect.New(t).Interface()
	err = json.Unmarshal(e.Data, &cmd)
	if err != nil {
		return nil, err
	}

	return reflect.ValueOf(cmd).Elem().Interface(), nil
}
