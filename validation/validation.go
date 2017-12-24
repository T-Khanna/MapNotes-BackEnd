package Validation

import (
	"time"

	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
)

const firstId int64 = 1

func ValidateNote(note *models.Note) bool {
	return validateTitle(note.Title) &&
		validateComment(note.Comment) &&
		validateLongitude(note.Longitude) &&
		validateLatitude(note.Latitude) &&
		validateTime(note.StartTime) &&
		validateTime(note.EndTime)
}

func ValidatePartialNote(note *models.Note) bool {
	return note.Id != nil && validateId(*note.Id) &&
		(note.Title == nil || validateTitle(note.Title)) &&
		(note.Comment == nil || validateComment(note.Comment)) &&
		(note.StartTime == nil || validateTime(note.StartTime)) &&
		(note.EndTime == nil || validateTime(note.EndTime)) &&
		(note.Longitude == nil || validateLongitude(note.Longitude)) &&
		(note.Latitude == nil || validateLatitude(note.Latitude)) &&
		(note.Tags == nil || validateTags(*note.Tags)) &&
		(note.Images == nil || validateImages(*note.Images))
}

func validateNotNil(object interface{}) bool {
	return object != nil
}

func validateId(id int64) bool {
	return id >= firstId
}

func validateLongitude(longitude *float64) bool {
	return longitude != nil && boundCheck(-180, 180, *longitude)
}

func validateLatitude(latitude *float64) bool {
	return latitude != nil && boundCheck(-90, 90, *latitude)
}

func validateTime(timestamp *string) bool {
	if timestamp == nil {
		return false
	}
	_, err := time.Parse(models.NoteTimeFormat, *timestamp)
	return err == nil
}

func validateTitle(title *string) bool {
	return title != nil
}

func validateComment(comments *string) bool {
	return comments != nil
}

func validateTags(tags []string) bool {
	return true
}

func validateImages(images []string) bool {
	return true
}

func boundCheck(low, high, value float64) bool {
	return value >= low && value <= high
}
