package models

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

//Struct to hold the insertion count
type SynchronisedNoteCounter struct {
	sync.RWMutex
	counter int
	target  int
}

/*
  A counter that is used to keep track of how many notes we have inserted
  during the run time of the server.
*/
var insertionNoteCounter = SynchronisedNoteCounter{counter: 0, target: 2}

// TODO: Change to StartTime and EndTime, and add json tags in camel case.
type Note struct {
	Title     *string   `json:"title,omitempty"`
	Comment   *string   `json:"comment,omitempty"`
	StartTime *string   `json:"start_time,omitempty"`
	EndTime   *string   `json:"end_time,omitempty"`
	Longitude *float64  `json:"longitude,omitempty"`
	Latitude  *float64  `json:"latitude,omitempty"`
	Id        *int64    `json:"id,omitempty"`
	Users     *[]User   `json:"users,omitempty"`
	Tags      *[]string `json:"tags,omitempty"`
}

// Possibly will be a similar struct for any future structs we perform CRUD on.
// Refactor when the time comes.
type NoteOperations struct {
	GetAll          func() ([]Note, error)
	GetActiveAtTime func(string) ([]Note, error)
	GetByUser       func(string) ([]Note, error)
	Create          func(*Note) (int64, error)
	Update          func(*Note) error
	Delete          func(int64) error
	Merge           func([]int64, Note) (int64, error)
}

// Exported API. Use as models.Notes.Create(..)
// FIXME: To compile, this is named plurarly whilst the actual Note struct must
//        be named singularly. This is inconsistent and an easy to get wrong
//        oddity of the code.
var Notes = NoteOperations{
	GetAll:          getAllNotes,
	GetActiveAtTime: getNotesActiveAtTime,
	GetByUser:       getNotesActiveByUser,
	Create:          createNote,
	Update:          updateNote,
	Delete:          deleteNote,
	Merge:           mergeNotes,
}

func mergeNotes(oldIds []int64, newNote Note) (id int64, err error) {

	id, err = createNote(&newNote)

	Comments.Merge(oldIds, id)
	deleteNotes(oldIds)

	return


}

/*
  Deletes notes by note_ids
*/
func deleteNotes(deleteids []int64) (err error) {

	for i := 0; i < len(deleteids); i++ {

		err = deleteNote(deleteids[i])

		if err != nil {

			return

		}

	}

	return

}

func createNote(note *Note) (int64, error) {
	// Prepare sql that inserts the note and returns the new id.
	stmt, err := db.Prepare("INSERT INTO notes(title, comments, startTime, endTime, longitude, latitude) VALUES($1, $2, $3, $4, $5, $6) RETURNING id")

	if err != nil {
		return -1, err
	}

	// Execute the INSERT statement, marshalling the returned id into an int64.
	var id int64
	err = stmt.QueryRow(note.Title, note.Comment, note.StartTime, note.EndTime,
		note.Longitude, note.Latitude).Scan(&id)

	if err != nil {
		return -1, err
	}

	// Get tags from note and insert each tag in database
	tags := note.Tags

	for _, t := range *tags {

		err := linkTag(t, id)

		if err != nil {
			return id, err
		}
	}

	users := note.Users

	for _, u := range *users {

		_, uid := GetUserId(u)

		err := NotesUsers.Insert(id, uid)

		if err != nil {
			return -1, err
		}
	}

	//Increment counter
	insertionNoteCounter.Lock()
	insertionNoteCounter.counter += 1
	insertionNoteCounter.Unlock()

	return id, nil
}

