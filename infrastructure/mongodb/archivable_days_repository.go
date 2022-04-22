package mongodb

import (
	"context"
	"time"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArchivableDayRepository struct {
	readmodel.ArchivableDaysRepository

	collection *mongo.Collection
}

func NewArchivableDayRepository(db *mongo.Database) *AvailableSlotsRepository {
	return &AvailableSlotsRepository{
		collection: db.Collection("archivable_day"),
	}
}

func (m *AvailableSlotsRepository) Add(s readmodel.ArchivableDay) error {
	_, err := m.collection.InsertOne(context.TODO(), s)
	return err
}

func (m *AvailableSlotsRepository) FindAll(dateThreshold time.Time) ([]readmodel.ArchivableDay, error) {
	slots := make([]readmodel.ArchivableDay, 0)
	cur, err := m.collection.Find(context.TODO(), bson.D{{"date", bson.M{
		"$gte": primitive.NewDateTimeFromTime(dateThreshold),
	}}})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var result readmodel.ArchivableDay
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}

		slots = append(slots, result)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}
