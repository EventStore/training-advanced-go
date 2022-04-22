package eventsourcing

import (
	"fmt"
	"reflect"
)

type DataToType func(map[string]interface{}) interface{}
type TypeToData func(interface{}) map[string]interface{}
type TypeToDataWithName func(interface{}) (string, map[string]interface{})
type CreateSnapshot func() interface{}

type TypeMapper struct {
	dataToType map[string]DataToType
	typeToData map[reflect.Type]TypeToDataWithName

	snapshotCreation   map[string]reflect.Type
	snapshotTypeToName map[reflect.Type]string
}

func NewTypeMapper() *TypeMapper {
	return &TypeMapper{
		dataToType:         make(map[string]DataToType),
		typeToData:         make(map[reflect.Type]TypeToDataWithName),
		snapshotCreation:   make(map[string]reflect.Type),
		snapshotTypeToName: make(map[reflect.Type]string),
	}
}

func (tm *TypeMapper) MapEvent(eventType reflect.Type, name string, dt DataToType, td TypeToData) error {

	if name == "" {
		return fmt.Errorf("need name for type mapping")
	}

	if _, exists := tm.typeToData[eventType]; exists {
		return nil
	}

	tm.dataToType[name] = dt
	tm.typeToData[eventType] = func(t interface{}) (string, map[string]interface{}) {
		data := td(t)
		return name, data
	}
	return nil
}

func (tm *TypeMapper) GetDataToType(typeName string) (DataToType, error) {
	if dt, exists := tm.dataToType[typeName]; exists {
		return dt, nil
	}
	return nil, fmt.Errorf("failed to find type mapped with '%s'", typeName)
}

func (tm *TypeMapper) GetTypeToData(t reflect.Type) (TypeToDataWithName, error) {
	if td, exists := tm.typeToData[t]; exists {
		return td, nil
	}
	return nil, fmt.Errorf("failed to find name mapped with '%s'", t)
}

func (tm *TypeMapper) RegisterType(t reflect.Type, typeName string, c CreateSnapshot) error {
	if typeName == "" {
		return fmt.Errorf("need type name for registration")
	}

	if _, exists := tm.snapshotTypeToName[t]; exists {
		return nil
	}

	tm.snapshotTypeToName[t] = typeName
	tm.snapshotCreation[typeName] = t
	return nil
}

func (tm *TypeMapper) GetTypeName(v interface{}) (string, error) {
	t := getValueType(v)
	if name, exists := tm.snapshotTypeToName[t]; exists {
		return name, nil
	}

	return "", fmt.Errorf("type '%v' not registered", t)
}

func (tm *TypeMapper) GetType(typeName string) (reflect.Type, error) {
	if typeName == "" {
		return nil, fmt.Errorf("need type name for type creation")
	}

	if createSnapshot, exists := tm.snapshotCreation[typeName]; exists {
		return createSnapshot, nil
	}

	return nil, fmt.Errorf("type '%s' not registered", typeName)
}
