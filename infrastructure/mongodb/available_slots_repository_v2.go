package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AvailableSlotsRepositoryV2 struct {
	readmodel.AvailableSlotsRepository

	db         *mongo.Database
	collection *mongo.Collection
}

func NewAvailableSlotsRepositoryV2(db *mongo.Database) *AvailableSlotsRepositoryV2 {
	return &AvailableSlotsRepositoryV2{
		db:         db,
		collection: db.Collection("available_slots"),
	}
}

func (m *AvailableSlotsRepositoryV2) AddSlot(s AvailableSlotV2) error {
	_, err := m.collection.InsertOne(context.TODO(), s)
	return err
}

func (m *AvailableSlotsRepositoryV2) HideSlot(slotId uuid.UUID) error {
	result, err := m.collection.UpdateOne(
		context.TODO(),
		bson.M{"id": slotId.String()},
		bson.D{{"$set", bson.D{{"isbooked", true}}}})

	if result.ModifiedCount == 0 {
		return fmt.Errorf("failed to hide slot")
	}

	return err
}

func (m *AvailableSlotsRepositoryV2) ShowSlot(slotId uuid.UUID) error {
	result, err := m.collection.UpdateOne(
		context.TODO(),
		bson.M{"id": slotId.String()},
		bson.D{{"$set", bson.D{{"isbooked", false}}}})

	if result.ModifiedCount == 0 {
		return fmt.Errorf("failed to show slot")
	}

	return err
}

func (m *AvailableSlotsRepositoryV2) DeleteSlot(slotId uuid.UUID) error {
	result, err := m.collection.DeleteOne(
		context.TODO(),
		bson.M{"id": slotId.String()})

	if result.DeletedCount == 0 {
		return fmt.Errorf("failed to delete slot")
	}

	return err
}

func (m *AvailableSlotsRepositoryV2) GetSlotsAvailableOn(date time.Time) ([]readmodel.AvailableSlot, error) {
	slots := make([]readmodel.AvailableSlot, 0)
	cur, err := m.collection.Find(context.TODO(), bson.D{{"date", date.Format("02-01-2006")}, {"isbooked", false}})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var result AvailableSlotV2
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}

		slots = append(slots, result.ToAvailableSlot())
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}
