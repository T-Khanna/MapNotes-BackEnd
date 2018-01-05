package Validation

import (
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"log"
	"time"
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
	if longitude != nil && boundCheck(-180, 180, *longitude) {
		return true
	} else {
		log.Println("Failed longitude check")
		return false
	}
}

func validateLatitude(latitude *float64) bool {
	if latitude != nil && boundCheck(-90, 90, *latitude) {
		return true
	} else {
		log.Println("Failed latitude check")
		return false
	}
}

func validateTime(timestamp *string) bool {
	if timestamp == nil {
		log.Println("Error validating timestamp: empty timestamp")
		return false
	}
	_, err := time.Parse(models.NoteTimeFormat, *timestamp)
	if err == nil {
		return true
	} else {
		log.Println("Error validating timestamp: parse error")
		return false
	}
}

func validateTitle(title *string) bool {
	if title != nil {
		return true
	} else {
		log.Println("Failed title validation")
		return false
	}
}

func validateComment(comments *string) bool {
	if comments != nil {
		return true
	} else {
		log.Println("Failed title validation")
		return false
	}
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
