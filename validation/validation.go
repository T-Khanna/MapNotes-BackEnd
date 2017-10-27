package Validation

import (
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
)

func ValidateNoteRequest(note *models.Note) bool {

	return validateTitle(*note.Title) &&
		validateComments(*note.Comment) &&
		validateLongitude(*note.Longitude) &&
		validateLatitude(*note.Latitude) &&
		validateTimes(*note.StartTime, *note.EndTime)
}

func validateNotNil(object interface{}) bool {

	return object != nil

}

func validateLongitude(longitude float64) bool {

	return boundCheck(-180, 180, longitude)

}

func validateLatitude(latitude float64) bool {

	return boundCheck(-90, 90, latitude)

}

func validateTimes(startTime, endTime string) bool {

	return true

}

func validateTitle(title string) bool {

	return true

}

func validateComments(comments string) bool {

	return true

}

func boundCheck(low, high, value float64) bool {

	return value >= low && value <= high

}
