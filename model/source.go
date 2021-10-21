package model

import (
	"log"
	"truth/storage"

	"github.com/lib/pq"
)

// Source ...
type Source struct {
	ID     uint           `json:"id,omitempty" gorm:"primaryKey,autoIncrement" swaggerignore:"true"`
	Name   string         `json:"name" bson:"name" validate:"required"`
	URL    string         `json:"url" bson:"url" validate:"required"`
	Tags   pq.StringArray `json:"tags" gorm:"type:text[]" validate:"required" swaggertype:"array,string"`
	Active bool           `json:"active" bson:"active" validate:"required"`
	Status string         `json:"status" bson:"status" validate:"required"`
}

// GetSources ...
func GetSources() []Source {
	db := storage.GetDBInstance()

	var sources []Source
	db.Find(&sources)

	return sources
}

// CreateSources ...
func CreateSources(source *Source) {
	db := storage.GetDBInstance()

	result := db.Create(&source)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

// DeleteSource ...
func DeleteSource(ID uint) {
	db := storage.GetDBInstance()

	db.Delete(&Source{}, ID)
}
