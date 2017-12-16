package main

import (
  "log"
  "time"
  "fmt"
  "strconv"

  "net/http"
  "encoding/json"
  "gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
)

type Venue struct {
  Address struct {
    Latitude *string  `json:"latitude,omitempty"`
    Longitude *string `json:"longitude,omitempty"`
  } `json:"address,omitempty"`
}

type EventBriteEvents struct {
  Pagination struct {
    Size  *int                `json:"page_size,omitempty"`
    Count *int                `json:"page_count,omitempty"`
  }                         `json:"pagination,omitempty"`
  Events []struct {
    Name struct {
      Text *string              `json:"text,omitempty"`
    }                         `json:"name,omitempty"`
    Description struct {
      Text *string              `json:"text,omitempty"`
    }                         `json:"description,omitempty"`
    VenueId *string           `json:"venue_id,omitempty"`
    Start struct {
      LocalTime *string         `json:"local,omitempty"`
    }                         `json:"start,omitempty"`
    End struct {
      LocalTime *string         `json:"local,omitempty"`
    }                         `json:"end,omitempty"`
    Url *string               `json:"url,omitempty"`
  }                         `json:"events,omitempty"`
}

func generateUrl(path string) string {
  return "https://www.eventbriteapi.com/v3/" + path + "token=GJ2IG7JJRDXTOF77BOTH"
}

func parseTimeWithT(t time.Time) string {
  exampleLayout := "2006-01-02T03:04:05"
  return t.Format(exampleLayout)
}

func parseTimeWithoutT(timeValue string) string {
  return timeValue[:10] + " " + timeValue[11:]
}

func printNote(n models.Note) {
  log.Println("Printing note...")
  log.Println(*n.Title)
  log.Println(*n.Comment)
  log.Println(*n.StartTime)
  log.Println(*n.EndTime)
  log.Println(*n.Longitude)
  log.Println(*n.Latitude)
  //log.Println(*n.Id)
  log.Println(*n.Users)
  //log.Println(*n.Tags)
}   

func GetEventBriteEvents() {
  t := time.Now()
  // Set startTime to be the currentTime
  startTime := parseTimeWithT(t)
  // Set endTime to be two days after startTime
  endTime := parseTimeWithT(t.AddDate(0, 0, 2))
  eventSearchUrl := "events/search/?location.address=london&start_date.range_start=" +
                startTime + "&start_date.range_end=" + endTime + "&"
  log.Println(generateUrl(eventSearchUrl))
  res, httpGetErr := http.Get(generateUrl(eventSearchUrl))
  if httpGetErr != nil {
    log.Println("HTTP GET request for EventBrite failed")
    return
  }
  var eventsJSON EventBriteEvents
  json.NewDecoder(res.Body).Decode(&eventsJSON)
  count := 0
  numberOfPages := *eventsJSON.Pagination.Count

  for i := 1; i <= numberOfPages; i++ {
    if i > 1 {
      newEventSearchUrl := fmt.Sprintf("%s%s%d%s", eventSearchUrl, "page=", i, "&")
      res, httpGetErr = http.Get(generateUrl(newEventSearchUrl))
      if httpGetErr != nil {
        log.Println("HTTP GET request for EventBrite failed")
        return
      }
      var newEventsJSON EventBriteEvents
      json.NewDecoder(res.Body).Decode(&newEventsJSON)
      eventsJSON = newEventsJSON
    }
    events := eventsJSON.Events
    for j := 0; j < len(events); j++ {
      event := events[j]
      var name string = *event.Name.Text
      var description *string = event.Description.Text
      var eventUrl string = *event.Url
      if description == nil {
        description = &eventUrl
      } else {
        *description += "\n" + eventUrl
      }
      var venueId string = *event.VenueId
      log.Println(venueId)
      venueSearchUrl := generateUrl("venues/" + venueId + "/?")
      log.Println("Venue Search URL: " + venueSearchUrl)
      venueRes, err := http.Get(venueSearchUrl)
      if err != nil {
        log.Println("Couldn't load event venue with id: " + venueId)
        continue
      }
      var venueJSON Venue
      json.NewDecoder(venueRes.Body).Decode(&venueJSON)
      var latitude string = *venueJSON.Address.Latitude
      var longitude string = *venueJSON.Address.Longitude
      var eventStartTime string = parseTimeWithoutT(*event.Start.LocalTime)
      var eventEndTime string = parseTimeWithoutT(*event.End.LocalTime)
      count++
      log.Println(count)
      var newNote models.Note
      newNote.Title = &name
      newNote.Comment = description
      newNote.StartTime = &eventStartTime
      newNote.EndTime = &eventEndTime
      lat, _ := strconv.ParseFloat(latitude, 64)
      long, _ := strconv.ParseFloat(longitude, 64)
      newNote.Latitude = &lat
      newNote.Longitude = &long
      var eventBriteUser models.User
      eventBriteUser.Id = -1
      eventBriteUser.Name = "EventBrite"
      eventBriteUser.Email = "eventbrite@example.com"
      eventBriteUser.Picture = "https://tctechcrunch2011.files.wordpress.com/2014/06/eventbrite_1.jpg?w=700"
      newUsers := []models.User{eventBriteUser}
      newNote.Users = &newUsers

      printNote(newNote)
      models.Notes.Create(&newNote)
    }
  }
}
