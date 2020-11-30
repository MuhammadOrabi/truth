package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Source ...
type Source struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" swaggerignore:"true"`
	Name   string             `json:"name" bson:"name" validate:"required"`
	URL    string             `json:"url" bson:"url" validate:"required"`
	Tags   []string           `json:"tags" bson:"tags" validate:"required"`
	Active bool               `json:"active" bson:"active" validate:"required"`
	Status string             `json:"status" bson:"status" validate:"required"`
}

// GetSources ...
func GetSources() ([]Source, error) {
	db := Load()

	cursor, err := db.Sources.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	sources := make([]Source, 0)
	if err = cursor.All(context.TODO(), &sources); err != nil {
		return nil, err
	}

	return sources, nil
}

// CreateSources ...
func CreateSources(source *Source) error {
	db := Load()

	_, err := db.Sources.InsertOne(context.TODO(), source)

	return err
}

// DeleteSource ...
func DeleteSource(SourceID string) error {
	db := Load()

	ID, err := primitive.ObjectIDFromHex(SourceID)
	if err != nil {
		return err
	}

	_, err = db.Sources.DeleteOne(context.TODO(), bson.M{"_id": ID})

	return err
}