func updateNote(note *Note) error {
	/*
	   To implement partial updates:

	   The fields of the Note struct must be pointers, so that they we can
	   distinguish when they've been ommitted from the JSON by checking if the
	   pointer is nil.

	   To dynamically construct the query based on what columns are included, uses a
	   bunch of if statements that check if a column is present and, if so, appends
	   "column name = $n" to the byte buffer.

	   Uses a byte buffer to avoid re-concatenating strings over and over.
	*/

	if note.Id == nil {
		return errors.New("Error: Attempting to update Note but ID not provided")
	}

	// This will be the parameter number of the column-to-update's value in the
	// query that is constructed.. If a column needs to be updated and it's the
	// 'numCols'th column to be added to the query, then it will become parameter
	// '$numCols' in the query.
	numCols := 1

	// Contains the values of the columns to be added. Each time a non-nil field
	// is found in note, that field will be appended to values and numCols
	// incremented. Thus, values[i] will match $i in the query.
	values := []interface{}{}

	// Initialise buffer in which to build the query string.
	var buffer bytes.Buffer
	buffer.WriteString("UPDATE notes SET ")

	if note.Title != nil {
		buffer.WriteString(fmt.Sprintf("title = $%d, ", numCols))
		numCols++
		values = append(values, *note.Title)
	}

	if note.Comment != nil {
		buffer.WriteString(fmt.Sprintf("comments = $%d, ", numCols))
		numCols++
		values = append(values, *note.Comment)
	}

	if note.StartTime != nil {
		buffer.WriteString(fmt.Sprintf("startTime = $%d, ", numCols))
		numCols++
		values = append(values, *note.StartTime)
	}

	if note.EndTime != nil {
		buffer.WriteString(fmt.Sprintf("endTime = $%d, ", numCols))
		numCols++
		values = append(values, *note.EndTime)
	}

	if note.Longitude != nil {
		buffer.WriteString(fmt.Sprintf("Longitude = $%d, ", numCols))
		numCols++
		values = append(values, *note.Longitude)
	}

	if note.Latitude != nil {
		buffer.WriteString(fmt.Sprintf("Latitude = $%d, ", numCols))
		numCols++
		values = append(values, *note.Latitude)
	}

	if note.Tags != nil {
		buffer.WriteString(fmt.Sprintf("Tags = $%d, ", numCols))
		numCols++
		values = append(values, *note.Tags)
	}

	// FIXME: For some reason, bytes.TrimSuffix does not exist, so the trailing
	// comma cannot be removed. Instead, add a superflous 'id = id'.
	buffer.WriteString(fmt.Sprintf("id = %d", *note.Id))

	buffer.WriteString(fmt.Sprintf(" WHERE id = %d;", *note.Id))

	query := buffer.String()

	_, err := db.Exec(query, values...)

	return err
}

//Duplication here with DeleteUser
func deleteNote(id int64) error {
	stmt, prepErr := db.Prepare("DELETE FROM Notes WHERE id = $1")

	if prepErr != nil {
		return prepErr
	}

	_, execErr := stmt.Exec(id)

	if execErr != nil {
		return execErr
	}

	return nil
}

func filterNotes(whereClause string) ([]Note, error) {
	notesWithUsersQuery := fmt.Sprintf(
		`SELECT comments, title, n.id,
      to_char(startTime AT TIME ZONE 'UTC', 'YYYY-MM-DD HH24:MI:SS'),
      to_char(endtime AT TIME ZONE 'UTC', 'YYYY-MM-DD HH24:MI:SS'),
      longitude, latitude, u.id, u.name, u.email
    FROM notes as n
    JOIN notesusers as nu ON n.id = nu.note_id
    JOIN users as u ON nu.user_id = u.id
    %s`, whereClause,
	)

	// Get rows of (...note, ...user)
	notesWithUsersRows, uErr := db.Query(notesWithUsersQuery)
	if uErr != nil {
		log.Println("getting user rows failed")
		return nil, uErr
	}
	defer notesWithUsersRows.Close()

	result, notesById, uErr := userRowsToNotes(notesWithUsersRows)
	if uErr != nil {
		log.Println("userRowsToNotes failed")
		return nil, uErr
	}
	log.Println(fmt.Sprintf("%s", result))

	notesWithTagsQuery := fmt.Sprintf(
		`SELECT n.id, t.tag
    FROM notes as n
    LEFT JOIN notestags as nt ON n.id = nt.note_id
    LEFT JOIN tags as t ON nt.tag_id = t.id
    WHERE n.id IN %s`, result,
	)

	// Get rows of (note.id, tag)
	notesWithTagsRows, tErr := db.Query(notesWithTagsQuery)
	if tErr != nil {
		log.Println("getting note's tags rows failed")
		return nil, tErr
	}
	defer notesWithTagsRows.Close()

	return tagsRowsToNotes(notesById, notesWithTagsRows)
}

