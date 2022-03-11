package mongodb

import (
	"context"
	"fmt"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookedSlotsRepository struct {
	readmodel.BookedSlotsRepository

	db         *mongo.Database
	collection *mongo.Collection
}

func NewBookedSlotsRepository(db *mongo.Database) *BookedSlotsRepository {
	return &BookedSlotsRepository{
		db:         db,
		collection: db.Collection("booked_slots"),
	}
}

func (m *BookedSlotsRepository) AddSlot(s readmodel.BookedSlot) error {
	_, err := m.collection.InsertOne(context.TODO(), s)
	return err
}

func (m *BookedSlotsRepository) MarkSlotAsBooked(slotId, patientId string) error {
	result, err := m.collection.UpdateOne(
		context.TODO(),
		bson.M{"slotid": slotId},
		bson.D{{"$set", bson.D{{"isbooked", true}, {"patientid", patientId}}}})

	if result.UpsertedCount == 0 {
		return fmt.Errorf("failed to mark slot as booked")
	}

	return err
}

func (m *BookedSlotsRepository) MarkSlotAsAvailable(slotId string) error {
	result, err := m.collection.UpdateOne(
		context.TODO(),
		bson.M{"slotid": slotId},
		bson.D{{"$set", bson.D{{"isbooked", true}, {"patientid", ""}}}})

	if result.UpsertedCount == 0 {
		return fmt.Errorf("failed to mark slot as available")
	}

	return err
}

func (m *BookedSlotsRepository) GetSlot(slotId string) (readmodel.BookedSlot, error) {
	slot := readmodel.BookedSlot{}
	result := m.collection.FindOne(context.TODO(), bson.M{"slotid": slotId})
	if result == nil {
		return slot, fmt.Errorf("failed to find slot")
	}

	err := result.Decode(&slot)
	if err != nil {
		return slot, err
	}

	return slot, nil
}

func (m *BookedSlotsRepository) CountByPatientAndMonth(patientId string, month int) (int, error) {
	result, err := m.collection.CountDocuments(context.TODO(), bson.D{{"patientid", patientId}, {"month", month}})
	if err != nil {
		return 0, err
	}

	return int(result), nil
}