type reverseChronologicalOrder []Note

func (a reverseChronologicalOrder) Len() int {
	return len(a)
}

func (a reverseChronologicalOrder) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a reverseChronologicalOrder) Less(i, j int) bool {
	if *a[i].StartTime > *a[j].StartTime {
		return true
	}
	if *a[i].StartTime < *a[j].StartTime {
		return false
	}
	return *a[i].EndTime > *a[j].EndTime
}

func userRowsToNotes(notesWithUsersRows *sql.Rows) (string, map[int64]Note, error) {
	/*
	   Loop over notes with users, populating each note's Users field. As this is
	   done, insert the notes into a hash map by their id. When iterating over
	   notes with tags, insert the tags into the note struct taken from the map
	   with the correct id.
	*/
	notesById := make(map[int64]Note)

	// Iterate over the notesWithUsersRows, populating notesById and each note's
	// users field.
	var emptyNote Note
	var note Note
	var user User
	noteIds := make([]int64, 0)

	for notesWithUsersRows.Next() {
		err := notesWithUsersRows.Scan(
			&note.Comment,
			&note.Title,
			&note.Id,
			&note.StartTime,
			&note.EndTime,
			&note.Longitude,
			&note.Latitude,
			&user.Id,
			&user.Name,
			&user.Email,
		)
		if err != nil {
			return "", notesById, err
		}

		noteIds = append(noteIds, *note.Id)
		// If not already hit this note, add it to the map and initialise its users.
		// Else, get the note from the map and add this user to its users.
		noteWithUsers := notesById[*note.Id]
		if noteWithUsers == emptyNote {
			note.Users = &[]User{user}
			notesById[*note.Id] = note
		} else {
			noteUsers := noteWithUsers.Users
			*noteUsers = append(*noteUsers, user)
			notesById[*note.Id] = noteWithUsers
		}
	}
	result := ConvertIntArrayToString(noteIds)
	return result, notesById, nil
}

/**
 * Takes rows of (...note, ...user) and (note.id, tag) and constructs a slice
 * of note objects with the tag and user arrays filled in.
 */
func tagsRowsToNotes(notesById map[int64]Note, notesWithTagsRows *sql.Rows) ([]Note, error) {

	var note Note
	var emptyNote Note
	// Iterate over notes with tags rows, adding tags to the notes from the map,
	// whose Users field has been created.
	var tag *string
	for notesWithTagsRows.Next() {
		err := notesWithTagsRows.Scan(&note.Id, &tag)
		if err != nil {
			return nil, err
		}

		// If tags not yet constructed, construct it, else append.
		noteWithUsers := notesById[*note.Id]
		if noteWithUsers == emptyNote {
			// FIXME: wtf why is this happening?
			log.Printf("Error in models.rowsToResults(): Found note with tag but not user: %+v", note)
		} else if tag == nil && noteWithUsers.Tags == nil {
			var emptyTags []string = make([]string, 0)
			noteWithUsers.Tags = &emptyTags
			//noteWithUsers.Tags = &[]string{}
		} else if tag != nil && noteWithUsers.Tags == nil {
			noteWithUsers.Tags = &[]string{*tag}
		} else if tag != nil {
			*noteWithUsers.Tags = append(*noteWithUsers.Tags, *tag)
		}
		if noteWithUsers != emptyNote {
			notesById[*note.Id] = noteWithUsers
		}
	}

	// Convert map to slice
	// FIXME: Is there a way to do this without required this extra iteration
	//        over all notes?
	var notes []Note = make([]Note, 0)
	for _, note := range notesById {
		if note.Tags == nil {
		}
		notes = append(notes, note)
	}

	// Sorting notes
	sort.Sort(reverseChronologicalOrder(notes))

	return notes, nil
}

func getNotesActiveByUser(userEmail string) ([]Note, error) {
	s := fmt.Sprintf("WHERE u.email = '%s'", userEmail)
	log.Println(s)
	return filterNotes(s)
}

func getNotesActiveAtTime(time string) ([]Note, error) {
	s := fmt.Sprintf("WHERE (starttime <= '%[1]s' AND endtime >= '%[1]s')", time)
	log.Println(s)
	return filterNotes(s)
}

func getAllNotes() ([]Note, error) {
	return filterNotes("")
}

func printNote(n Note) {
	log.Println("Printing note...")
	log.Println(*n.Title)
	log.Println(*n.Comment)
	log.Println(*n.StartTime)
	log.Println(*n.EndTime)
	log.Println(*n.Longitude)
	log.Println(*n.Latitude)
	log.Println(*n.Id)
	log.Println(*n.Users)
	log.Println(*n.Tags)
}

func TimeForAggregation() bool {
	var valid bool = false
	insertionNoteCounter.Lock()
	if insertionNoteCounter.counter == insertionNoteCounter.target {
		insertionNoteCounter.target = randomRange(1, 5)
		log.Printf("Next merge attempt will occur in %d inserts\n", insertionNoteCounter.target)
		//Set the counter to -1 as after a merge, we call createNote which will
		//increment the counter
		insertionNoteCounter.counter = -1
		valid = true
	} else {
		log.Printf("Current insertion target: %d", insertionNoteCounter.target)
		log.Printf("Current insertion count: %d", insertionNoteCounter.counter)
	}
	insertionNoteCounter.Unlock()
	return valid
}

/* Gets all notes in Notes table that are within a certain range
radius - metres
longitude - degrees
latitude - degrees
*/
func GetNotesWithinRange(radius float64, note Note) (notes []Note, err error) {
	result := make([]Note, 0)
	result = append(result, note)
	latitude := *note.Latitude
	longitude := *note.Longitude

	notes, err = getAllNotes()
	if err != nil {
		return result, err
	}

	for i := 0; i < len(notes); i++ {
		if *notes[i].Id == *note.Id {
			continue
		}
		distance := greatCircleDistance(latitude, longitude, *notes[i].Latitude, *notes[i].Longitude)
		if distance <= radius {
			result = append(result, notes[i])
		}
	}

	return result, nil
}

func degToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

//calculates shortest distance of two spherical co-ordinates in metres using Haversine formula
func greatCircleDistance(plat1 float64, plong1 float64, plat2 float64, plong2 float64) float64 {
	var EARTH_RADIUS float64 = 6371000 //metres
	dLat := degToRadians(plat2 - plat1)
	dLon := degToRadians(plong2 - plong1)

	lat1 := degToRadians(plat1)
	lat2 := degToRadians(plat2)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}

func randomRange(min int, max int) (result int) {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

/* Gets all notes in the list that occur around the same time as the head of
the list
Policy - Get notes where start time and end time are within 15 mins of each
         other
*/
func GetAllNotesAroundSameTime(notes []Note) (filter []Note, err error) {
	result := make([]Note, 0)
	if len(notes) == 0 {
		filterErr := fmt.Errorf("Empty notes list was passed to GetAllNotesAroundSameTime")
		return result, filterErr
	}

	//Add original new note
	newNote := notes[0]
	result = append(result, newNote)

	for i := 1; i < len(notes); i++ {
		if withinTimeFrame(*notes[i].StartTime, *newNote.StartTime, 15*60) &&
			withinTimeFrame(*notes[i].EndTime, *newNote.EndTime, 15*60) {
			result = append(result, notes[i])
		}
	}
	return result, nil
}

/*
  Checks if two dates are within a time frame of each other
  t1 - First date
  t2 - Second date
  dt - time frame in seconds
*/
func withinTimeFrame(t1 string, t2 string, dt int64) bool {
	format := "2006-01-02 15:04:05"
	pT1, err := time.Parse(format, t1)
	if err != nil {
		log.Println("Error parsing time t1 ", t1)
		log.Println(err)
		return false
	}
	pT2, err := time.Parse(format, t2)
	if err != nil {
		log.Println("Error parsing time t2 ", t2)
		log.Println(err)
		return false
	}
	absDiff := time.Duration(math.Abs(float64(pT1.Sub(pT2)))) * time.Nanosecond
	return absDiff < time.Duration(dt)*time.Second
}

/* Gets all notes in Notes table that have similar titles/comments */
func GetNotesWithSimilarText(notes []Note) []Note {
	result := make([]Note, 0)
	comparatorNote := notes[0]
	for i := 1; i < len(notes); i++ {
		note := notes[i]
		if areSimilarStrings(*note.Title, *comparatorNote.Title) {
			result = append(result, notes[i])
		}
	}
	if len(result) > 0 {
		result = append([]Note{comparatorNote}, result...)
	}
	return result
}

func areSimilarStrings(s1 string, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	if len(s1) == 0 {
		return true
	}
	s1, s2 = strings.Title(s1), strings.Title(s2)
	return strings.ToLower(s1) == strings.ToLower(s2)
}

type noteOccurs struct {
	note   Note
	occurs int
}

/*
  Iterates through notes and returns notes that have similar tags given that
  they are aroung the same location and have roughly the same title
*/
func GetNotesWithSimilarTags(notes []Note) (filtered []Note, err error) {
	// For each index i in upperlist,
	// note i has notes that have similar tags to it
	var upperList [][]noteOccurs = make([][]noteOccurs, len(notes))
	lowerList := []noteOccurs{}

	filtered = make([]Note, 0)

	// Store maximum number of notes with similar tags
	index_of_notes_similar := 0
	max_num_of_notes_similar := 0

	// Go through each note,
	// for each note i, check all the other notes to see if they have similar tags
	// lowerlist contains all tags similar to note i
	// And then upperList[i] contains lowerlist

	for i := 0; i < len(notes); i++ {
		for j := i + 1; j < len(notes); j++ {
			similar, occ := notesHaveSimilarTags(notes[i], notes[j])
			if similar {
				noteOccurs := noteOccurs{note: notes[j], occurs: occ}
				lowerList = append(lowerList, noteOccurs)
			}
		}
		upperList[i] = lowerList

		// Check and keep track of maximum number of notes with similar tags
		len := len(lowerList)
		if len >= max_num_of_notes_similar {
			max_num_of_notes_similar = len
			index_of_notes_similar = i
		}
		// Reset lowerlist to refill in next forloop
		lowerList = nil
	}

	if max_num_of_notes_similar > 0 {
		// Return maximum number of notes with similar tags
		lowerList = upperList[index_of_notes_similar]
		filtered = append(filtered, notes[index_of_notes_similar])
		for _, nList := range lowerList {
			filtered = append(filtered, nList.note)
		}
	}
	// else max_num_of_notes_similar = 0
	// Every tag is different - so no notes are similar, so return empty struct

	return filtered, nil
}

/*
  Function which returns true if each note has a least
  one tag in common with the other.
  It also returns number of tags it has in common.
*/
func notesHaveSimilarTags(n Note, m Note) (bool, int) {
	tags_n := *n.Tags
	tags_m := *m.Tags

	similar := false
	num_common := 0

	for _, tn := range tags_n {
		for _, tm := range tags_m {
			if tn == tm {
				num_common++
				similar = true
			}
		}
	}
	return similar, num_common
}

/*Returns a list of note ids and a new note using predefined policies
  notes - list of notes to be aggregated. Must not be empty
*/
func ConstructAggregatedNote(notes []Note) (note_ids []int64, note Note) {
	length := len(notes)
	for i := 0; i < length; i++ {

		note_ids = append(note_ids, int64(*notes[i].Id))
	}
	log.Println("MERGING NOTES WITH IDS: ", note_ids)

	var n Note
	n.Title = aggregateTitle(notes)
	n.Comment = aggregateComments(notes, length)
	lat, long := aggregateCoordinates(notes, length)
	n.Latitude = lat
	n.Longitude = long
	startTime, endTime := aggregateTimes(notes, length)
	n.StartTime = startTime
	n.EndTime = endTime
	n.Users = aggregateUsers(notes, length)
	log.Println("Cons Users: ", *n.Users)
	n.Tags = aggregateTags(notes, length)
	log.Println("Cons Tags: ", *n.Tags)
	log.Println("Number of notes merged: ", length)

	return note_ids, n
}

//Our policy for title aggregation
func aggregateTitle(notes []Note) *string {
	//Takes the first note's title. Every note must contain a title so we can
	//do this
	return notes[0].Title
}

//Our policy for comments aggregation
func aggregateComments(notes []Note, length int) *string {
	//Takes the first comment we find
	i := 0
	for i < length && *notes[i].Comment == "" {
		i++
	}
	if i < length {
		return notes[i].Comment
	} else {
		//None of the notes had any comments
		res := ""
		return &res
	}

}

//Our policy for latitude and longitude aggregation
func aggregateCoordinates(notes []Note, length int) (*float64, *float64) {
	//Take the average latitude and longitudes
	var accLat float64 = 0
	var accLong float64 = 0
	for j := 0; j < length; j++ {
		accLat += *notes[j].Latitude
		accLong += *notes[j].Longitude
	}
	lat := accLat / float64(length)
	long := accLong / float64(length)
	return &lat, &long
}

type dateSlice []Date

func (p dateSlice) Len() int {
	return len(p)
}

func (p dateSlice) Less(i, j int) bool {
	return p[i].time.Before(p[j].time)
}

func (p dateSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Date struct {
	str  string
	time time.Time
}

func aggregateTimes(notes []Note, length int) (finalStartTime *string, finalEndTime *string) {
	format := "2006-01-02 15:04:05"

	//Sort the time and select the median
	startTimes := make([]Date, 0)
	endTimes := make([]Date, 0)
	for i := 0; i < length; i++ {
		startTime, _ := time.Parse(format, *notes[i].StartTime)
		date := Date{str: *notes[i].StartTime, time: startTime}
		startTimes = append(startTimes, date)
		endTime, _ := time.Parse(format, *notes[i].EndTime)
		date = Date{str: *notes[i].EndTime, time: endTime}
		endTimes = append(endTimes, date)
	}
	medST := medianTimes(startTimes)
	medET := medianTimes(endTimes)
	return &medST, &medET
}

func medianTimes(dates []Date) string {
	length := len(dates)
	sort.Sort(dateSlice(dates))
	if length%2 == 0 {
		lowerBound := dates[length/2-1]
		upperBound := dates[length/2]
		diff := upperBound.time.Sub(lowerBound.time)
		halfDur := diff / 2
		median := lowerBound.time.Add(halfDur)
		format := "2006-01-02 15:04:05"
		return median.Format(format)
	} else {
		return dates[length/2].str
	}
}

//Our policy for aggregating users
func aggregateUsers(notes []Note, length int) *[]User {
	//Combine all distinct users
	var users = make([]User, 0)
	seenSet := make(map[string]struct{})
	for i := 0; i < length; i++ {
		for j := 0; j < len(*notes[i].Users); j++ {
			_, found := seenSet[(*notes[i].Users)[j].Email]
			if !found {
				users = append(users, (*notes[i].Users)[j])
				seenSet[(*notes[i].Users)[j].Email] = struct{}{}
			}
		}
	}
	return &users
}

//Our policy for aggregating tags
func aggregateTags(notes []Note, length int) *[]string {
	//Combine all distinct tags
	var tags = make([]string, 0)
	seenSet := make(map[string]struct{})
	for i := 0; i < length; i++ {
		for j := 0; j < len(*notes[i].Tags); j++ {
			_, found := seenSet[(*notes[i].Tags)[j]]
			if !found {
				tags = append(tags, (*notes[i].Tags)[j])
				seenSet[(*notes[i].Tags)[j]] = struct{}{}
			}
		}
	}
	return &tags
}
